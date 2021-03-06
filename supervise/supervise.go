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
	"runtime"
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
	UnpackPath string
	Lang       string
	torcmd     *exec.Cmd
	//tbcmd           *exec.Cmd
	//ibcmd           *exec.Cmd
	Profile         *embed.FS
	PassThroughArgs []string
}

// PTAS is the validator for the pass-through arguments
func (s *Supervisor) PTAS() []string {
	return s.PassThroughArgs
}

// TBPath returns the path to the Tor Browser Bundle launcher
func (s *Supervisor) TBPath() string {
	switch OS() {
	case "linux":
		return filepath.Join(s.TBUnpackPath(), "Browser", "start-tor-browser")
	case "osx":
		//return filepath.Join(s.TBUnpackPath(), "Browser", "Tor Browser.app", "Contents", "MacOS", "start-tor-browser")
		return filepath.Join(s.TBUnpackPath(), "Tor Browser.app", "Contents", "MacOS", "firefox")
	case "windows":
		return filepath.Join(s.TBDirectory(), "firefox.exe")
	default:
		return filepath.Join(s.TBDirectory(), "firefox")
	}
}

// FirefoxPath returns the path to the Firefox executable inside Tor Browser
func (s *Supervisor) FirefoxPath() string {
	return s.SpecificFirefoxPath(s.TBUnpackPath())
}

// FirefoxPath returns the path to the Firefox executable inside Tor Browser
func (s *Supervisor) SpecificFirefoxPath(unpackedFirefox string) string {
	switch OS() {
	case "linux":
		return filepath.Join(s.SpecificTBDirectory(unpackedFirefox), "firefox.real")
	case "windows":
		return filepath.Join(s.SpecificTBDirectory(unpackedFirefox), "firefox.exe")
	default:
		return filepath.Join(s.SpecificTBDirectory(unpackedFirefox), "firefox")
	}
}

// SpecificTBDirectory returns the path to the Tor Browser firefox directory within an unpacked TBB
func (s *Supervisor) SpecificTBDirectory(unpacked string) string {
	return filepath.Join(unpacked, "Browser")
}

// TBDirectory returns the path to the Tor Browser firefox directory
func (s *Supervisor) TBDirectory() string {
	return filepath.Join(s.TBUnpackPath(), "Browser")
}

// TorPath returns the path to the Tor executable
func (s *Supervisor) TorPath() string {
	if OS() == "osx" {
		return filepath.Join(s.TBUnpackPath(), "Tor Browser.app", "Contents", "Resources", "TorBrowser", "Tor", "tor")
	}
	return filepath.Join(s.TBUnpackPath(), "Browser", "TorBrowser", "Tor", "tor")
}

// TorDataPath returns the path to the Tor Browser Bundle Data directory
func (s *Supervisor) TorDataPath() string {
	return filepath.Join(s.TBUnpackPath(), "Browser", "TorBrowser", "Data")
}

// I2PProfilePath returns the path to the I2P profile
func (s *Supervisor) I2PProfilePath() string {
	fp := filepath.Join(filepath.Dir(s.IBBUnpackPath()), ".i2p.firefox")
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
	fp := filepath.Join(filepath.Dir(s.IBBUnpackPath()), ".i2p.firefox.config")
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
	up := filepath.Join(filepath.Dir(s.IBBUnpackPath()), "i2p.firefox")
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

// UnpackI2PData unpacks the I2P data into the s.IBBUnpackPath()
func (s *Supervisor) UnpackI2PData() error {
	return fs.WalkDir(s.Profile, ".", func(embedpath string, d fs.DirEntry, err error) error {
		fp := filepath.Join(filepath.Dir(s.IBBUnpackPath()), ".i2p.firefox")
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(embedpath, filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox", "", -1)))
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
	up := filepath.Join(filepath.Dir(s.IBBUnpackPath()), "i2p.firefox.config")
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

// UnpackI2PAppData unpacks the I2P application data into the s.IBBUnpackPath()
func (s *Supervisor) UnpackI2PAppData() error {
	return fs.WalkDir(s.Profile, ".", func(embedpath string, d fs.DirEntry, err error) error {
		fp := filepath.Join(filepath.Dir(s.IBBUnpackPath()), ".i2p.firefox.config")
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(embedpath, filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox.config", "", -1)))
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
	return s.ibbail()
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

	log.Println("running tor browser with lang", s.Lang, s.TBUnpackPath(), OS())
	switch OS() {
	case "linux":
		if tbget.FileExists(s.TBUnpackPath()) {
			log.Println("running tor browser with lang", s.Lang, s.TBUnpackPath())
			args := []string{}
			args = append(args, s.PTAS()...)
			bcmd := exec.Command(s.TBPath(), args...)
			bcmd.Stdout = os.Stdout
			bcmd.Stderr = os.Stderr
			return bcmd.Run()
		}
		log.Println("tor browser not found at", s.TBPath())
		return fmt.Errorf("tor browser not found at %s", s.TBPath())
	case "osx":
		firefoxPath := s.TBPath() //FirefoxPath
		bcmd := exec.Command(firefoxPath)
		bcmd.Dir = s.TBUnpackPath()
		bcmd.Stdout = os.Stdout
		bcmd.Stderr = os.Stderr

		return bcmd.Run()
	case "win":
		log.Println("Running Windows EXE", s.TBDirectory(), "firefox.exe")
		args := []string{}
		args = append(args, s.PTAS()...)
		bcmd := exec.Command(s.TBPath(), args...)
		bcmd.Dir = s.TBDirectory()
		return bcmd.Run()
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

	log.Println("running tor browser with lang", s.Lang, s.TBUnpackPath(), OS())
	switch OS() {
	case "linux":
		if tbget.FileExists(s.TBUnpackPath()) {
			log.Println("running tor browser with lang", s.Lang, s.TBUnpackPath())
			bcmd := exec.Command(s.TBPath(), "--help")
			bcmd.Stdout = os.Stdout
			bcmd.Stderr = os.Stderr
			return bcmd.Run()
		}
		log.Println("tor browser not found at", s.TBPath())
		return fmt.Errorf("tor browser not found at %s", s.TBPath())
	case "osx":
		firefoxPath := s.TBPath()
		bcmd := exec.Command(firefoxPath, "--help")
		bcmd.Stdout = os.Stdout
		bcmd.Stderr = os.Stderr
		bcmd.Dir = s.TBDirectory()

		return bcmd.Run()
	case "win":
		log.Println("Running Windows EXE", s.TBDirectory(), "firefox.exe")
		bcmd := exec.Command(s.TBPath(), "--help")
		bcmd.Dir = s.TBDirectory()
		return bcmd.Run()
	default:
	}

	return nil
}

func (s *Supervisor) ibbail() error {
	processes, err := ps.Processes()
	if err != nil {
		return nil
	}
	for _, p := range processes {
		if p.Executable() == s.TorPath() {
			var err error
			s.torcmd.Process, err = os.FindProcess(p.Pid())
			if err == nil {
				return fmt.Errorf("Already running")
			}
		}
	}
	return nil
}

// RunI2PBWithLang runs the I2P Browser with the given language
func (s *Supervisor) RunI2PBWithLang() error {
	if s.ibbail() != nil {
		return nil
	}
	// export TOR_HIDE_BROWSER_LOGO=1
	os.Setenv("TOR_HIDE_BROWSER_LOGO", "1")
	return s.RunSpecificTBBWithOfflineClearnetProfile(s.I2PDataPath(), s.IBBUnpackPath(), false, false, false)
}

// RunI2PBAppWithLang runs the I2P Browser with the given language
func (s *Supervisor) RunI2PBAppWithLang() error {
	if s.ibbail() != nil {
		return nil
	}
	// export TOR_HIDE_BROWSER_LOGO=1
	os.Setenv("TOR_HIDE_BROWSER_LOGO", "1")
	return s.RunSpecificTBBWithOfflineClearnetProfile(s.I2PAppDataPath(), s.IBBUnpackPath(), true, false, false)
}

func (s *Supervisor) generateOfflineProfile(profiledata string) error {
	apath, err := filepath.Abs(profiledata)
	if err != nil {
		return err
	}
	if !tbget.FileExists(filepath.Join(apath, "user.js")) {
		err := ioutil.WriteFile(filepath.Join(apath, "user.js"), offlinebrowserjs, 0644)
		if err != nil {
			return err
		}
	}
	if !tbget.FileExists(filepath.Join(apath, "pref.js")) {
		err := ioutil.WriteFile(filepath.Join(apath, "pref.js"), offlinebrowserjs, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Supervisor) GenerateClearnetProfile(profiledata string) error {
	apath, err := filepath.Abs(profiledata)
	if err != nil {
		return err
	}

	// see: https://github.com/Whonix/tb-starter/blob/b5d2280ad445bc1fbdb613424664bf8503e6f395/usr/share/secbrowser/variables.bsh
	// export TOR_NO_DISPLAY_NETWORK_SETTINGS=1
	os.Setenv("TOR_NO_DISPLAY_NETWORK_SETTINGS", "1")
	// export TOR_HIDE_BROWSER_LOGO=1
	os.Setenv("TOR_HIDE_BROWSER_LOGO", "1")
	// export TOR_SKIP_CONTROLPORTTEST=1
	os.Setenv("TOR_SKIP_CONTROLPORTTEST", "1")
	// export TOR_SKIP_LAUNCH=1
	os.Setenv("TOR_SKIP_LAUNCH", "1")
	// export TOR_TRANSPROXY=1
	os.Setenv("TOR_TRANSPROXY", "1")

	odir := filepath.Join(apath, "extensions")
	if err := os.MkdirAll(odir, 0755); err != nil {
		return err
	}
	if !tbget.FileExists(filepath.Join(apath, "user.js")) {
		err := ioutil.WriteFile(filepath.Join(apath, "user.js"), secbrowserjs, 0644)
		if err != nil {
			return err
		}
	}
	if !tbget.FileExists(filepath.Join(apath, "pref.js")) {
		err := ioutil.WriteFile(filepath.Join(apath, "pref.js"), secbrowserjs, 0644)
		if err != nil {
			return err
		}
	}
	htmlfile := filepath.Join(apath, "index.html")
	if !tbget.FileExists(htmlfile) {
		err := ioutil.WriteFile(htmlfile, []byte(secbrowserhtml), 0644)
		if err != nil {
			return err
		}
	}

	opath := filepath.Join(odir, "uBlock0@raymondhill.net.xpi")

	ipath := filepath.Join(filepath.Dir(s.TBUnpackPath()), ".i2p.firefox", "extensions", "uBlock0@raymondhill.net.xpi")
	if err := copy.Copy(ipath, opath); err != nil {
		return err
	}
	return nil
}

func (s *Supervisor) CopyAWOXPI(profiledata string) error {
	// export TOR_HIDE_BROWSER_LOGO=1
	os.Setenv("TOR_HIDE_BROWSER_LOGO", "1")
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

	htmlfile := filepath.Join(apath, "index.html")
	if !tbget.FileExists(htmlfile) {
		err := ioutil.WriteFile(htmlfile, []byte(offlinehtml), 0644)
		if err != nil {
			return err
		}
	}

	ipath := filepath.Join(filepath.Dir(s.TBUnpackPath()), "awo@eyedeekay.github.io.xpi")
	if err := copy.Copy(ipath, opath); err != nil {
		return err
	}
	return nil
}

// RunTBBWithOfflineProfile runs the I2P Browser with the given language
func (s *Supervisor) RunTBBWithOfflineClearnetProfile(profiledata string, offline, clearnet bool) error {
	return s.RunSpecificTBBWithOfflineClearnetProfile(profiledata, s.TBUnpackPath(), offline, clearnet, false)
}

func (s *Supervisor) RunSpecificTBBWithOfflineClearnetProfile(profiledata, torbrowserdata string, offline, clearnet, editor bool) error {
	defaultpage := "about:blank"
	if clearnet {
		log.Print("Generating Clearnet Profile")
		if err := s.GenerateClearnetProfile(profiledata); err != nil {
			log.Println("Error generating Clearnet Profile", err)
			return err
		}
		defaultpage = profiledata + "/index.html"
	}
	if offline {
		if err := s.CopyAWOXPI(profiledata); err != nil {
			log.Println("Error copying AWO XPI", err)
			return err
		}
		if err := s.generateOfflineProfile(profiledata); err != nil {
			log.Println("Error generating Offline Profile", err)
			return err
		}
		if !strings.Contains(filepath.Base(profiledata), "i2p") {
			defaultpage = profiledata + "/index.html"
		}
	}
	if editor {
		defaultpage = "http://127.0.0.1:7685"
		clearnet := true
		offline := true
		if clearnet {
			log.Print("Generating Clearnet Profile")
			if err := s.GenerateClearnetProfile(profiledata); err != nil {
				log.Println("Error generating Clearnet Profile", err)
				return err
			}
		}
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
	}
	return s.RunSpecificTBBWithOfflineClearnetProfileAndPage(profiledata, torbrowserdata, offline, clearnet, defaultpage)
}

func (s *Supervisor) RunSpecificTBBWithOfflineClearnetProfileAndPage(profiledata, torbrowserdata string, offline, clearnet bool, defaultpage string) error {
	tbget.ARCH = ARCH()
	if s.Lang == "" {
		s.Lang = DEFAULT_TB_LANG
	}
	if torbrowserdata == "" {
		torbrowserdata = UNPACK_URL()
	}

	log.Println("running i2p in tor browser with lang", s.Lang, torbrowserdata, OS())
	switch OS() {
	case "linux":
		if tbget.FileExists(torbrowserdata) {
			args := []string{"--profile", profiledata, defaultpage}
			args = append(args, s.PTAS()...)
			log.Println("running Tor browser with lang and Custom Profile", s.Lang, torbrowserdata, s.SpecificFirefoxPath(torbrowserdata), args)
			bcmd := exec.Command(s.SpecificFirefoxPath(torbrowserdata), args...)
			bcmd.Stdout = os.Stdout
			bcmd.Stderr = os.Stderr
			return bcmd.Run()
		}
		log.Println("tor browser not found at", s.SpecificFirefoxPath(torbrowserdata))
		return fmt.Errorf("tor browser not found at %s", s.SpecificFirefoxPath(torbrowserdata))
	case "osx":
		firefoxPath := s.TBPath()
		args := []string{"--profile", profiledata, defaultpage}
		args = append(args, s.PTAS()...)
		log.Println("running Tor browser with lang and Custom Profile", s.Lang, torbrowserdata, firefoxPath, args)
		bcmd := exec.Command(firefoxPath, args...)
		bcmd.Dir = profiledata
		bcmd.Stdout = os.Stdout
		bcmd.Stderr = os.Stderr

		return bcmd.Run()
	case "win":
		args := []string{"--profile", profiledata, defaultpage}
		args = append(args, s.PTAS()...)
		log.Println("running Tor browser with lang and Custom Profile", s.Lang, torbrowserdata, s.SpecificFirefoxPath(torbrowserdata), args)
		bcmd := exec.Command(s.SpecificFirefoxPath(torbrowserdata), args...)
		bcmd.Dir = profiledata
		bcmd.Stdout = os.Stdout
		bcmd.Stderr = os.Stderr
		return bcmd.Run()
	default:
	}
	return nil
}

// RunTBBWithProfile runs the I2P Browser with the given language
func (s *Supervisor) RunTBBWithProfile(profiledata string) error {
	return s.RunTBBWithOfflineClearnetProfile(profiledata, false, false)
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

	log.Println("running tor with lang", s.Lang, s.TBUnpackPath())
	switch OS() {
	case "linux":
		if tbget.FileExists(s.TBUnpackPath()) {
			log.Println("running tor with lang", s.Lang, s.TBUnpackPath())
			s.torcmd = exec.Command(s.TorPath())
			s.torcmd.Stdout = os.Stdout
			s.torcmd.Stderr = os.Stderr
			return s.torcmd.Run()
		}
		log.Println("tor not found at", s.TorPath())
		return fmt.Errorf("tor not found at %s", s.TorPath())
	case "osx":
		torPath := filepath.Join(s.TBUnpackPath(), "Tor Browser.app", "Contents", "Resources", "TorBrowser", "Tor", "tor")
		s.torcmd = exec.Command(torPath)
		s.torcmd.Dir = filepath.Dir(torPath)
		s.torcmd.Stdout = os.Stdout
		s.torcmd.Stderr = os.Stderr
		defer s.torcmd.Process.Kill()
		return s.torcmd.Run()
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
	if s.torcmd != nil && s.torcmd.Process != nil && s.torcmd.ProcessState != nil {
		return s.torcmd.Process.Kill()
	}
	return nil
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

func (s *Supervisor) TBUnpackPath() string {
	return s.UnpackPath
}

func (s *Supervisor) IBBUnpackPath() string {
	return strings.Replace(s.TBUnpackPath(), "tor-browser", "i2p-browser", -1)
}

func FindEepsiteDocroot() (string, error) {
	// eepsite docroot could be at:
	// or: $I2P_CONFIG/eepsite/docroot/
	// or: $I2P/eepsite/docroot/
	// or: $HOME/.i2p/eepsite/docroot/
	// or: /var/lib/i2p/i2p-config/eepsite/docroot/
	// or: %LOCALAPPDATA\i2p\eepsite\docroot\
	// or: %APPDATA\i2p\eepsite\docroot\
	// or: %USERPROFILE\i2p\eepsite\docroot\
	SNARK_CONFIG := os.Getenv("SNARK_CONFIG")
	if SNARK_CONFIG != "" {
		checkfori2pcustom := filepath.Join(SNARK_CONFIG)
		if tbget.FileExists(checkfori2pcustom) {
			return checkfori2pcustom, nil
		}
	}

	I2P_CONFIG := os.Getenv("I2P_CONFIG")
	if I2P_CONFIG != "" {
		checkfori2pcustom := filepath.Join(I2P_CONFIG, "eepsite", "docroot")
		if tbget.FileExists(checkfori2pcustom) {
			return checkfori2pcustom, nil
		}
	}

	I2P := os.Getenv("I2P")
	if I2P != "" {
		checkfori2p := filepath.Join(I2P, "eepsite", "docroot")
		if tbget.FileExists(checkfori2p) {
			return checkfori2p, nil
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// Start by getting the home directory
	switch runtime.GOOS {
	case "windows":
		checkfori2plocal := filepath.Join(home, "AppData", "Local", "i2p", "eepsite", "docroot")
		if tbget.FileExists(checkfori2plocal) {
			return checkfori2plocal, nil
		}
		checkfori2proaming := filepath.Join(home, "AppData", "Roaming", "i2p", "eepsite", "docroot")
		if tbget.FileExists(checkfori2proaming) {
			return checkfori2proaming, nil
		}
	case "linux":
		checkfori2phome := filepath.Join(home, ".i2p", "eepsite", "docroot")
		if tbget.FileExists(checkfori2phome) {
			return checkfori2phome, nil
		}
		checkfori2pservice := filepath.Join("/var/lib/i2p/i2p-config", "eepsite", "docroot")
		if tbget.FileExists(checkfori2pservice) {
			return checkfori2pservice, nil
		}
	case "darwin":
		return "", fmt.Errorf("FindSnarkDirectory: Automatic torrent generation is not supported on MacOS, for now copy the files manually")
	}
	return "", fmt.Errorf("FindSnarkDirectory: Unable to find snark directory")

}

// RunTBBWithOfflineProfile runs the I2P Browser with the given language
func (s *Supervisor) RunI2PSiteEditorWithOfflineClearnetProfile(profiledata string) error {
	return s.RunSpecificTBBWithOfflineClearnetProfile(profiledata, s.IBBUnpackPath(), true, true, true)
}
