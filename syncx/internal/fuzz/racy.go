//go:build race
package fuzz

import "runtime"

//go:noinline
func MaybeYield() {
	if (*Manager).Fault() {
		runtime.Gosched()
	}
}
