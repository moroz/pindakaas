<script lang="ts">
  import { cn } from "../lib/cn";

  interface Props {
    text: string;
    class?: string;
  }

  let { text, class: className = "" }: Props = $props();

  let success = $state(false);

  function copy(e: MouseEvent) {
    e.stopPropagation();
    navigator.clipboard.writeText(text);
    success = true;
    setTimeout(() => {
      success = false;
    }, 3000);
  }
</script>

<button
  class={cn(
    "button secondary relative font-sans h-7 px-0 text-xs font-normal",
    className,
  )}
  onclick={copy}
>
  <!-- Invisible ghost of the wider state — drives button width -->
  <span class="invisible flex items-center gap-1">
    <svg class="fill-current w-5 h-5" viewBox="0 0 640 640">
      <use href="/assets/check.svg#icon" />
    </svg>
    Copied!
  </span>

  <!-- Copy state -->
  <span
    class={cn(
      "absolute inset-0 flex items-center justify-center gap-1 transition-opacity duration-150",
      success ? "opacity-0 pointer-events-none" : "opacity-100",
    )}
  >
    <svg class="fill-current w-5 h-5" viewBox="0 0 640 640">
      <use href="/assets/copy.svg#icon" />
    </svg>
    Copy
  </span>

  <!-- Copied state -->
  <span
    class={cn(
      "absolute inset-0 flex items-center justify-center gap-1 transition-opacity duration-300",
      success ? "opacity-100" : "opacity-0 pointer-events-none",
    )}
  >
    Copied!
  </span>
</button>
