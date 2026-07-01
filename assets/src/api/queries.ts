import type { TunnelListDTO, TunnelListResponse } from "../interfaces";

export async function fetchTunnels(): Promise<ReadonlyArray<TunnelListDTO>> {
  const response = await fetch("/api/tunnels", { credentials: "include" });
  if (!response.ok) {
    throw new Error(`Failed to fetch tunnels: ${response.status}`);
  }
  const body: TunnelListResponse = await response.json();
  return body.data;
}
