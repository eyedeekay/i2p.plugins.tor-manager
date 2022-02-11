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
	i2cpcheck "github.com/eyedeekay/checki2cp"
	"github.com/itchio/damage"
	"github.com/itchio/damage/hdiutil"
	"github.com/itchio/headway/state"
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
	tbserve "i2pgit.org/idk/i2p.plugins.tor-manager/serve"
)

/*
TODO: A "Default" config file which uses hardened Tor Browser for clearnet
(or default-route) browsing.
*/

//go:embed tor-browser/unpack/i2p.firefox/*
//go:embed tor-browser/unpack/i2p.firefox.config/*
//go:embed tor-browser/unpack/awo@eyedeekay.github.io.xpi
//go:embed tor-browser/TPO-signing-key.pub
//go:embed garliconion.png
//go:embed onion.png
//go:embed torbrowser.desktop
//go:embed i2ptorbrowser.desktop
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
	profile    = flag.String("profile", "", "use a custom profile path, normally blank")
	help       = flag.Bool("help", false, "Print help")
	/*onion    = flag.Bool("onion", false, "Serve an onion site which shows some I2P propaganda, magnet links, your I2P mirror URL if configured")*/
	/*torrent  = flag.Bool("torrent", false, "Create a torrent of the downloaded files and seed it over I2P using an Open Tracker")*/
	/*ptop     = flag.Bool("p2p", false, "Use bittorrent over I2P to download the initial copy of Tor Browser")*/
	/*mirror   = flag.String("mirror", "", "Mirror to use")*/
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
	tbget.WORKING_DIR = *directory
	if filename == "i2pbrowser" {
		log.Println("Starting I2P in Tor Browser")
		*i2pbrowser = true
	} else if filename == "torbrowser" {
		log.Println("Starting Tor Browser")
		*torbrowser = true
	} else if filename == "i2pconfig" {
		log.Println("Starting I2P routerconsole in Tor Browser")
		*i2pconfig = true
	} else if filename == "firefox" {
		log.Println("Starting Firefox")
		if *profile != "" {
			*profile = filepath.Join(tbget.WORKING_DIR, "profile.firefox")
		}
		log.Println("Using profile", *profile)
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
	var err error
	client, err = tbserve.NewClient(*verbose, *lang, *system, *arch, &content)
	if err != nil {
		log.Fatal("Couldn't create client", err)
	}
	if *i2pbrowser || *i2pconfig {
		if tbget.TestHTTPDefaultProxy() {
			log.Println("I2P HTTP proxy OK")
		} else {
			log.Println("I2P HTTP proxy not OK")
			run, err := i2cpcheck.ConditionallyLaunchI2P()
			if err != nil {
				log.Fatal("Couldn't launch I2P", err)
			}
			if run {
				if tbget.TestHTTPDefaultProxy() {
					log.Println("I2P HTTP proxy OK after launching I2P")
				} else {
					go proxy()
					if !tbget.TestHTTPBackupProxy() {
						log.Fatal("Please set the I2P HTTP proxy on localhost:4444", err)
					}
				}
			} else {
				log.Fatal("Failed to run I2P", err)
				//TODO: Link libi2pd and start our own router if we cant find one anywhere.
				//TODO: loop again until TestHTTPDefaultProxy is up
			}
		}
	}
	if *apparmor {
		err := GenerateAppArmor()
		if err != nil {
			log.Fatal("Couldn't generate apparmor rules", err)
		}
		log.Println("################################################################")
		log.Println("#             AppArmor rules generated successfully            #")
		log.Println("################################################################")
		log.Println("!IMPORTANT! You must now run the following commands:")
		log.Println("sudo mkdir -p /etc/apparmor.d/tunables/")
		log.Println("sudo cp tunables.torbrowser.apparmor /etc/apparmor.d/tunables/torbrowser")
		log.Println("sudo cp torbrowser.Tor.tor.apparmor /etc/apparmor.d/torbrowser.Tor.tor")
		log.Println("sudo cp torbrowser.Browser.firefox.apparmor /etc/apparmor.d/torbrowser.Browser.firefox")
		log.Println("sudo apparmor_parser -r /etc/apparmor.d/tunables/torbrowser")
		log.Println("sudo apparmor_parser -r /etc/apparmor.d/torbrowser.Tor.tor")
		log.Println("sudo apparmor_parser -r /etc/apparmor.d/torbrowser.Browser.firefox")
		log.Println("To copy them to apparmor profiles directory and reload AppArmor")
		return
	}
	if *shortcuts {
		err := CreateShortcuts()
		if err != nil {
			log.Fatal("Couldn't create desktop shortcuts", err)
		}
	}
	client.Host = *host
	client.Port = *port
	client.TBS.Profile = &content
	client.TBS.PassThroughArgs = flag.Args()
	consumer := &state.Consumer{
		OnMessage: func(lvl string, msg string) {
			log.Printf("[%s] %s", lvl, msg)
		},
	}
	host := hdiutil.NewHost(consumer)
	defer damage.Unmount(host, client.TBD.BrowserDir())
	//	log.Fatalf("%s", client.TBS.PassThroughArgs)
	if *help {
		flag.Usage()
		if err := client.TBS.RunTBHelpWithLang(); err != nil {
			log.Fatal(err)
		}
		return
	}
	if *profile != "" && !*offline {
		log.Println("Using a custom profile")
		if err := client.TBS.RunTBBWithProfile(*profile); err != nil {
			log.Fatal(err)
		}
	} else if *offline {
		if *profile == "" {
			*profile = "firefox.offline"
		}
		log.Println("Working offline")

		if err := client.TBS.RunTBBWithOfflineProfile(*profile, *offline); err != nil {
			log.Fatal(err)
		}
	} else if *i2pbrowser {
		if err := client.TBS.RunI2PBWithLang(); err != nil {
			log.Fatal(err)
		}
	} else if *i2pconfig {
		if err := client.TBS.RunI2PBAppWithLang(); err != nil {
			log.Fatal(err)
		}
	} else if *torbrowser {
		if err := client.TBS.RunTBWithLang(); err != nil {
			log.Fatal(err)
		}
	} else {
		if *bemirror {
			go client.TBD.Serve()
		}
		if err := client.Serve(); err != nil {
			log.Fatal(err)
		}
	}
}

func pathToMe() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath, err := filepath.Abs(ex)
	if err != nil {
		return "", err
	}
	return exPath, nil
}
