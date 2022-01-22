package tbsupervise

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

var UNPACK_URL = tbget.UNPACK_PATH
var DEFAULT_TB_LANG = tbget.DefaultIETFLang

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

func (s *Supervisor) FirefoxPath() string {
	return filepath.Join(s.UnpackPath, "Browser", "firefox.real")
}

func (s *Supervisor) TBDirectory() string {
	return filepath.Join(s.UnpackPath, "Browser")
}

func (s *Supervisor) TorPath() string {
	return filepath.Join(s.UnpackPath, "Browser", "TorBrowser", "Tor", "tor")
}

func (s *Supervisor) TorDataPath() string {
	return filepath.Join(s.UnpackPath, "Browser", "TorBrowser", "Data")
}

func (s *Supervisor) I2PDataPath() string {
	//if tbget.FileExists(filepath.Join(s.UnpackPath, "i2p.firefox")) {
	return filepath.Join(filepath.Dir(s.UnpackPath), "i2p.firefox")
	//}
}

func (s *Supervisor) RunTBWithLang() error {
	tbget.ARCH = ARCH
	if s.Lang == "" {
		s.Lang = DEFAULT_TB_LANG
	}
	if s.UnpackPath == "" {
		s.UnpackPath = UNPACK_URL
	}

	log.Println("running tor browser with lang", s.Lang, s.UnpackPath)
	switch OS {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running tor browser with lang", s.Lang, s.UnpackPath)
			s.cmd = exec.Command(s.TBPath())
			s.cmd.Stdout = os.Stdout
			s.cmd.Stderr = os.Stderr
			return s.cmd.Run()
		} else {
			log.Println("tor browser not found at", s.TBPath())
			return fmt.Errorf("tor browser not found at %s", s.TBPath())
		}
	case "darwin":
		cmd := exec.Command("/usr/bin/env", "open", "-a", "\"Tor Browser.app\"")
		cmd.Dir = s.TBDirectory()
		return cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", "\""+s.TBDirectory()+"\"", "\"Tor Browser.exe\"")
		cmd.Dir = s.TBDirectory()
		return cmd.Run()
	default:
	}

	return nil
}

func (s *Supervisor) RunI2PBWithLang() error {
	tbget.ARCH = ARCH
	if s.Lang == "" {
		s.Lang = DEFAULT_TB_LANG
	}
	if s.UnpackPath == "" {
		s.UnpackPath = UNPACK_URL
	}

	log.Println("running tor browser with lang", s.Lang, s.UnpackPath)
	switch OS {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running Tor browser with lang and I2P Profile", s.Lang, s.UnpackPath, s.FirefoxPath(), "--profile", s.I2PDataPath())
			s.cmd = exec.Command(s.FirefoxPath(), "--profile", s.I2PDataPath())
			s.cmd.Stdout = os.Stdout
			s.cmd.Stderr = os.Stderr
			return s.cmd.Run()
		} else {
			log.Println("tor browser not found at", s.FirefoxPath())
			return fmt.Errorf("tor browser not found at %s", s.FirefoxPath())
		}
	case "darwin":
		cmd := exec.Command("/usr/bin/env", "open", "-a", "\"Tor Browser.app\"")
		cmd.Dir = s.TBDirectory()
		return cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", "\""+s.TBDirectory()+"\"", "\"Tor Browser.exe\"")
		cmd.Dir = s.TBDirectory()
		return cmd.Run()
	default:
	}

	return nil
}

func (s *Supervisor) RunTorWithLang() error {
	tbget.ARCH = ARCH
	if s.Lang == "" {
		s.Lang = DEFAULT_TB_LANG
	}
	if s.UnpackPath == "" {
		s.UnpackPath = UNPACK_URL
	}

	log.Println("running tor with lang", s.Lang, s.UnpackPath)
	switch OS {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running tor with lang", s.Lang, s.UnpackPath)
			s.cmd = exec.Command(s.TorPath())
			s.cmd.Stdout = os.Stdout
			s.cmd.Stderr = os.Stderr
			return s.cmd.Run()
		} else {
			log.Println("tor not found at", s.TorPath())
			return fmt.Errorf("tor not found at %s", s.TorPath())
		}
	case "darwin":
		cmd := exec.Command("/usr/bin/env", "open", "-a", "\"Tor Browser.app\"")
		cmd.Dir = s.TBDirectory()
		return cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", "\""+s.TBDirectory()+"\"", "\"Tor Browser.exe\"")
		cmd.Dir = s.TBDirectory()
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
