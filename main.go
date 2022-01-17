package main

import (
	"flag"
	"log"

	"github.com/cloudfoundry/jibber_jabber"
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
	tbsupervise "i2pgit.org/idk/i2p.plugins.tor-manager/supervise"
)

var runtimePair = tbget.GetRuntimePair()

var (
	lang = flag.String("lang", "", "Language to download")
	os   = flag.String("os", tbget.OS, "OS/arch to download")
	arch = flag.String("arch", tbget.ARCH, "OS/arch to download")
	/*mirror   = flag.String("mirror", "", "Mirror to use")*/
	bemirror = flag.Bool("bemirror", false, "Act as an in-I2P mirror when you're done downloading")
)

func main() {
	flag.Parse()
	tbget.OS = *os
	tbget.ARCH = *arch
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
	if err := tbsupervise.RunTBWithLang(*lang); err != nil {
		log.Fatal(err)
	}
}
