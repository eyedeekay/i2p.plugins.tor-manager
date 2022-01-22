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

var UNPACK_URL = tbget.UNPACK_URL

//var DEFAULT_TB_LANG = tbget.DefaultIETFLang
//var DEFAULT_TB_DIRECTORY = filepath.Join(UNPACK_URL, "tor-browser"+"_"+DEFAULT_TB_LANG)
//var DEFAULT_TB_PATH = filepath.Join(DEFAULT_TB_DIRECTORY, "Browser", "start-tor-browser")

var (
	OS   = tbget.OS
	ARCH = tbget.ARCH
)

type Supervisor struct {
	UnpackPath string
	Lang       string
	cmd        *exec.Cmd
}

func (s *Supervisor) TBPath() string {
	return filepath.Join(s.UnpackPath, "Browser", "start-tor-browser")
}

func (s *Supervisor) TorPath() string {
	return filepath.Join(s.UnpackPath, "Browser", "TorBrowser", "tor")
}

func (s *Supervisor) TorDataPath() string {
	return filepath.Join(s.UnpackPath, "Browser", "TorBrowser", "Data")
}

func (s *Supervisor) RunTBWithLang() error {
	tbget.ARCH = ARCH
	if s.Lang == "" {
		s.Lang = DEFAULT_TB_LANG
	}
	if s.UnpackPath == "" {
		s.UnpackPath = DEFAULT_TB_PATH
	}

	log.Println("running tor browser with lang", s.Lang, s.UnpackPath)
	switch OS {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running tor browser with lang", s.Lang, s.UnpackPath)
			s.cmd = exec.Command(s.UnpackPath)
			s.cmd.Stdout = os.Stdout
			s.cmd.Stderr = os.Stderr
			return s.cmd.Run()
		} else {
			log.Println("tor browser not found at", s.UnpackPath)
			return fmt.Errorf("tor browser not found at %s", s.UnpackPath)
		}
	case "darwin":
		cmd := exec.Command("/usr/bin/env", "open", "-a", "\"Tor Browser.app\"")
		cmd.Dir = DEFAULT_TB_DIRECTORY
		return cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", "\""+s.UnpackPath+"\"", "\"Tor Browser.exe\"")
		cmd.Dir = DEFAULT_TB_DIRECTORY
		return cmd.Run()
	default:
	}

	return nil
}

func NewSupervisor(tbPath, lang string) *Supervisor {
	return &Supervisor{
		UnpackPath: tbPath,
		Lang:       lang,
	}
}
