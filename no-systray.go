//go:build nosystray
// +build nosystray

package main

import (
	"flag"
)

var shutdown = false

func Password() string {
	return ""
}

func runSysTray(down bool) {

}

func Snowflake() {

}

func SnowflakeFlag() {
	snowflake = flag.Bool("snowflake", false, "Disabled on this build")
}
