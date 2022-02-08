package tbsupervise

import (
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-ps"
	"github.com/otiai10/copy"
	cp "github.com/otiai10/copy"
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

// UNPACK_URL is the URL to place to unpack the Browser Bundle
var UNPACK_URL = tbget.UNPACK_PATH

// DEFAULT_TB_LANG is the default language to use for the Tor Browser Bundle
var DEFAULT_TB_LANG = tbget.DefaultIETFLang

// OS returns the current OS
func OS() string {
	return tbget.OS
}

// ARCH returns the current architecture
func ARCH() string {
	return tbget.ARCH
}

// Supervisor is the main struct for the Tor Browser Bundle Supervisor
type Supervisor struct {
	UnpackPath      string
	Lang            string
	torcmd          *exec.Cmd
	tbcmd           *exec.Cmd
	ibcmd           *exec.Cmd
	Profile         *embed.FS
	PassThroughArgs []string
}

// PTAS is the validator for the pass-through arguments
func (s *Supervisor) PTAS() []string {
	// loop over the arguments and make sure that we remove any --profile, -P args
	// and blank them out.
	var args []string
	for index, arg := range s.PassThroughArgs {
		if arg == "--profile" || arg == "-P" || arg == "-profile" {
			continue
		}
		if index > 0 {
			if s.PassThroughArgs[index-1] == "--profile" || s.PassThroughArgs[index-1] == "-P" || s.PassThroughArgs[index-1] == "-profile" {
				continue
			}
		}
		args = append(args, arg)
	}
	return args
}

// TBPath returns the path to the Tor Browser Bundle launcher
func (s *Supervisor) TBPath() string {
	return filepath.Join(s.UnpackPath, "Browser", "start-tor-browser")
}

// FirefoxPath returns the path to the Firefox executable inside Tor Browser
func (s *Supervisor) FirefoxPath() string {
	switch OS() {
	case "linux":
		return filepath.Join(s.UnpackPath, "Browser", "firefox.real")
	case "windows":
		return filepath.Join(s.UnpackPath, "Browser", "firefox.exe")
	default:
		return filepath.Join(s.UnpackPath, "Browser", "firefox")
	}
}

// TBDirectory returns the path to the Tor Browser Bundle directory
func (s *Supervisor) TBDirectory() string {
	return filepath.Join(s.UnpackPath, "Browser")
}

// TorPath returns the path to the Tor executable
func (s *Supervisor) TorPath() string {
	return filepath.Join(s.UnpackPath, "Browser", "TorBrowser", "Tor", "tor")
}

// TorDataPath returns the path to the Tor Browser Bundle Data directory
func (s *Supervisor) TorDataPath() string {
	return filepath.Join(s.UnpackPath, "Browser", "TorBrowser", "Data")
}

// I2PProfilePath returns the path to the I2P profile
func (s *Supervisor) I2PProfilePath() string {
	fp := filepath.Join(filepath.Dir(s.UnpackPath), ".i2p.firefox")
	if !tbget.FileExists(fp) {
		log.Printf("i2p data not found at %s, unpacking", fp)
		if s.Profile != nil {
			if err := s.UnpackI2PData(); err != nil {
				log.Fatal(err)
			}
		}
	}
	return fp
}

// I2PProfilePath returns the path to the I2P profile
func (s *Supervisor) I2PAppProfilePath() string {
	fp := filepath.Join(filepath.Dir(s.UnpackPath), ".i2p.firefox.config")
	if !tbget.FileExists(fp) {
		log.Printf("i2p app data not found at %s, unpacking", fp)
		if s.Profile != nil {
			if err := s.UnpackI2PAppData(); err != nil {
				log.Fatal(err)
			}
		}
	}
	return fp
}

// I2PDataPath returns the path to the I2P data directory
func (s *Supervisor) I2PDataPath() string {
	fp := s.I2PProfilePath()
	up := filepath.Join(filepath.Dir(s.UnpackPath), "i2p.firefox")
	if tbget.FileExists(up) {
		return up
	}
	log.Printf("i2p workdir not found at %s, copying", up)
	if s.Profile != nil {
		if err := cp.Copy(fp, up); err != nil {
			log.Fatal(err)
		}
	}
	return up
}

// UnpackI2PData unpacks the I2P data into the s.UnpackPath
func (s *Supervisor) UnpackI2PData() error {
	return fs.WalkDir(s.Profile, ".", func(embedpath string, d fs.DirEntry, err error) error {
		fp := filepath.Join(filepath.Dir(s.UnpackPath), ".i2p.firefox")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(embedpath, filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox", "", -1)))
		if d.IsDir() {
			os.MkdirAll(filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox", "", -1)), 0755)
		} else {
			fullpath := path.Join(embedpath)
			bytes, err := s.Profile.ReadFile(fullpath)
			if err != nil {
				return err
			}
			unpack := filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox", "", -1))
			if err := ioutil.WriteFile(unpack, bytes, 0644); err != nil {
				return err
			}
		}
		return nil
	})
}

// I2PAppDataPath returns the path to the I2P application data directory
func (s *Supervisor) I2PAppDataPath() string {
	fp := s.I2PAppProfilePath()
	up := filepath.Join(filepath.Dir(s.UnpackPath), "i2p.firefox.config")
	if tbget.FileExists(up) {
		return up
	}
	log.Printf("i2p workdir not found at %s, copying", up)
	if s.Profile != nil {
		if err := cp.Copy(fp, up); err != nil {
			log.Fatal(err)
		}
	}
	return up
}

// UnpackI2PAppData unpacks the I2P application data into the s.UnpackPath
func (s *Supervisor) UnpackI2PAppData() error {
	return fs.WalkDir(s.Profile, ".", func(embedpath string, d fs.DirEntry, err error) error {
		fp := filepath.Join(filepath.Dir(s.UnpackPath), ".i2p.firefox.config")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(embedpath, filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox.config", "", -1)))
		if d.IsDir() {
			os.MkdirAll(filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox.config", "", -1)), 0755)
		} else {
			fullpath := path.Join(embedpath)
			bytes, err := s.Profile.ReadFile(fullpath)
			if err != nil {
				return err
			}
			unpack := filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox.config", "", -1))
			if err := ioutil.WriteFile(unpack, bytes, 0644); err != nil {
				return err
			}
		}
		return nil
	})
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

// RunTBWithLang runs the Tor Browser with the given language
func (s *Supervisor) RunTBWithLang() error {
	tbget.ARCH = ARCH()
	if s.Lang == "" {
		s.Lang = DEFAULT_TB_LANG
	}
	if s.UnpackPath == "" {
		s.UnpackPath = UNPACK_URL()
	}

	if s.tbbail() != nil {
		return nil
	}

	log.Println("running tor browser with lang", s.Lang, s.UnpackPath, OS())
	switch OS() {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running tor browser with lang", s.Lang, s.UnpackPath)
			args := []string{}
			args = append(args, s.PTAS()...)
			s.tbcmd = exec.Command(s.TBPath(), args...)
			s.tbcmd.Stdout = os.Stdout
			s.tbcmd.Stderr = os.Stderr
			return s.tbcmd.Run()
		}
		log.Println("tor browser not found at", s.TBPath())
		return fmt.Errorf("tor browser not found at %s", s.TBPath())
	case "osx":
		firefoxPath := filepath.Join(s.UnpackPath, "Tor Browser.app", "Contents", "MacOS", "firefox")
		s.tbcmd = exec.Command(firefoxPath)
		s.tbcmd.Dir = s.UnpackPath
		s.tbcmd.Stdout = os.Stdout
		s.tbcmd.Stderr = os.Stderr
		defer s.tbcmd.Process.Kill()
		return s.tbcmd.Run()
	case "win":
		log.Println("Running Windows EXE", s.TBDirectory(), "firefox.exe")
		args := []string{}
		args = append(args, s.PTAS()...)
		s.tbcmd = exec.Command(filepath.Join(s.TBDirectory(), "firefox.exe"), args...)
		s.tbcmd.Dir = s.TBDirectory()
		return s.tbcmd.Run()
	default:
	}

	return nil
}

// RunTBWithLang runs the Tor Browser with the given language
func (s *Supervisor) RunTBHelpWithLang() error {
	tbget.ARCH = ARCH()
	if s.Lang == "" {
		s.Lang = DEFAULT_TB_LANG
	}
	if s.UnpackPath == "" {
		s.UnpackPath = UNPACK_URL()
	}

	if s.tbbail() != nil {
		return nil
	}

	log.Println("running tor browser with lang", s.Lang, s.UnpackPath, OS())
	switch OS() {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running tor browser with lang", s.Lang, s.UnpackPath)
			s.tbcmd = exec.Command(s.TBPath(), "--help")
			s.tbcmd.Stdout = os.Stdout
			s.tbcmd.Stderr = os.Stderr
			return s.tbcmd.Run()
		}
		log.Println("tor browser not found at", s.TBPath())
		return fmt.Errorf("tor browser not found at %s", s.TBPath())
	case "osx":
		firefoxPath := filepath.Join(s.UnpackPath, "Tor Browser.app", "Contents", "MacOS", "firefox")
		s.tbcmd = exec.Command(firefoxPath, "--help")
		s.tbcmd.Stdout = os.Stdout
		s.tbcmd.Stderr = os.Stderr
		s.tbcmd.Dir = s.TBDirectory()
		defer s.tbcmd.Process.Kill()
		return s.tbcmd.Run()
	case "win":
		log.Println("Running Windows EXE", s.TBDirectory(), "firefox.exe")
		s.tbcmd = exec.Command(filepath.Join(s.TBDirectory(), "firefox.exe"), "--help")
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

func (s *Supervisor) CopyAWOXPI(profiledata string) error {
	apath, err := filepath.Abs(profiledata)
	if err != nil {
		return err
	}
	odir := filepath.Join(apath, "extensions")
	if err := os.MkdirAll(odir, 0755); err != nil {
		return err
	}
	if !tbget.FileExists(filepath.Join(apath, "user.js")) {
		err := ioutil.WriteFile(filepath.Join(apath, "user.js"), []byte("#\n"), 0644)
		if err != nil {
			return err
		}
	}
	if !tbget.FileExists(filepath.Join(apath, "pref.js")) {
		err := ioutil.WriteFile(filepath.Join(apath, "pref.js"), []byte("#\n"), 0644)
		if err != nil {
			return err
		}
	}
	opath := filepath.Join(odir, "awo@eyedeekay.github.io.xpi")

	ipath := filepath.Join(filepath.Dir(s.UnpackPath), "awo@eyedeekay.github.io.xpi")
	if err := copy.Copy(ipath, opath); err != nil {
		return err
	}
	return nil
}

// RunTBBWithOfflineProfile runs the I2P Browser with the given language
func (s *Supervisor) RunTBBWithOfflineProfile(profiledata string, offline bool) error {
	if offline {
		if err := s.CopyAWOXPI(profiledata); err != nil {
			log.Println("Error copying AWO XPI", err)
			return err
		}
	}
	tbget.ARCH = ARCH()
	if s.Lang == "" {
		s.Lang = DEFAULT_TB_LANG
	}
	if s.UnpackPath == "" {
		s.UnpackPath = UNPACK_URL()
	}

	log.Println("running i2p in tor browser with lang", s.Lang, s.UnpackPath, OS())
	switch OS() {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running Tor browser with lang and I2P Profile", s.Lang, s.UnpackPath, s.FirefoxPath(), "--profile", profiledata)
			args := []string{"--profile", profiledata}
			args = append(args, s.PTAS()...)
			s.ibcmd = exec.Command(s.FirefoxPath(), args...)
			s.ibcmd.Stdout = os.Stdout
			s.ibcmd.Stderr = os.Stderr
			return s.ibcmd.Run()
		}
		log.Println("tor browser not found at", s.FirefoxPath())
		return fmt.Errorf("tor browser not found at %s", s.FirefoxPath())
	case "osx":
		firefoxPath := filepath.Join(s.UnpackPath, "Tor Browser.app", "Contents", "MacOS", "firefox")
		s.ibcmd = exec.Command(firefoxPath, "--profile", profiledata)
		s.ibcmd.Dir = profiledata
		s.ibcmd.Stdout = os.Stdout
		s.ibcmd.Stderr = os.Stderr
		defer s.ibcmd.Process.Kill()
		return s.ibcmd.Run()
	case "win":
		log.Println("Running Windows EXE", filepath.Join(s.TBDirectory(), "firefox.exe"), "--profile", profiledata)
		args := []string{"--profile", "."}
		args = append(args, s.PTAS()...)
		s.ibcmd = exec.Command(filepath.Join(s.TBDirectory(), "firefox.exe"), args...)
		s.ibcmd.Dir = profiledata
		s.ibcmd.Stdout = os.Stdout
		s.ibcmd.Stderr = os.Stderr
		return s.ibcmd.Run()
	default:
	}
	return nil
}

// RunTBBWithProfile runs the I2P Browser with the given language
func (s *Supervisor) RunTBBWithProfile(profiledata string) error {
	return s.RunTBBWithOfflineProfile(profiledata, false)
}

// RunI2PBWithLang runs the I2P Browser with the given language
func (s *Supervisor) RunI2PBWithLang() error {
	if s.ibbail() != nil {
		return nil
	}
	return s.RunTBBWithProfile(s.I2PDataPath())
}

// RunI2PBAppWithLang runs the I2P Browser with the given language
func (s *Supervisor) RunI2PBAppWithLang() error {
	if s.ibbail() != nil {
		return nil
	}
	return s.RunTBBWithOfflineProfile(s.I2PDataPath(), true)
}

func (s *Supervisor) torbail() error {
	_, err := net.Listen("tcp", "127.0.0.1:9050")
	if err != nil {
		log.Println("Already Running on 9050", err)
		return fmt.Errorf("Already running")
	}
	if s.torcmd != nil && s.torcmd.Process != nil && s.torcmd.ProcessState != nil {
		if s.torcmd.ProcessState.Exited() {
			log.Println("Tor exited, restarting")
			return nil
		}
		log.Println("Already Running")
		return fmt.Errorf("Already running")
	}
	log.Println("Starting Tor")
	return nil
}

// RunTorWithLang runs the Tor Exe with the given language
func (s *Supervisor) RunTorWithLang() error {
	tbget.ARCH = ARCH()
	if s.Lang == "" {
		s.Lang = DEFAULT_TB_LANG
	}
	if s.UnpackPath == "" {
		s.UnpackPath = UNPACK_URL()
	}
	if err := s.torbail(); err != nil {
		return nil
	}

	log.Println("running tor with lang", s.Lang, s.UnpackPath)
	switch OS() {
	case "linux":
		if tbget.FileExists(s.UnpackPath) {
			log.Println("running tor with lang", s.Lang, s.UnpackPath)
			s.torcmd = exec.Command(s.TorPath())
			s.torcmd.Stdout = os.Stdout
			s.torcmd.Stderr = os.Stderr
			return s.torcmd.Run()
		}
		log.Println("tor not found at", s.TorPath())
		return fmt.Errorf("tor not found at %s", s.TorPath())
	case "osx":
		torPath := filepath.Join(s.UnpackPath, "Tor Browser.app", "Contents", "Resources", "TorBrowser", "Tor", "tor")
		s.ibcmd = exec.Command(torPath)
		s.ibcmd.Dir = filepath.Dir(torPath)
		s.ibcmd.Stdout = os.Stdout
		s.ibcmd.Stderr = os.Stderr
		defer s.ibcmd.Process.Kill()
		return s.tbcmd.Run()
	case "win":
		log.Println("Running Windows EXE", filepath.Join(s.TBDirectory(), "TorBrowser", "Tor", "tor.exe"))
		s.torcmd = exec.Command(filepath.Join(s.TBDirectory(), "TorBrowser", "Tor", "tor.exe"))
		s.torcmd.Dir = s.TBDirectory()
		return s.torcmd.Run()
	default:
	}

	return nil
}

// StopTor stops tor
func (s *Supervisor) StopTor() error {
	return s.torcmd.Process.Kill()
}

// TorIsAlive returns true,true if tor is alive and belongs to us, true,false
// if it's alive and doesn't belong to us, false,false if no Tor can be found
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

// NewSupervisor creates a new supervisor
func NewSupervisor(tbPath, lang string) *Supervisor {
	return &Supervisor{
		UnpackPath: tbPath,
		Lang:       lang,
	}
}
