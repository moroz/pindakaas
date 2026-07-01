import "temporal-polyfill/global";
import { mount } from "svelte";
import "./style.css";
import TunnelIndex from "./tunnel-index.svelte";

document.querySelectorAll<HTMLElement>("[data-copy]").forEach((btn) => {
  btn.addEventListener("click", () => navigator.clipboard.writeText(btn.dataset.copy!));
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
