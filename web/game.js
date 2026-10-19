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
    let bytes = await response.arrayBuffer();
    const magic = new Uint8Array(bytes, 0, Math.min(2, bytes.byteLength));
    if (magic[0] === 0x1f && magic[1] === 0x8b) {
      if (typeof DecompressionStream === "undefined") throw new Error("This browser cannot decompress the game. Use a current desktop browser.");
      bytes = await new Response(new Blob([bytes]).stream().pipeThrough(new DecompressionStream("gzip"))).arrayBuffer();
    }
    const result = await WebAssembly.instantiate(bytes, go.importObject);
    await go.run(result.instance);
    throw new Error("Game stopped. Reload the page to restart.");
  } catch (error) { reportFailure(error); }
})();
