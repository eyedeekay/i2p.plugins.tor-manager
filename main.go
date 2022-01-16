package main

import (
	"flag"
	"log"

	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

var (
	lang = flag.String("lang", "", "Language to download")
)

func main() {
	flag.Parse()
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
