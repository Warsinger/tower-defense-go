import {startGoGame} from './shared/go-runtime.mjs';
startGoGame({
  wasmURL: new URL('game.wasm', import.meta.url),
  readyCallback: 'towerDefenseReady',
  stoppedMessage: 'Game stopped. Reload the page to restart.',
  preventKeys: ['Space', 'ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight']
});
