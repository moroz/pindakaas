// Supports weights 100-700
import "@fontsource-variable/ibm-plex-sans/wght.css";
import "temporal-polyfill/global";
import { mount } from "svelte";
import "./style.css";
import TunnelIndex from "./components/TunnelIndex.svelte";
import CopyButton from "./components/CopyButton.svelte";

document.querySelectorAll<HTMLElement>("[data-copy]").forEach((btn) => {
  const classNames = btn.dataset.class ?? "";

  mount(CopyButton, {
    props: { text: btn.dataset.copy!, class: classNames },
    target: btn,
  });
});

function initTable() {
  const propsEl = document.querySelector("#index-table-props");
  if (!propsEl) return;

  const target = document.querySelector("#svelte-root");
  if (!target) return;

  let initialProps;
  try {
    initialProps = JSON.parse(propsEl.innerHTML);
  } catch (e) {
    console.error(e);
    return;
  }

  mount(TunnelIndex, {
    props: initialProps,
    target,
  });
}

initTable();
