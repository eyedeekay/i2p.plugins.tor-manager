package main

import (
	"log"
	"os"
	"os/exec"

	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

func mirrorAll() error {
	log.Println("Mirroring all languages, platforms, and architectures")
	for x, l := range tbget.Languages() {
		log.Println("Mirroring language:", l, "(", x, ")", "of", len(tbget.Languages()), "languages complete")
		err := mirrorLang(l)
		if err != nil {
			return err
		}
	}
	return nil
}

func mirrorLang(ietf string) error {
	// get the path to myself(the executable)
	path, err := os.Executable()
	if err != nil {
		return err
	}
	// set the environment variables
	//TOR_MANAGER_CLEARNET_MIRROR=true
	err = os.Setenv("TOR_MANAGER_CLEARNET_MIRROR", "true")
	if err != nil {
		return err
	}
	//TOR_MANAGER_REQUIRE_PASSWORD=false
	err = os.Setenv("TOR_MANAGER_REQUIRE_PASSWORD", "false")
	if err != nil {
		return err
	}
	err = mirrorPlatform(path, ietf, "linux", "64")
	if err != nil {
		return err
	}
	err = mirrorPlatform(path, ietf, "linux", "32")
	if err != nil {
		return err
	}
	err = mirrorPlatform(path, ietf, "win", "64")
	if err != nil {
		return err
	}
	err = mirrorPlatform(path, ietf, "win", "32")
	if err != nil {
		return err
	}
	err = mirrorPlatform(path, ietf, "osx", "64")
	if err != nil {
		return err
	}
	return nil
}

func mirrorPlatform(path, ietf, platform, arch string) error {
	cmd := exec.Command(path, "-nounpack", "-notor", "-os", platform, "-lang="+ietf, "-arch="+arch, "-p2p=false", "-nevertor")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "TOR_MANAGER_CLEARNET_MIRROR=true")
	cmd.Env = append(cmd.Env, "TOR_MANAGER_REQUIRE_PASSWORD=false")
	return cmd.Run()
}
