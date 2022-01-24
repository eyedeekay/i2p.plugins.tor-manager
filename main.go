package main

import (
	"embed"
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cloudfoundry/jibber_jabber"
	tbserve "i2pgit.org/idk/i2p.plugins.tor-manager/serve"
)

//go:embed tor-browser/unpack/i2p.firefox/*
//go:embed tor-browser/TPO-signing-key.pub
var content embed.FS

func OS() string {
	switch runtime.GOOS {
	case "darwin":
		return "mac"
	case "linux":
		return "linux"
	case "windows":
		return "win"
	default:
		return "unknown"
	}
}

func ARCH() string {
	switch runtime.GOARCH {
	case "386":
		return "32"
	case "amd64":
		return "64"
	default:
		return "unknown"
	}
}

var (
	lang       = flag.String("lang", "", "Language to download")
	system     = flag.String("os", OS(), "OS/arch to download")
	arch       = flag.String("arch", ARCH(), "OS/arch to download")
	i2pbrowser = flag.Bool("i2pbrowser", false, "Open I2P in Tor Browser")
	torbrowser = flag.Bool("torbrowser", false, "Open Tor Browser")
	/*mirror   = flag.String("mirror", "", "Mirror to use")*/
	/*bemirror = flag.Bool("bemirror", false, "Act as an in-I2P mirror when you're done downloading")*/
)

func main() {
	filename := filepath.Base(os.Args[0])
	flag.Parse()
	if filename == "i2pbrowser" {
		log.Println("Starting I2P in Tor Browser")
		*i2pbrowser = true
	} else if filename == "torbrowser" {
		log.Println("Starting Tor Browser")
		*torbrowser = true
	}
	if *i2pbrowser && *torbrowser {
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
	client, err := tbserve.NewClient("", *lang, *system, *arch, &content)
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
