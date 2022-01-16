package tbget

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cloudfoundry/jibber_jabber"

	"github.com/jchavannes/go-pgp/pgp"
	"golang.org/x/crypto/openpgp"
)

var wd, _ = os.Getwd()

var DOWNLOAD_PATH = filepath.Join(wd, "tor-browser")

const TOR_UPDATES_URL string = "https://aus1.torproject.org/torbrowser/update_3/release/downloads.json"

var (
	defaultIETFLang, _ = jibber_jabber.DetectIETF()
)

func GetRuntimePair() string {
	var OS, ARCH string
	switch runtime.GOOS {
	case "darwin":
		OS = "osx"
	case "linux":
		OS = "linux"
	case "windows":
		OS = "win"
	default:
		OS = "unknown"
	}
	switch runtime.GOARCH {
	case "amd64":
		ARCH = "64"
	case "386":
		ARCH = "32"
	default:
		ARCH = "unknown"
	}
	return fmt.Sprintf("%s%s", OS, ARCH)
}

func GetUpdater() (string, string, error) {
	return GetUpdaterForLang(defaultIETFLang)
}

func GetUpdaterForLang(ietf string) (string, string, error) {
	jsonText, err := http.Get(TOR_UPDATES_URL)
	if err != nil {
		return "", "", fmt.Errorf("GetUpdaterForLang: %s", err)
	}
	defer jsonText.Body.Close()
	return GetUpdaterForLangFromJson(jsonText.Body, ietf)
}

func GetUpdaterForLangFromJson(body io.ReadCloser, ietf string) (string, string, error) {
	jsonBytes, err := io.ReadAll(body)
	if err != nil {
		return "", "", fmt.Errorf("GetUpdaterForLangFromJson: %s", err)
	}
	if err = ioutil.WriteFile(filepath.Join(DOWNLOAD_PATH, "downloads.json"), jsonBytes, 0644); err != nil {
		return "", "", fmt.Errorf("GetUpdaterForLangFromJson: %s", err)
	}
	return GetUpdaterForLangFromJsonBytes(jsonBytes, ietf)
}

func GetUpdaterForLangFromJsonBytes(jsonBytes []byte, ietf string) (string, string, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &dat); err != nil {
		return "", "", fmt.Errorf("FuncName: %s", err)
	}
	if platform, ok := dat["downloads"]; ok {
		rtp := GetRuntimePair()
		if updater, ok := platform.(map[string]interface{})[rtp]; ok {
			if langUpdater, ok := updater.(map[string]interface{})[ietf]; ok {
				return langUpdater.(map[string]interface{})["binary"].(string), langUpdater.(map[string]interface{})["sig"].(string), nil
			}
			// If we didn't find the language, try splitting at the hyphen
			lang := strings.Split(ietf, "-")[0]
			if langUpdater, ok := updater.(map[string]interface{})[lang]; ok {
				return langUpdater.(map[string]interface{})["binary"].(string), langUpdater.(map[string]interface{})["sig"].(string), nil
			}
			// If we didn't find the language after splitting at the hyphen, try the default
			return GetUpdaterForLangFromJsonBytes(jsonBytes, defaultIETFLang)
		}
	}
	return "", "", fmt.Errorf("GetUpdaterForLangFromJsonBytes: no updater for language %s", ietf)
}

func SingleFileDownload(url, name string) (string, error) {
	path := filepath.Join(DOWNLOAD_PATH, name)
	if !BotherToDownload(url, name) {
		fmt.Printf("No updates required, skipping download of %s\n", name)
		return path, nil
	}
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
	return path, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func BotherToDownload(url, name string) bool {
	path := filepath.Join(DOWNLOAD_PATH, name)
	if !fileExists(path) {
		return true
	}
	defer ioutil.WriteFile(filepath.Join(DOWNLOAD_PATH, name+".last-url"), []byte(url), 0644)
	lastUrl, err := ioutil.ReadFile(filepath.Join(DOWNLOAD_PATH, name+".last-url"))
	if err != nil {
		return true
	}
	if string(lastUrl) == url {
		return false
	}
	return true
}

func DownloadUpdater() (string, string, error) {
	binary, sig, err := GetUpdater()
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdater: %s", err)
	}
	sigpath, err := SingleFileDownload(sig, "tor-browser-"+GetRuntimePair()+"-"+defaultIETFLang+".tar.xz.asc")
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdater: %s", err)
	}
	binpath, err := SingleFileDownload(binary, "tor-browser-"+GetRuntimePair()+"-"+defaultIETFLang+".tar.xz")
	if err != nil {
		return "", sigpath, fmt.Errorf("DownloadUpdater: %s", err)
	}
	return binpath, sigpath, nil
}

func DownloadUpdaterForLang(ietf string) (string, string, error) {
	binary, sig, err := GetUpdaterForLang(ietf)
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdaterForLang: %s", err)
	}
	sigpath, err := SingleFileDownload(sig, "tor-browser-"+GetRuntimePair()+"-"+ietf+".tar.xz.asc")
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdaterForLang: %s", err)
	}
	binpath, err := SingleFileDownload(binary, "tor-browser-"+GetRuntimePair()+"-"+ietf+".tar.xz")
	if err != nil {
		return "", sigpath, fmt.Errorf("DownloadUpdaterForLang: %s", err)
	}
	return binpath, sigpath, nil
}

func CheckSignature(binpath, sigpath string) error {
	var pkBytes []byte
	var pk *openpgp.Entity
	var sig []byte
	var bin []byte
	var err error
	if pkBytes, err = ioutil.ReadFile(filepath.Join(DOWNLOAD_PATH, "TPO-signing-key.pub")); err != nil {
		return fmt.Errorf("CheckSignature pkBytes: %s", err)
	}
	if pk, err = pgp.GetEntity(pkBytes, nil); err != nil {
		return fmt.Errorf("CheckSignature pk: %s", err)
	}
	if bin, err = ioutil.ReadFile(binpath); err != nil {
		return fmt.Errorf("CheckSignature bin: %s", err)
	}
	if sig, err = ioutil.ReadFile(sigpath); err != nil {
		return fmt.Errorf("CheckSignature sig: %s", err)
	}
	if err = pgp.Verify(pk, sig, bin); err != nil {
		return nil
	}
	err = fmt.Errorf("signature check failed")
	return fmt.Errorf("CheckSignature: %s", err)
}

func BoolCheckSignature(binpath, sigpath string) bool {
	err := CheckSignature(binpath, sigpath)
	return err == nil
}
