package main

import (
	"os"
	"path/filepath"
)

func PluginStat() bool {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	exDir := filepath.Base(exPath)
	return exDir == "plugins"
}

func DefaultDir() string {
	if PluginStat() {
		return ""
	}
	return "tmp-i2pbrowser"
}
