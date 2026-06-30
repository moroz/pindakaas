#!/usr/bin/env -S deno run --allow-net

// WebSocket client that stands in for Twilio: it streams messages to the server
// and reports, per direction, whether traffic actually flowed. Pair it with
// scripts/ws-echo-server.ts.
//
// It distinguishes two things the server does:
//   - "echo:" replies   -> proves the client->server->client round trip works.
//   - "push" messages   -> proves the server->client (unprompted) stream works.
//                          This is the direction a Twilio Media Streams bot uses
//                          to send audio, and the one we suspect is broken.
//
// Usage:
//   deno run --allow-net scripts/ws-demo-client.ts <url> [--interval 250] [--seconds 8]
//
// Examples:
//   # direct, no tunnel (baseline):
//   deno run --allow-net scripts/ws-demo-client.ts ws://localhost:8787/ws
//   # through the tunnel (the real test):
//   deno run --allow-net scripts/ws-demo-client.ts wss://atrocious-jaguar.pindakaas.virtualq.run/ws

const url = Deno.args.find((a) => a.startsWith("ws://") || a.startsWith("wss://"));
if (!url) {
  console.error("usage: ws-demo-client.ts <ws(s)://url> [--interval MS] [--seconds N]");
  Deno.exit(2);
}

const flag = (name: string, fallback: number): number => {
  const i = Deno.args.indexOf(`--${name}`);
  return i >= 0 ? Number(Deno.args[i + 1]) : fallback;
};

const interval = flag("interval", 250);
const seconds = flag("seconds", 8);
const size = flag("size", 0); // bytes of payload padding, to mimic audio frames
const pad = "x".repeat(size);

let sent = 0;
let echoes = 0;
let pushes = 0;
const latencies: number[] = []; // round-trip ms, from echoed timestamps
let sendTimer: number | undefined;

console.log(`connecting to ${url} (stream ${seconds}s, every ${interval}ms) ...`);
const ws = new WebSocket(url);
const t0 = performance.now();

ws.addEventListener("open", () => {
  console.log(`✓ connected in ${(performance.now() - t0).toFixed(0)}ms`);
  sendTimer = setInterval(() => {
    sent += 1;
    ws.send(JSON.stringify({ kind: "client", seq: sent, t: Date.now(), pad }));
  }, interval);
  setTimeout(() => ws.close(1000, "done"), seconds * 1000);
});

ws.addEventListener("message", (e) => {
  const data = String(e.data);
  if (data.startsWith("echo:")) {
    echoes += 1;
    try {
      const orig = JSON.parse(data.slice("echo:".length));
      if (typeof orig.t === "number") latencies.push(Date.now() - orig.t);
    } catch { /* ignore non-JSON echoes */ }
  } else {
    // server-initiated push
    pushes += 1;
  }
});

ws.addEventListener("error", (e) => {
  console.error(`✗ error after ${(performance.now() - t0).toFixed(0)}ms: ${(e as ErrorEvent).message ?? e}`);
});

ws.addEventListener("close", (e) => {
  clearInterval(sendTimer);

  const avg = latencies.length
    ? (latencies.reduce((a, b) => a + b, 0) / latencies.length).toFixed(0)
    : "n/a";
  const max = latencies.length ? Math.max(...latencies).toFixed(0) : "n/a";

  console.log("\n=== summary ===");
  console.log(`close: code=${e.code} clean=${e.wasClean}${e.reason ? ` reason="${e.reason}"` : ""}`);
  console.log(`client -> server -> client (echoes): sent ${sent}, got ${echoes} back  ${verdict(echoes)}`);
  console.log(`server -> client (unprompted pushes): got ${pushes}            ${verdict(pushes)}`);
  console.log(`echo round-trip latency: avg ${avg}ms, max ${max}ms`);

  if (echoes > 0 && pushes === 0) {
    console.log("\n→ client->server works, but server->client push is DEAD. This is the Twilio-bot symptom.");
  } else if (echoes === 0 && pushes > 0) {
    console.log("\n→ server->client works, but client->server is DEAD.");
  } else if (echoes === 0 && pushes === 0) {
    console.log("\n→ no data flowed either way after the upgrade (handshake-only).");
  } else {
    console.log("\n→ both directions carried data. If this is through the tunnel, the proxy is fine for streaming.");
  }

  Deno.exit(echoes > 0 && pushes > 0 ? 0 : 1);
});

function verdict(n: number): string {
  return n > 0 ? "OK" : "FAIL";
}
