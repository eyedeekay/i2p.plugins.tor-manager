//go:build darwin
// +build darwin

package main

import (
	flag "github.com/spf13/pflag"
)

func SnowflakeFlag() {
	snowflake = flag.Bool("snowflake", false, "Disabled in this build")
}

func Snowflake() {
}
