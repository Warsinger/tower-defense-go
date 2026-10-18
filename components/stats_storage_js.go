//go:build js && wasm

package components

import (
	"fmt"
	"syscall/js"
)

const browserStatsKey = "tower-defense.stats.v1"

var sessionStats []byte

// Storage can be unavailable in private sessions or when site storage is blocked.
// Keep the current session usable and never touch the native score file.
// Reflect.get must run through Call: unlike Value.Get, Call translates a thrown
// JavaScript property getter into a Go panic that the deferred recover can catch.
func readStats() (data []byte, err error) {
	data = sessionStats
	defer func() {
		if recover() != nil {
			fmt.Println("Browser score storage unavailable; scores last for this session only.")
		}
	}()
	value := js.Global().Get("Reflect").Call("get", js.Global(), "localStorage").Call("getItem", browserStatsKey)
	if !value.IsNull() && !value.IsUndefined() {
		data = []byte(value.String())
		sessionStats = data
	}
	return data, nil
}

func writeStats(data []byte) (err error) {
	sessionStats = append(sessionStats[:0], data...)
	defer func() {
		if recover() != nil {
			fmt.Println("Browser score storage unavailable; scores last for this session only.")
		}
	}()
	js.Global().Get("Reflect").Call("get", js.Global(), "localStorage").Call("setItem", browserStatsKey, string(data))
	return nil
}
