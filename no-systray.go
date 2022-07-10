//go:build nosystray
// +build nosystray

package main

import (
	flag "github.com/spf13/pflag"
)

var shutdown = false

func runSysTray(down bool) {

}

func Snowflake() {

}

func SnowflakeFlag() {
	snowflake = flag.Bool("snowflake", false, "Disabled on this build")
}
