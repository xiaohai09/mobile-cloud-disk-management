package utils

import (
	"log"
	"runtime/debug"
)

// SafeGo starts a goroutine that recovers from panics and logs them.
func SafeGo(name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] %s: %v\n%s", name, r, debug.Stack())
			}
		}()
		fn()
	}()
}
