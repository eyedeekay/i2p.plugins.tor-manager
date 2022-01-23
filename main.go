package main

import (
	"embed"
	"flag"
	"log"

	"github.com/cloudfoundry/jibber_jabber"
	tbserve "i2pgit.org/idk/i2p.plugins.tor-manager/serve"
)

//go:embed tor-browser/unpack/i2p.firefox/*
//go:embed tor-browser/TPO-signing-key.pub
var content embed.FS

//var runtimePair = tbget.GetRuntimePair()

var (
	lang       = flag.String("lang", "", "Language to download")
	os         = flag.String("os", "linux", "OS/arch to download")
	arch       = flag.String("arch", "64", "OS/arch to download")
	i2pbrowser = flag.Bool("i2pbrowser", false, "Open I2P in Tor Browser")
	torbrowser = flag.Bool("torbrowser", false, "Open Tor Browser")
	/*mirror   = flag.String("mirror", "", "Mirror to use")*/
	/*bemirror = flag.Bool("bemirror", false, "Act as an in-I2P mirror when you're done downloading")*/
)

func main() {
	flag.Parse()
	if *i2pbrowser == true && *torbrowser == true {
		log.Fatal("Please don't open I2P and Tor Browser at the same time when running from the terminal.")
	}
	if *lang == "" {
		var err error
		*lang, err = jibber_jabber.DetectIETF()
		if err != nil {
			log.Fatal("Please specify a language", err)
		}
		log.Println("Using auto-detected language", *lang)
	}
	client, err := tbserve.NewClient("", *lang, *os, *arch, &content)
	if err != nil {
		log.Fatal("Couldn't create client", err)
	}
	//client.TBD.Profile = &content
	client.TBS.Profile = &content
	if *i2pbrowser {
		client.TBS.RunI2PBWithLang()
	} else if *torbrowser {
		client.TBS.RunTBWithLang()
	} else {
		if err := client.Serve(); err != nil {
			log.Fatal(err)
		}
	}
}
