<script lang="ts">
  import { onMount } from "svelte";
  import type { TunnelListDTO } from "./interfaces";
  import CopyButton from "./components/CopyButton.svelte";
  import { fetchTunnels } from "./api/queries";

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
        <th class="w-32">Status</th>
        <th>Subdomain</th>
        <th class="w-36">Username</th>
        <th class="w-52">Created at</th>
      </tr>
    </thead>
    <tbody>
      {#each tunnels as tunnel}
        <tr
          data-url="/tunnels/{tunnel.id}"
          onclick={() => (location.href = `/tunnels/${tunnel.id}`)}
        >
          <td class="text-center">
            <span class="badge" class:active={tunnel.active}>
              {#if tunnel.active}
                <svg class="fill-current w-5 h-5" viewBox="0 0 640 640">
                  <use href="/assets/person-running.svg#icon" />
                </svg>
                Online
              {:else}
                <svg class="fill-current w-5 h-5" viewBox="0 0 640 640">
                  <use href="/assets/bed.svg#icon" />
                </svg>
                Inactive
              {/if}
            </span>
          </td>
          <td>
            <span title={fqdn(tunnel.subdomain)}>{tunnel.subdomain}</span>
          </td>
          <td class="font-mono">{tunnel.username}</td>
          <td class="font-mono text-slate-400"
            >{formatDate(tunnel.insertedAt)}</td
          >
        </tr>
      {/each}
    </tbody>
  </table>
</div>
