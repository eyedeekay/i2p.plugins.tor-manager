package tbget

import (
	"archive/tar"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cloudfoundry/jibber_jabber"
	sam "github.com/eyedeekay/sam3/helper"
	"github.com/ulikunitz/xz"

	"github.com/jchavannes/go-pgp/pgp"
	"golang.org/x/crypto/openpgp"
)

var WORKING_DIR = ""

func DefaultDir() string {
	if WORKING_DIR == "" {
		WORKING_DIR, _ = os.Getwd()
	}
	if !FileExists(WORKING_DIR) {
		os.MkdirAll(WORKING_DIR, 0755)
	}
	return WORKING_DIR
}

var UNPACK_PATH = filepath.Join(DefaultDir(), "unpack")
var DOWNLOAD_PATH = filepath.Join(DefaultDir(), "tor-browser")

const TOR_UPDATES_URL string = "https://aus1.torproject.org/torbrowser/update_3/release/downloads.json"

var (
	DefaultIETFLang, _ = jibber_jabber.DetectIETF()
)

type TBDownloader struct {
	UnpackPath   string
	DownloadPath string
	Lang         string
	OS, ARCH     string
	Verbose      bool
	Profile      *embed.FS
}

var OS = "linux"
var ARCH = "64"

func NewTBDownloader(lang string, os, arch string, content *embed.FS) *TBDownloader {
	return &TBDownloader{
		Lang:         lang,
		DownloadPath: DOWNLOAD_PATH,
		UnpackPath:   UNPACK_PATH,
		OS:           os,
		ARCH:         arch,
		Verbose:      false,
		Profile:      content,
	}
}

func (t *TBDownloader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.Replace(r.URL.Path, "..", "", -1)
	ext := filepath.Ext(r.URL.Path)
	if ext == ".json" {
		w.Header().Set("Content-Type", "application/json")
		if FileExists(filepath.Join(t.DownloadPath, "mirror.json")) {
			http.ServeFile(w, r, filepath.Join(t.DownloadPath, "mirror.json"))
		}
	}
	if FileExists(filepath.Join(t.DownloadPath, r.URL.Path)) {
		http.ServeFile(w, r, filepath.Join(t.DownloadPath, r.URL.Path))
		return
	}
}

func (t *TBDownloader) Serve() {
	samlistener, err := sam.I2PListener("tor-mirror", "127.0.0.1:7656", "tor-mirror")
	if err != nil {
		log.Fatal(err)
	}
	defer samlistener.Close()
	http.Serve(samlistener, t)
}

func (t *TBDownloader) GetRuntimePair() string {
	if t.OS != "" && t.ARCH != "" {
		return fmt.Sprintf("%s%s", t.OS, t.ARCH)
	}
	switch runtime.GOOS {
	case "darwin":
		t.OS = "osx"
	case "linux":
		t.OS = "linux"
	case "windows":
		t.OS = "win"
	default:
		t.OS = "unknown"
	}
	switch runtime.GOARCH {
	case "amd64":
		t.ARCH = "64"
	case "386":
		t.ARCH = "32"
	default:
		t.ARCH = "unknown"
	}
	return fmt.Sprintf("%s%s", t.OS, t.ARCH)
}

func (t *TBDownloader) GetUpdater() (string, string, error) {
	return t.GetUpdaterForLang(t.Lang)
}

func (t *TBDownloader) GetUpdaterForLang(ietf string) (string, string, error) {
	jsonText, err := http.Get(TOR_UPDATES_URL)
	if err != nil {
		return "", "", fmt.Errorf("t.GetUpdaterForLang: %s", err)
	}
	defer jsonText.Body.Close()
	return t.GetUpdaterForLangFromJson(jsonText.Body, ietf)
}

func (t *TBDownloader) GetUpdaterForLangFromJson(body io.ReadCloser, ietf string) (string, string, error) {
	jsonBytes, err := io.ReadAll(body)
	if err != nil {
		return "", "", fmt.Errorf("t.GetUpdaterForLangFromJson: %s", err)
	}
	t.MakeTBDirectory()
	if err = ioutil.WriteFile(filepath.Join(t.DownloadPath, "downloads.json"), jsonBytes, 0644); err != nil {
		return "", "", fmt.Errorf("t.GetUpdaterForLangFromJson: %s", err)
	}
	return t.GetUpdaterForLangFromJsonBytes(jsonBytes, ietf)
}

func (t *TBDownloader) Log(function, message string) {
	if t.Verbose {
		log.Println(fmt.Sprintf("%s: %s", function, message))
	}
}

func (t *TBDownloader) MakeTBDirectory() {
	os.MkdirAll(t.DownloadPath, 0755)

	path := filepath.Join("", "tor-browser", "TPO-signing-key.pub")
	if !FileExists(path) {
		t.Log("MakeTBDirectory()", "Initial TPO signing key not found, using the one embedded in the executable")
		bytes, err := t.Profile.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		t.Log("MakeTBDirectory()", "Writing TPO signing key to disk")
		ioutil.WriteFile(filepath.Join(t.DownloadPath, "TPO-signing-key.pub"), bytes, 0644)
		t.Log("MakeTBDirectory()", "Writing TPO signing key to disk complete")
	}
}

func (t *TBDownloader) GetUpdaterForLangFromJsonBytes(jsonBytes []byte, ietf string) (string, string, error) {
	t.MakeTBDirectory()
	var dat map[string]interface{}
	t.Log("GetUpdaterForLangFromJsonBytes()", "Parsing JSON")
	if err := json.Unmarshal(jsonBytes, &dat); err != nil {
		return "", "", fmt.Errorf("func (t *TBDownloader)Name: %s", err)
	}
	t.Log("GetUpdaterForLangFromJsonBytes()", "Parsing JSON complete")
	if platform, ok := dat["downloads"]; ok {
		rtp := t.GetRuntimePair()
		if updater, ok := platform.(map[string]interface{})[rtp]; ok {
			if langUpdater, ok := updater.(map[string]interface{})[ietf]; ok {
				t.Log("GetUpdaterForLangFromJsonBytes()", "Found updater for language")
				return langUpdater.(map[string]interface{})["binary"].(string), langUpdater.(map[string]interface{})["sig"].(string), nil
			}
			// If we didn't find the language, try splitting at the hyphen
			lang := strings.Split(ietf, "-")[0]
			if langUpdater, ok := updater.(map[string]interface{})[lang]; ok {
				t.Log("GetUpdaterForLangFromJsonBytes()", "Found updater for backup language")
				return langUpdater.(map[string]interface{})["binary"].(string), langUpdater.(map[string]interface{})["sig"].(string), nil
			}
			// If we didn't find the language after splitting at the hyphen, try the default
			t.Log("GetUpdaterForLangFromJsonBytes()", "Last attempt, trying default language")
			return t.GetUpdaterForLangFromJsonBytes(jsonBytes, t.Lang)
		} else {
			return "", "", fmt.Errorf("t.GetUpdaterForLangFromJsonBytes: no updater for platform %s", rtp)
		}
	}
	return "", "", fmt.Errorf("t.GetUpdaterForLangFromJsonBytes: %s", ietf)
}

func (t *TBDownloader) SingleFileDownload(url, name string) (string, error) {
	t.MakeTBDirectory()
	path := filepath.Join(t.DownloadPath, name)
	t.Log("SingleFileDownload()", fmt.Sprintf("Checking for updates %s to %s", url, path))
	if !t.BotherToDownload(url, name) {
		t.Log("SingleFileDownload()", "File already exists, skipping download")
		return path, nil
	}
	t.Log("SingleFileDownload()", "Downloading file")
	file, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("SingleFileDownload: %s", err)
	}
	defer file.Body.Close()
	outFile, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("SingleFileDownload: %s", err)
	}
	defer outFile.Close()
	io.Copy(outFile, file.Body)
	t.Log("SingleFileDownload()", "Downloading file complete")
	return path, nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (t *TBDownloader) BotherToDownload(url, name string) bool {
	path := filepath.Join(t.DownloadPath, name)
	if !FileExists(path) {
		return true
	}
	defer ioutil.WriteFile(filepath.Join(t.DownloadPath, name+".last-url"), []byte(url), 0644)
	lastUrl, err := ioutil.ReadFile(filepath.Join(t.DownloadPath, name+".last-url"))
	if err != nil {
		return true
	}
	if string(lastUrl) == url {
		return false
	}
	return true
}

func (t *TBDownloader) NamePerPlatform(ietf string) string {
	extension := "tar.xz"
	windowsonly := ""
	switch runtime.GOOS {
	case "darwin":
		extension = "dmg"
	case "windows":
		windowsonly = "-installer"
		extension = "exe"
	}
	return fmt.Sprintf("torbrowser%s-%s-%s.%s", windowsonly, t.GetRuntimePair(), ietf, extension)
}

func (t *TBDownloader) DownloadUpdater() (string, string, error) {
	binary, sig, err := t.GetUpdater()
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdater: %s", err)
	}
	sigpath, err := t.SingleFileDownload(sig, t.NamePerPlatform(t.Lang)+".asc")
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdater: %s", err)
	}
	binpath, err := t.SingleFileDownload(binary, t.NamePerPlatform(t.Lang))
	if err != nil {
		return "", sigpath, fmt.Errorf("DownloadUpdater: %s", err)
	}
	return binpath, sigpath, nil
}

func (t *TBDownloader) DownloadUpdaterForLang(ietf string) (string, string, error) {
	binary, sig, err := t.GetUpdaterForLang(ietf)
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdaterForLang: %s", err)
	}

	sigpath, err := t.SingleFileDownload(sig, t.NamePerPlatform(ietf)+".asc")
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdaterForLang: %s", err)
	}
	binpath, err := t.SingleFileDownload(binary, t.NamePerPlatform(ietf))
	if err != nil {
		return "", sigpath, fmt.Errorf("DownloadUpdaterForLang: %s", err)
	}
	return binpath, sigpath, nil
}

func (t *TBDownloader) UnpackUpdater(binpath string) (string, error) {
	t.Log("UnpackUpdater()", fmt.Sprintf("Unpacking %s", binpath))
	if t.OS == "win" {
		installPath := filepath.Join(t.UnpackPath, "tor-browser_"+t.Lang)
		t.Log("UnpackUpdater()", "Windows updater, running silent NSIS installer")
		t.Log("UnpackUpdater()", fmt.Sprintf("Running %s %s %s", binpath, "/S", "/D="+installPath))
		cmd := exec.Command(binpath, "/S", "/D="+installPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("UnpackUpdater: windows exec fail %s", err)
		}
		return installPath, nil
	}
	if t.OS == "osx" {
		cmd := exec.Command("open", "-W", "-n", "-a", "\""+t.UnpackPath+"\"", "\""+binpath+"\"")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("UnpackUpdater: osx open/mount fail %s", err)
		}
	}
	if FileExists(filepath.Join(t.UnpackPath, "tor-browser_"+t.Lang)) {
		return filepath.Join(t.UnpackPath, "tor-browser_"+t.Lang), nil
	}
	fmt.Printf("Unpacking %s %s\n", binpath, t.UnpackPath)
	os.MkdirAll(t.UnpackPath, 0755)
	UNPACK_DIRECTORY, err := os.Open(t.UnpackPath)
	if err != nil {
		return "", fmt.Errorf("UnpackUpdater: %s", err)
	}
	defer UNPACK_DIRECTORY.Close()
	xzfile, err := os.Open(binpath)
	if err != nil {
		return "", fmt.Errorf("UnpackUpdater: %s", err)
	}
	defer xzfile.Close()
	xzReader, err := xz.NewReader(xzfile)
	if err != nil {
		return "", fmt.Errorf("UnpackUpdater: %s", err)
	}
	tarReader := tar.NewReader(xzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("UnpackUpdater: %s", err)
		}
		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(filepath.Join(UNPACK_DIRECTORY.Name(), header.Name), 0755)
			continue
		}
		filename := filepath.Join(UNPACK_DIRECTORY.Name(), header.Name)
		file, err := os.Create(filename)
		if err != nil {
			return "", fmt.Errorf("UnpackUpdater: %s", err)
		}
		defer file.Close()
		io.Copy(file, tarReader)
		mode := header.FileInfo().Mode()
		//remember to chmod the file afterwards
		file.Chmod(mode)
		if t.Verbose {
			fmt.Printf("Unpacked %s\n", header.Name)
		}
	}
	return filepath.Join(t.UnpackPath, "tor-browser_"+t.Lang), nil
}

func (t *TBDownloader) CheckSignature(binpath, sigpath string) (string, error) {
	var pkBytes []byte
	var pk *openpgp.Entity
	var sig []byte
	var bin []byte
	var err error
	if pkBytes, err = ioutil.ReadFile(filepath.Join(t.DownloadPath, "TPO-signing-key.pub")); err != nil {
		return "", fmt.Errorf("CheckSignature pkBytes: %s", err)
	}
	if pk, err = pgp.GetEntity(pkBytes, nil); err != nil {
		return "", fmt.Errorf("CheckSignature pk: %s", err)
	}
	if bin, err = ioutil.ReadFile(binpath); err != nil {
		return "", fmt.Errorf("CheckSignature bin: %s", err)
	}
	if sig, err = ioutil.ReadFile(sigpath); err != nil {
		return "", fmt.Errorf("CheckSignature sig: %s", err)
	}
	if err = pgp.Verify(pk, sig, bin); err != nil {
		return t.UnpackUpdater(binpath)
		//return nil
	}
	err = fmt.Errorf("signature check failed")
	return "", fmt.Errorf("CheckSignature: %s", err)
}

func (t *TBDownloader) BoolCheckSignature(binpath, sigpath string) bool {
	_, err := t.CheckSignature(binpath, sigpath)
	return err == nil
}
