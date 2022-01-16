package main

import (
	"flag"
	"log"

	"github.com/cloudfoundry/jibber_jabber"
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

var (
	lang = flag.String("lang", "", "Language to download")
)

func main() {
	flag.Parse()
	if *lang == "" {
		var err error
		*lang, err = jibber_jabber.DetectIETF()
		if err != nil {
			log.Fatal("Please specify a language", err)
		}
		log.Println("Using auto-detected language", *lang)
	}
	bin, sig, err := tbget.DownloadUpdaterForLang(*lang)
	if err != nil {
		panic(err)
	}
	if err := tbget.CheckSignature(bin, sig); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Signature check passed")
	}

}
