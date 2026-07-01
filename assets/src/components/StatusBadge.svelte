<script lang="ts">
  interface Props {
    active: boolean;
  }

  let { active }: Props = $props();
</script>

<span class="badge" class:active>
  <svg class="fill-current w-5 h-5" viewBox="0 0 640 640">
    <use href={`/assets/${active ? "person-running" : "bed"}.svg#icon`} />
  </svg>
  <span class="inner">
    {#if active}
      Online
    {:else}
      Inactive
    {/if}
  </span>
</span>

<style>
  @reference "../style.css";

  .badge {
    @apply my-1 inline-flex gap-1 rounded-full bg-slate-200 px-3 py-1 font-sans font-bold text-slate-800 uppercase w-32;

    :is(&).active {
      @apply bg-green-100 text-green-800;
    }
  }

  .inner {
    @apply text-center flex-1;
  }
</style>
