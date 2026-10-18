"use strict";
function reportFailure(error) {
  console.error("Tower Defense runtime failure", error);
  parent.postMessage({type: "game-error", message: String(error.message || error)}, location.origin);
}
window.addEventListener("error", event => reportFailure(event.error || event.message));
window.addEventListener("unhandledrejection", event => reportFailure(event.reason));
window.towerDefenseReady = () => parent.postMessage({type: "game-ready"}, location.origin);
window.addEventListener("keydown", event => {
  if (["Space", "ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight"].includes(event.code)) event.preventDefault();
});
(async () => {
  try {
    const go = new Go();
    const response = await fetch("game.wasm");
    if (!response.ok) throw new Error("Game download failed: HTTP " + response.status);
    const result = await WebAssembly.instantiate(await response.arrayBuffer(), go.importObject);
    await go.run(result.instance);
    throw new Error("Game stopped. Reload the page to restart.");
  } catch (error) { reportFailure(error); }
})();
