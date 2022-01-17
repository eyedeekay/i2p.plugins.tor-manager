package tbsupervise

import (
	"os"
	"os/exec"
	"path/filepath"

	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

var wd, _ = os.Getwd()

var UNPACK_URL = filepath.Join(wd, "unpack")
var DEFAULT_TB_LANG = tbget.DefaultIETFLang
var DEFAULT_TB_PATH = filepath.Join(UNPACK_URL, "tor-browser"+"_"+DEFAULT_TB_LANG)

var (
	OS   = tbget.OS
	ARCH = tbget.ARCH
)

func RunTBWithLang(lang string) error {
	tbget.ARCH = ARCH
	if lang == "" {
		lang = DEFAULT_TB_LANG
	}
	switch OS {
	case "linux":
		cmd := exec.Command("/usr/bin/env", "sh", "-c", "./tor-browser_"+lang+"/Browser/start-tor-browser")
		cmd.Dir = DEFAULT_TB_PATH
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("/usr/bin/env", "open", "-a", "Tor Browser.app")
		cmd.Dir = DEFAULT_TB_PATH
		return cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", "Tor Browser.exe")
		cmd.Dir = DEFAULT_TB_PATH
		return cmd.Run()
	default:
	}

	return nil
}
