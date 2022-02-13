//go:build !i2pd
// +build !i2pd

package main

func InitI2PSAM() func() {
	return nullFunc
}

func nullFunc() {

}
func StartI2P() {

}

func StopI2P() {

}
