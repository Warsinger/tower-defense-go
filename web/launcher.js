"use strict";
// Loaded before the iframe so even a cached game's first ready message is seen.
window.addEventListener("message", (event) => {
  const frame = document.getElementById("game");
  const status = document.getElementById("status");
  if (!frame || !status || event.origin !== location.origin || event.source !== frame.contentWindow) return;
  if (event.data?.type === "game-ready" && !status.dataset.error) {
    status.textContent = "Game ready. Click inside, then press Space or select Start Game to play.";
    frame.focus({preventScroll: true});
  }
  if (event.data?.type === "game-error") {
    status.dataset.error = "true";
    status.textContent = "Game could not run. Reload to retry. Details: " + event.data.message;
  }
});
