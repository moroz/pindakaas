#!/usr/bin/env -S deno run --allow-net

// Minimal WebSocket server for diagnosing the tunnel's WS proxy.
//
// It does two independent things, to mirror a streaming backend like a Twilio
// Media Streams bot:
//   1. Echoes every message it receives, prefixed with "echo:".
//   2. Pushes a "push" message every --interval ms WITHOUT being prompted.
//
// Run it locally, then expose it through a tunnel, e.g.:
//   deno run --allow-net scripts/ws-echo-server.ts            # listens on :8787
//   ssh -R 0:localhost:8787 <tunnel-user>@pindakaas.virtualq.run -tt
//
// Usage: deno run --allow-net scripts/ws-echo-server.ts [--port 8787] [--interval 250]

const flag = (name: string, fallback: number): number => {
  const i = Deno.args.indexOf(`--${name}`);
  return i >= 0 ? Number(Deno.args[i + 1]) : fallback;
};

const port = flag("port", 8787);
const interval = flag("interval", 250);
const size = flag("size", 0); // bytes of payload padding, to mimic audio frames
const pad = "x".repeat(size);

Deno.serve({ port, hostname: "0.0.0.0" }, (req) => {
  const url = new URL(req.url);

  // Plain HTTP hits (health checks, manual curl) get a simple 200.
  if (req.headers.get("upgrade")?.toLowerCase() !== "websocket") {
    return new Response(`ok ${url.pathname}\n`);
  }

  const { socket, response } = Deno.upgradeWebSocket(req);
  let pushSeq = 0;
  let received = 0;
  let timer: number | undefined;

  socket.onopen = () => {
    console.log(`[server] client connected on ${url.pathname}`);
    timer = setInterval(() => {
      pushSeq += 1;
      socket.send(JSON.stringify({ kind: "push", seq: pushSeq, t: Date.now(), pad }));
      console.log(`[server] → push #${pushSeq}`);
    }, interval);
  };

  socket.onmessage = (e) => {
    received += 1;
    console.log(`[server] ← ${e.data}`);
    socket.send(`echo:${e.data}`);
  };

  socket.onclose = (e) => {
    clearInterval(timer);
    console.log(
      `[server] closed (code=${e.code}, received ${received} msgs, pushed ${pushSeq})`,
    );
  };

  socket.onerror = (e) => {
    console.error(`[server] error:`, (e as ErrorEvent).message ?? e);
  };

  return response;
});

console.log(`ws echo+push server listening on :${port} (push every ${interval}ms)`);
