package tbsupervise

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

var wd, _ = os.Getwd()

var UNPACK_URL = filepath.Join(wd, "unpack")
var DEFAULT_TB_LANG = tbget.DefaultIETFLang
var DEFAULT_TB_DIRECTORY = filepath.Join(UNPACK_URL, "tor-browser"+"_"+DEFAULT_TB_LANG)
var DEFAULT_TB_PATH = filepath.Join(DEFAULT_TB_DIRECTORY, "Browser", "start-tor-browser")

var (
	OS   = tbget.OS
	ARCH = tbget.ARCH
)

func RunTBWithLang(lang string) error {
	tbget.ARCH = ARCH
	if lang == "" {
		lang = DEFAULT_TB_LANG
	}
	log.Println("running tor browser with lang", lang, DEFAULT_TB_PATH)
	switch OS {
	case "linux":
		if tbget.FileExists(DEFAULT_TB_PATH) {
			log.Println("running tor browser with lang", lang, DEFAULT_TB_PATH)
			cmd := exec.Command(DEFAULT_TB_PATH)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		} else {
			log.Println("tor browser not found at", DEFAULT_TB_PATH)
			return fmt.Errorf("tor browser not found at %s", DEFAULT_TB_PATH)
		}
	case "darwin":
		cmd := exec.Command("/usr/bin/env", "open", "-a", "\"Tor Browser.app\"")
		cmd.Dir = DEFAULT_TB_DIRECTORY
		return cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", "\""+DEFAULT_TB_PATH+"\"", "\"Tor Browser.exe\"")
		cmd.Dir = DEFAULT_TB_DIRECTORY
		return cmd.Run()
	default:
	}

	return nil
}
