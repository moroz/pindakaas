#!/usr/bin/env -S deno run --allow-net

// Try connecting to a WebSocket tunnel and exercise it: send a message on open,
// keep sending one every few seconds, log everything received, and report any
// error or close. Useful for verifying the reverse-proxy WebSocket upgrade path
// end to end.
//
// Usage:
//   deno run --allow-net scripts/ws-test.ts [url] [--count N] [--interval MS]
//   ./scripts/ws-test.ts wss://succulent-chinese-meal.pindakaas.virtualq.run/ws
//
// Exits 0 on a clean close, 1 on error or unclean close.

const args = Deno.args.filter((a) => !a.startsWith("--"));
const url = args[0] ??
  "wss://succulent-chinese-meal.pindakaas.virtualq.run/ws";

const flag = (name: string, fallback: number): number => {
  const i = Deno.args.indexOf(`--${name}`);
  return i >= 0 ? Number(Deno.args[i + 1]) : fallback;
};

const count = flag("count", 3); // messages to send before closing
const intervalMs = flag("interval", 2000); // delay between messages

console.log(`Connecting to ${url} ...`);

const ws = new WebSocket(url);
const started = performance.now();
let sent = 0;
let timer: number | undefined;

const elapsed = () => `${(performance.now() - started).toFixed(0)}ms`;

ws.addEventListener("open", () => {
  console.log(`✓ connected in ${elapsed()}`);

  const send = () => {
    sent += 1;
    const msg = `hello #${sent} from deno`;
    console.log(`→ ${msg}`);
    ws.send(msg);
    if (sent >= count) {
      clearInterval(timer);
      // Give the last echo a moment to come back, then close cleanly.
      setTimeout(() => ws.close(1000, "done"), 500);
    }
  };

  send();
  timer = setInterval(send, intervalMs);
});

ws.addEventListener("message", (e) => {
  console.log(`← ${e.data}`);
});

ws.addEventListener("error", (e) => {
  const msg = (e as ErrorEvent).message ?? String(e);
  console.error(`✗ websocket error after ${elapsed()}: ${msg}`);
});

ws.addEventListener("close", (e) => {
  clearInterval(timer);
  const reason = e.reason ? `, reason="${e.reason}"` : "";
  console.log(`connection closed (code=${e.code}${reason}, clean=${e.wasClean})`);
  Deno.exit(e.wasClean ? 0 : 1);
});

// Hard safety net so the script can never hang forever.
setTimeout(() => {
  console.error("timeout reached without a clean close, giving up");
  try {
    ws.close(1000, "client timeout");
  } catch { /* ignore */ }
  Deno.exit(1);
}, intervalMs * (count + 2) + 5000);
