"use strict";
const gameURL = new URL("game.html", document.currentScript.src);
const frame = document.getElementById("game");
const status = document.getElementById("status");
document.getElementById("launch").addEventListener("click", (event) => {
  event.currentTarget.hidden = true;
  frame.hidden = false;
  status.textContent = "Loading game…";
  frame.src = gameURL.href;
});
window.addEventListener("message", (event) => {
  if (event.origin !== location.origin || event.source !== frame.contentWindow) return;
  if (event.data?.type === "game-ready") status.textContent = "Game ready. Click inside to play.";
  if (event.data?.type === "game-error") {
    status.dataset.error = "true";
    status.textContent = "Game could not run. Reload to retry. Details: " + event.data.message;
  }
});
