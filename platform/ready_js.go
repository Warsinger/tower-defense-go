//go:build js && wasm

package platform

import "syscall/js"

var ready bool

// Ready reports the first frame, rather than merely a successful Wasm download.
func Ready() {
	if ready {
		return
	}
	ready = true
	callback := js.Global().Get("towerDefenseReady")
	if callback.Type() == js.TypeFunction {
		callback.Invoke()
	}
}
