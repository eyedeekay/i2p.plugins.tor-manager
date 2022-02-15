//go:build !noi2pd
// +build !noi2pd

package main

import (
	i2pd "github.com/eyedeekay/go-i2pd/goi2pd"
)

func InitI2PSAM() func() {
	return i2pd.InitI2PSAM()
}

func StartI2P() {
	i2pd.StartI2P()
}

func StopI2P() {
	defer i2pd.StopI2P()
}
