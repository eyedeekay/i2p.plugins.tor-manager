//go:build darwin
// +build darwin

package main

var running = false
var shutdown = false

func runSysTray(down bool) {
	if !running {
		running = true
		shutdown = down
	}
}
