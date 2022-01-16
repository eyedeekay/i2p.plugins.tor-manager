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

	"aead.dev/minisign"
	"github.com/cloudfoundry/jibber_jabber"
	//"encoding/json"
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
		return "", "", err
	}
	defer jsonText.Body.Close()
	return GetUpdaterForLangFromJson(jsonText.Body, ietf)
}

func GetUpdaterForLangFromJson(body io.ReadCloser, ietf string) (string, string, error) {
	jsonBytes, err := io.ReadAll(body)
	if err != nil {
		return "", "", err
	}
	return GetUpdaterForLangFromJsonBytes(jsonBytes, ietf)
}

func GetUpdaterForLangFromJsonBytes(jsonBytes []byte, ietf string) (string, string, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &dat); err != nil {
		return "", "", err
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
	return "", "", fmt.Errorf("no updater for language %s", ietf)
}

func SingleFileDownload(url, name string) (string, error) {
	file, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer file.Body.Close()
	path := filepath.Join(DOWNLOAD_PATH, name)
	outFile, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	io.Copy(outFile, file.Body)
	return path, nil
}

func DownloadUpdater() (string, string, error) {
	binary, sig, err := GetUpdater()
	if err != nil {
		return "", "", err
	}
	sigpath, err := SingleFileDownload(sig, "tor-browser-"+GetRuntimePair()+"-"+defaultIETFLang+".tar.xz.asc")
	if err != nil {
		return "", "", err
	}
	binpath, err := SingleFileDownload(binary, "tor-browser-"+GetRuntimePair()+"-"+defaultIETFLang+".tar.xz")
	if err != nil {
		return "", sigpath, err
	}
	return binpath, sigpath, nil
}

func DownloadUpadterForLang(ietf string) (string, string, error) {
	binary, sig, err := GetUpdaterForLang(ietf)
	if err != nil {
		return "", "", err
	}
	sigpath, err := SingleFileDownload(sig, "tor-browser-"+GetRuntimePair()+"-"+ietf+".tar.xz.asc")
	if err != nil {
		return "", "", err
	}
	binpath, err := SingleFileDownload(binary, "tor-browser-"+GetRuntimePair()+"-"+ietf+".tar.xz")
	if err != nil {
		return "", sigpath, err
	}
	return binpath, sigpath, nil
}

func CheckSignature(binpath, sigpath string) error {
	var pk minisign.PublicKey
	var sig minisign.Signature
	var bin []byte
	var sigBytes []byte
	var err error
	if pk, err = minisign.PublicKeyFromFile(filepath.Join(DOWNLOAD_PATH, "TPO-signing-key.pub")); err != nil {
		return err
	}
	if bin, err = ioutil.ReadFile(binpath); err != nil {
		return err
	}
	if sig, err = minisign.SignatureFromFile(sigpath); err != nil {
		return err
	}
	if sigBytes, err = sig.MarshalText(); err != nil {
		return err
	}
	if minisign.Verify(pk, bin, sigBytes) {
		return nil
	}
	err = fmt.Errorf("signature check failed")
	return err
}

func BoolCheckSignature(binpath, sigpath string) bool {
	err := CheckSignature(binpath, sigpath)
	return err == nil
}
