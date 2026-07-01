<script lang="ts">
  import { onMount } from "svelte";
  import type { TunnelListDTO } from "../interfaces";
  import { fetchTunnels } from "../api/queries";
  import StatusBadge from "./StatusBadge.svelte";

  let {
    tunnels: initialTunnels,
    baseDomain,
  }: { tunnels: TunnelListDTO[]; baseDomain: string } = $props();

  let tunnels = $state<ReadonlyArray<TunnelListDTO>>([]);

  onMount(() => {
    tunnels = initialTunnels;

    const interval = setInterval(async () => {
      tunnels = await fetchTunnels();
    }, 5000);
    return () => clearInterval(interval);
  });

  function fqdn(subdomain: string): string {
    return `https://${subdomain}.${baseDomain}`;
  }

  function formatDate(iso: string): string {
    const zdt = Temporal.Instant.from(iso).toZonedDateTimeISO("Europe/Berlin");
    const mo = String(zdt.month).padStart(2, "0");
    const d = String(zdt.day).padStart(2, "0");
    const h = String(zdt.hour).padStart(2, "0");
    const mi = String(zdt.minute).padStart(2, "0");
    const s = String(zdt.second).padStart(2, "0");
    return `${zdt.year}-${mo}-${d} ${h}:${mi}:${s}`;
  }
</script>

<div
  class="overflow-hidden rounded-lg border border-slate-300 bg-white shadow-sm"
>
  <table class="index-table w-full">
    <thead>
      <tr>
        <th class="w-32 text-center">Status</th>
        <th class="w-32 text-center">Username</th>
        <th>Subdomain</th>
        <th class="w-58 text-right">Created at</th>
      </tr>
    </thead>
    <tbody>
      {#each tunnels as tunnel}
        <tr
          data-url="/tunnels/{tunnel.id}"
          onclick={() => (location.href = `/tunnels/${tunnel.id}`)}
        >
          <td class="text-center"><StatusBadge active={tunnel.active} /></td>
          <td class="font-mono text-center">{tunnel.username}</td>
          <td class="font-mono">
            <span title={fqdn(tunnel.subdomain)}>{tunnel.subdomain}</span>
          </td>
          <td class="text-slate-500 text-right"
            >{formatDate(tunnel.insertedAt)}</td
          >
        </tr>
      {/each}
    </tbody>
  </table>
</div>
