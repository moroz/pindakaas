document.querySelectorAll<HTMLElement>("[data-copy]").forEach((btn) => {
  btn.addEventListener("click", () => navigator.clipboard.writeText(btn.dataset.copy!));
});
