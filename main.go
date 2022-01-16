package main

import tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"

func main() {
	bin, sig, err := tbget.DownloadUpdaterForLang("")
	if err != nil {
		panic(err)
	}
	if tbget.CheckSignature(bin, sig) {
	} else {
		panic("Signature check failed")
	}

}
