// Package platform describes capabilities of the current game target.
package platform

import "runtime"

const Browser = runtime.GOOS == "js"
