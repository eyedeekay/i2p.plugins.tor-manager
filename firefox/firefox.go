package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cloudfoundry/jibber_jabber"
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
	tbserve "i2pgit.org/idk/i2p.plugins.tor-manager/serve"
)

//go:embed tor-browser/unpack/i2p.firefox/*
//go:embed tor-browser/unpack/i2p.firefox.config/*
//go:embed tor-browser/unpack/awo@eyedeekay.github.io.xpi
//go:embed tor-browser/TPO-signing-key.pub
//go:embed garliconion.png
//go:embed onion.png
//go:embed torbrowser.desktop
//go:embed i2ptorbrowser.desktop

/*
TODO: A "Default" config file which uses hardened Tor Browser for clearnet
(or default-route) browsing.
*/

var content embed.FS

func OS() string {
	switch runtime.GOOS {
	case "darwin":
		return "osx"
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
	case "arm64":
		if OS() == "osx" {
			return "64"
		}
		return ""
	default:
		return "unknown"
	}
}

var (
	lang       = flag.String("lang", "", "Language to download")
	system     = flag.String("os", OS(), "OS/arch to download")
	arch       = flag.String("arch", ARCH(), "OS/arch to download")
	i2pbrowser = flag.Bool("i2pbrowser", false, "Open I2P in Tor Browser")
	i2pconfig  = flag.Bool("i2pconfig", false, "Open I2P routerconsole in Tor Browser with javscript enabled and non-routerconsole sites disabled")
	torbrowser = flag.Bool("torbrowser", false, "Open Tor Browser")
	verbose    = flag.Bool("verbose", false, "Verbose output")
	directory  = flag.String("directory", "", "Directory operate in")
	host       = flag.String("host", "127.0.0.1", "Host to serve on")
	port       = flag.Int("port", 7695, "Port to serve on")
	bemirror   = flag.Bool("bemirror", false, "Act as an in-I2P mirror when you're done downloading")
	shortcuts  = flag.Bool("shortcuts", false, "Create desktop shortcuts")
	apparmor   = flag.Bool("apparmor", false, "Generate apparmor rules")
	offline    = flag.Bool("offline", false, "Work offline. Differs from Firefox's offline mode in that cannot be disabled until the browser is closed.")
	clearnet   = flag.Bool("clearnet", false, "Use clearnet (no Tor or I2P)")
	profile    = flag.String("profile", "", "use a custom profile path, normally blank")
	help       = flag.Bool("help", false, "Print help")
	mirror     = flag.String("mirror", "https://download.mozilla.org/?product=firefox-latest", "Mirror to use")
	/*onion    = flag.Bool("onion", false, "Serve an onion site which shows some I2P propaganda, magnet links, your I2P mirror URL if configured")*/
	/*torrent  = flag.Bool("torrent", false, "Create a torrent of the downloaded files and seed it over I2P using an Open Tracker")*/
	/*ptop     = flag.Bool("p2p", false, "Use bittorrent over I2P to download the initial copy of Tor Browser")*/

)

var client *tbserve.Client

func main() {
	filename := filepath.Base(os.Args[0])
	usage := flag.Usage
	flag.Usage = func() {
		fmt.Printf("Usage: %s %s\n", filename, "[options]")
		fmt.Printf("\n")
		fmt.Printf("Downloads, verifies and unpacks Tor Browser. Manages the Tor Browser\n")
		fmt.Printf("system in environments where Tor is not in use.\n")
		fmt.Printf("\n")
		fmt.Printf("Options:\n")
		fmt.Printf("\n")
		usage()
	}
	flag.Parse()
	if *lang == "" {
		var err error
		*lang, err = jibber_jabber.DetectIETF()
		if err != nil {
			log.Fatal("Please specify a language", err)
		}
		log.Println("Using auto-detected language", *lang)
	}
	tbget.WORKING_DIR = *directory
	var err error
	if client, err = tbserve.NewFirefoxClient(*verbose, *lang, *system, *arch, *mirror, &content); err != nil {
		log.Fatal(err)
	}
}
