//go:build !nosystray
// +build !nosystray

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eyedeekay/go-i2pcontrol"
	"github.com/getlantern/systray"
	"github.com/ncruces/zenity"
	"i2pgit.org/idk/i2p.plugins.tor-manager/icon"
)

var running = false
var shutdown = false

func Password() string {
	require_password := os.Getenv("TOR_MANAGER_REQUIRE_PASSWORD")
	if !PluginStat() {
		require_password = "true"
	}
	if require_password == "true" || require_password == "1" {
		passwd, err := zenity.Entry(
			"Enter a password if you want to encrypt the working directory",
			zenity.Title("Work Directory Encryption"),
			zenity.CancelLabel("Don't encrypt"),
			zenity.OKLabel("Encrypt"),
			zenity.Width(400),
			zenity.EntryText("password"),
		)
		if err != nil {
			if !strings.Contains(err.Error(), "canceled") {
				log.Panicln(err)
			}
			log.Println("Password dialog canceled")
			return ""
		}
		return passwd
	}
	return ""
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Tor Manager for I2P")
	systray.SetTooltip("Tor and I2P integrated")

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("Tor Managing I2P Browser Plugin")
		systray.SetTooltip("Configures I2P and Tor in the same browser")
		mEnabled := systray.AddMenuItem("Online", "I2P and Tor are both running")
		// Sets the icon of a menu item. Only available on Mac.
		mEnabled.SetTemplateIcon(icon.Data, icon.Data)

		subMenuTop := systray.AddMenuItem("Launch a Browser", "Launch a browser")
		subMenuBottom := subMenuTop.AddSubMenuItem("Launch Tor Browser configured for I2P", "Modify and launch the Tor Browser Bundle for I2P")
		subMenuBottom2 := subMenuTop.AddSubMenuItem("Launch the Tor Browser", "Launch the standard Tor Browser bundle")
		subMenuBottom3 := subMenuTop.AddSubMenuItem("Launch Hardened Firefox in Clearnet Mode", "Launch the Tor Browser bundle, but without Tor")
		systray.AddSeparator()
		go onSnowflakeReady()
		systray.AddSeparator()

		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		// Sets the icon of a menu item. Only available on Mac.
		mQuit.SetIcon(icon.Data)

		for {
			select {
			case <-mEnabled.ClickedCh:
				mEnabled.SetTitle("I2P and Tor are both running")
				time.Sleep(time.Second * 3)
				mEnabled.SetTitle("Online")
			case <-subMenuBottom.ClickedCh:
				fmt.Println("Launching Tor Browser configured for I2P")
				if err := client.TBS.RunI2PBWithLang(); err != nil {
					log.Println(err)
				}
			case <-subMenuBottom2.ClickedCh:
				fmt.Println("Launching the Tor Browser")
				if err := client.TBS.RunTBWithLang(); err != nil {
					log.Println(err)
				}
			case <-subMenuBottom3.ClickedCh:
				fmt.Println("Launching Hardened Firefox in Clearnet Mode")
				if err := client.TBS.RunTBBWithOfflineClearnetProfile(*profile, false, true); err != nil {
					log.Println(err)
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit now...")
				return
			}

			time.Sleep(time.Second * 1)
		}
	}()
}

func onExit() {
	if *snowflake {
		snowflakeProxy.Stop()
	}
	if *password != "" {
		log.Println("Encrypting directory with password")
		os.Remove(*directory + ".tar.xz")
		EncryptTarXZip(*directory, *password)
	}
	if shutdown {
		i2pcontrol.Initialize("127.0.0.1", "7657", "")
		_, err := i2pcontrol.Authenticate("itoopie")
		if err != nil {
			log.Println(err)
		}
		message, err := i2pcontrol.ShutdownGraceful()
		if err != nil {
			log.Println(err)
		}
		ltc := 0
		for {
			tunnels, err := i2pcontrol.ParticipatingTunnels()
			if err != nil {
				log.Println(err)
			}
			if ltc != tunnels {
				log.Println("Participating tunnels:", tunnels)
			}
			ltc = tunnels
			if tunnels <= 0 {
				break
			}
		}
		log.Println(message)
	}
	running = false
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func runSysTray(down bool) {
	if !running {
		running = true
		shutdown = down
		systray.Run(onReady, onExit)
	}
}
