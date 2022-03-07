package main

import (
	"log"
	"os"
	"path/filepath"

	trayirc "i2pgit.org/idk/libbrb"
)

func BRBClient(directory, server string) {
	brbdirectory := filepath.Join(directory, "brb")
	brbdirectory, err := filepath.Abs(brbdirectory)
	if err != nil {
		log.Println(err)
		return
	}
	os.MkdirAll(brbdirectory, 0755)
	brb, err := trayirc.NewBRBFromOptions(
		trayirc.SetBRBConfigDirectory(brbdirectory),
		trayirc.SetBRBServerName(server),
		trayirc.SetBRBServerConfig("ircd.yml"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = brb.IRC(); err != nil {
		log.Fatal(err)
	}
}
