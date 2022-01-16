package main

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
