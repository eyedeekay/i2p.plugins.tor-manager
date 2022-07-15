package main

import (
	"fmt"
	"os"
)

var VERSION string = "0.0.15"

func printversion() {
	fmt.Fprintf(os.Stdout, "Version: ", VERSION)
}
