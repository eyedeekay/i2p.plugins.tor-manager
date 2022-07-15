package main

import (
	"embed"
	"io"
	"log"
	"net/http"
	"os"

	flag "github.com/spf13/pflag"

	"fyne.io/systray"

	"i2pgit.org/idk/blizzard/icon"

	"git.torproject.org/pluggable-transports/snowflake.git/common/safelog"
	sf "git.torproject.org/pluggable-transports/snowflake.git/proxy/lib"
)

//go:embed tor-browser/www/home.css
//go:embed tor-browser/www/index.html
//go:embed tor-browser/www/blizzard.png
var snowflakeContent embed.FS

var snowflakeProxy sf.SnowflakeProxy

var (
	capacity           = flag.Uint("snowflake-capacity", 0, "maximum concurrent clients")
	stunURL            = flag.String("snowflake-stun", sf.DefaultSTUNURL, "broker URL")
	logFilename        = flag.String("snowflake-log", "", "log filename")
	rawBrokerURL       = flag.String("snowflake-broker", sf.DefaultBrokerURL, "broker URL")
	unsafeLogging      = flag.Bool("snowflake-unsafe-logging", false, "prevent logs from being scrubbed")
	keepLocalAddresses = flag.Bool("snowflake-keep-local-addresses", false, "keep local LAN address ICE candidates")
	relayURL           = flag.String("snowflake-relay", sf.DefaultRelayURL, "websocket relay URL")
	snowflakeDirectory = flag.String("snowflake-directory", "", "directory with a page to serve locally for your snowflake. If empty, no local page is served.")
	snowflakePort      = flag.String("snowflake-port", "7676", "port to serve info page(directory) on")
)

func SnowflakeFlag() {
	snowflake = flag.Bool("snowflake", false, "Offer a Snowflake to other Tor Browser users")
}

func Snowflake() {
	snowflakeProxy = sf.SnowflakeProxy{
		Capacity:           uint(*capacity),
		STUNURL:            *stunURL,
		BrokerURL:          *rawBrokerURL,
		KeepLocalAddresses: *keepLocalAddresses,
		RelayURL:           *relayURL,
	}

	var logOutput io.Writer = os.Stderr
	log.SetFlags(log.LstdFlags | log.LUTC)

	log.SetFlags(log.LstdFlags | log.LUTC)
	if *logFilename != "" {
		f, err := os.OpenFile(*logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		logOutput = io.MultiWriter(os.Stderr, f)
	}
	if *unsafeLogging {
		log.SetOutput(logOutput)
	} else {
		log.SetOutput(&safelog.LogScrubber{Output: logOutput})
	}

	go func() {
		if *directory != "" {
			http.Handle("/", http.FileServer(http.Dir(*snowflakeDirectory)))
		} else {
			http.Handle("/", http.FileServer(http.FS(snowflakeContent)))
		}

		log.Printf("Serving %s on HTTP localhost:snowflakePort: %s\n", *snowflakeDirectory, *snowflakePort)
		log.Fatal(http.ListenAndServe("localhost:"+*snowflakePort, nil))
	}()

	err := snowflakeProxy.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func onSnowflakeReady() {
	if !*snowflake {
		return
	}
	mSnowflakeQuit := systray.AddMenuItem("Stop Snowflake", "Close the application and stop your snowflake.")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mSnowflakeQuit.SetIcon(icon.Data)
	runloop := true
	for runloop {
		select {
		case <-mSnowflakeQuit.ClickedCh:
			snowflakeProxy.Stop()
			runloop = false
			log.Println("Snowflake stopped")
		}
	}
}

func onSnowflakeExit() {
	if !*snowflake {
		return
	}
	log.Println("Stopping the Snowflake")
	snowflakeProxy.Stop()
}
