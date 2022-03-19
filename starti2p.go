package main

import (
	"log"

	i2cpcheck "github.com/eyedeekay/checki2cp"
	"github.com/eyedeekay/go-I2P-jpackage"
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

func StartI2P(directory string) error {
	if tbget.TestHTTPDefaultProxy() {
		log.Println("I2P HTTP proxy OK")
	} else {
		log.Println("I2P HTTP proxy not OK")
		run, err := i2cpcheck.ConditionallyLaunchI2P()
		if err != nil {
			log.Println("Couldn't launch I2P", err)
		}
		if run {
			if tbget.TestHTTPDefaultProxy() {
				log.Println("I2P HTTP proxy OK after launching I2P")
			} else {
				go proxy()
				if !tbget.TestHTTPBackupProxy() {
					log.Println("Please set the I2P HTTP proxy on localhost:4444", err)
					return err
				}
			}
		} else {
			I2Pdaemon, err := I2P.NewDaemon(directory, false)
			if err != nil {
				log.Println(err)
				return err
			}
			if err = I2Pdaemon.Start(); err != nil {
				log.Println(err)
				return err
			}
			shutdown = true
			defer I2Pdaemon.Stop()
			go runSysTray(true)
			if tbget.TestHTTPDefaultProxy() {
				log.Println("I2P HTTP proxy OK")
			} else {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}
