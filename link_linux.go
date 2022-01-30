package main

import (
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
)

func DesktopDirectory() (string, error) {
	myself, err := user.Current()
	if err != nil {
		return "", err
	}
	homedir := myself.HomeDir
	desktop := filepath.Join(homedir, ".local", "share", "applications")
	return desktop, nil
}

func CreateShortcuts() error {
	desktopDir, err := DesktopDirectory()
	if err != nil {
		return err
	}
	torBrowserPath, err := pathToMe()
	if err != nil {
		return err
	}
	tordesktop := filepath.Join(desktopDir, "torbrowser.desktop")
	torbrowserShortcut := torBrowserPath + " -torbrowser"
	if err := makeLink(torbrowserShortcut, tordesktop); err != nil {
		return err
	}
	i2pbrowserPath := torBrowserPath + " -i2pbrowser"
	i2pdesktop := filepath.Join(desktopDir, "i2ptorbrowser.desktop")
	if err := makeLink(i2pbrowserPath, i2pdesktop); err != nil {
		return err
	}
	return nil
}

func desktopTemplate(command string) string {
	return `[Desktop Entry]
Encoding=UTF-8
Version=1.0
Type=Application
Terminal=false
Exec=/bin/sh -c "` + command + `"
Name=Tor Browser
Categories=Network;WebBrowser;
Icon=/var/lib/i2pbrowser/icons/onion.png
`
}

func makeLink(src, dst string) error {
	log.Println("Creating desktop shortcut:", dst)
	return ioutil.WriteFile(dst, []byte(desktopTemplate(src)), 0644)
}
