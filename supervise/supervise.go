package tbsupervise

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mitchellh/go-ps"
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
	torcmd     *exec.Cmd
	tbcmd      *exec.Cmd
	ibcmd      *exec.Cmd
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

func (s *Supervisor) tbbail() error {
	if s.tbcmd != nil && s.tbcmd.Process != nil && s.tbcmd.ProcessState != nil {
		if s.tbcmd.ProcessState.Exited() {
			return nil
		}
		return fmt.Errorf("Already running")
	}
	return nil
}

func (s *Supervisor) RunTBWithLang() error {
	tbget.ARCH = ARCH
	if s.Lang == "" {
		s.Lang = DEFAULT_TB_LANG
	}
	if s.UnpackPath == "" {
		s.UnpackPath = UNPACK_URL
	}

	if s.tbbail() != nil {
		return nil
	}

	log.Println("running tor browser with lang", s.Lang, s.UnpackPath)
	switch OS {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running tor browser with lang", s.Lang, s.UnpackPath)
			s.tbcmd = exec.Command(s.TBPath())
			s.tbcmd.Stdout = os.Stdout
			s.tbcmd.Stderr = os.Stderr
			return s.tbcmd.Run()
		} else {
			log.Println("tor browser not found at", s.TBPath())
			return fmt.Errorf("tor browser not found at %s", s.TBPath())
		}
	case "darwin":
		s.tbcmd = exec.Command("/usr/bin/env", "open", "-a", "\"Tor Browser.app\"")
		s.tbcmd.Dir = s.TBDirectory()
		return s.tbcmd.Run()
	case "windows":
		s.tbcmd = exec.Command("cmd", "/c", "start", "\""+s.TBDirectory()+"\"", "\"Tor Browser.exe\"")
		s.tbcmd.Dir = s.TBDirectory()
		return s.tbcmd.Run()
	default:
	}

	return nil
}

func (s *Supervisor) ibbail() error {
	if s.ibcmd != nil && s.ibcmd.Process != nil && s.ibcmd.ProcessState != nil {
		if s.ibcmd.ProcessState.Exited() {
			return nil
		}
		return fmt.Errorf("Already running")
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

	if s.ibbail() != nil {
		return nil
	}

	log.Println("running tor browser with lang", s.Lang, s.UnpackPath)
	switch OS {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running Tor browser with lang and I2P Profile", s.Lang, s.UnpackPath, s.FirefoxPath(), "--profile", s.I2PDataPath())
			s.ibcmd = exec.Command(s.FirefoxPath(), "--profile", s.I2PDataPath())
			s.ibcmd.Stdout = os.Stdout
			s.ibcmd.Stderr = os.Stderr
			return s.ibcmd.Run()
		} else {
			log.Println("tor browser not found at", s.FirefoxPath())
			return fmt.Errorf("tor browser not found at %s", s.FirefoxPath())
		}
	case "darwin":
		s.ibcmd = exec.Command("/usr/bin/env", "open", "-a", "\"Tor Browser.app\"")
		s.ibcmd.Dir = s.TBDirectory()
		return s.ibcmd.Run()
	case "windows":
		s.ibcmd = exec.Command("cmd", "/c", "start", "\""+s.TBDirectory()+"\"", "\"Tor Browser.exe\"")
		s.ibcmd.Dir = s.TBDirectory()
		return s.ibcmd.Run()
	default:
	}

	return nil
}

func (s *Supervisor) torbail() error {
	_, err := net.Listen("TCP", "127.0.0.1:9050")
	if err != nil {
		return fmt.Errorf("Already running")
	}
	if s.torcmd != nil && s.torcmd.Process != nil && s.torcmd.ProcessState != nil {
		if s.torcmd.ProcessState.Exited() {
			return nil
		}
		return fmt.Errorf("Already running")
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
	if err := s.torbail(); err != nil {
		return nil
	}

	log.Println("running tor with lang", s.Lang, s.UnpackPath)
	switch OS {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running tor with lang", s.Lang, s.UnpackPath)
			s.torcmd = exec.Command(s.TorPath())
			s.torcmd.Stdout = os.Stdout
			s.torcmd.Stderr = os.Stderr
			return s.torcmd.Run()
		} else {
			log.Println("tor not found at", s.TorPath())
			return fmt.Errorf("tor not found at %s", s.TorPath())
		}
	case "darwin":
		s.torcmd = exec.Command("/usr/bin/env", "open", "-a", "\"Tor Browser.app\"")
		s.torcmd.Dir = s.TBDirectory()
		return s.torcmd.Run()
	case "windows":
		s.torcmd = exec.Command("cmd", "/c", "start", "\""+s.TBDirectory()+"\"", "\"Tor Browser.exe\"")
		s.torcmd.Dir = s.TBDirectory()
		return s.torcmd.Run()
	default:
	}

	return nil
}

func (s *Supervisor) StopTor() error {
	return s.torcmd.Process.Kill()
}

func (s *Supervisor) TorIsAlive() (bool, bool) {
	_, err := net.Listen("TCP", "127.0.0.1:9050")
	if err != nil {
		return true, false
	}
	if s.torcmd != nil && s.torcmd.Process != nil && s.torcmd.ProcessState != nil {
		return !s.torcmd.ProcessState.Exited(), true
	}
	processes, err := ps.Processes()
	if err != nil {
		return false, true
	}
	for _, p := range processes {
		if p.Executable() == s.TorPath() {
			var err error
			s.torcmd.Process, err = os.FindProcess(p.Pid())
			if err == nil {
				return true, true
			}
		}
	}
	return false, true
}

func NewSupervisor(tbPath, lang string) *Supervisor {
	return &Supervisor{
		UnpackPath: tbPath,
		Lang:       lang,
	}
}
