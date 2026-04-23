document.addEventListener("htmx:afterSwap", () => {
  document.querySelectorAll(".copyable").forEach((el) => {
    el.addEventListener("click", () => {
      navigator.clipboard.writeText(el.innerText);
    });
  });
});

document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll(".file-upload").forEach((zone) => {
    zone.addEventListener("dragover", (e) => {
      e.preventDefault();
      zone.classList.add("dragging");
    });
    zone.addEventListener("dragleave", () => {
      zone.classList.remove("dragging");
    });
    zone.addEventListener("drop", (e) => {
      e.preventDefault();
      zone.classList.remove("dragging");
      const files = e.dataTransfer.files;
      console.log("hello from dragging");
    });
  });
});
