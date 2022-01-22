package main

import (
	"flag"
	"log"

	"github.com/cloudfoundry/jibber_jabber"
	tbserve "i2pgit.org/idk/i2p.plugins.tor-manager/serve"
)

//var runtimePair = tbget.GetRuntimePair()

var (
	lang   = flag.String("lang", "", "Language to download")
	os     = flag.String("os", "linux", "OS/arch to download")
	arch   = flag.String("arch", "64", "OS/arch to download")
	browse = flag.Bool("browse", false, "Open the browser")
	/*mirror   = flag.String("mirror", "", "Mirror to use")*/
	bemirror = flag.Bool("bemirror", false, "Act as an in-I2P mirror when you're done downloading")
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
	client, err := tbserve.NewClient("", *lang, *os, *arch)
	if err != nil {
		log.Fatal("Couldn't create client", err)
	}
	if err := client.Serve(); err != nil {
		log.Fatal(err)
	}
}
