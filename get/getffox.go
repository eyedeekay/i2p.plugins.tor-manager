package tbget

import (
	"archive/tar"
	"compress/bzip2"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FFOX_UPDATES_URL is the URL to the Firefox updates page
const FFOX_UPDATES_URL string = "https://download.mozilla.org/?product=firefox-latest&os=%s&lang=%s"

type FFDownloader TBDownloader

// NewFirefoxDownloader returns a new FFDownloader with the given language, using the FFDownloader's OS/ARCH pair
func NewFirefoxDownloader(lang string, os, arch string, content *embed.FS) *FFDownloader {
	OS = os
	ARCH = arch
	return &FFDownloader{
		Lang:         lang,
		DownloadPath: DOWNLOAD_FIREFOX_PATH(),
		UnpackPath:   UNPACK_FIREFOX_PATH(),
		OS:           os,
		ARCH:         arch,
		Verbose:      false,
		Profile:      content,
		Mirror:       "https://download.mozilla.org/?product=firefox-latest",
	}
}

// DOWNLOAD_FIREFOX_PATH returns the path to the downloads.
func DOWNLOAD_FIREFOX_PATH() string {
	var DOWNLOAD_PATH = filepath.Join(DefaultDir(), "firefox")
	return DOWNLOAD_PATH
}

// UNPACK_FIREFOX_PATH returns the path to the unpacked files.
func UNPACK_FIREFOX_PATH() string {
	var UNPACK_FIREFOX_PATH = filepath.Join(DefaultDir(), "unpack-firefox")
	return UNPACK_FIREFOX_PATH
}

func (t FFDownloader) GetRuntimePair() string {
	tbd := TBDownloader(t)
	return tbd.GetRuntimePair()
}

// GetLatestFirefoxVersionURL returns the URL to the latest Firefox version for the given os and lang
func (t *FFDownloader) GetLatestFirefoxVersionURL(os, lang string) string {
	return fmt.Sprintf(FFOX_UPDATES_URL, t.GetRuntimePair(), lang)
}

// GetLatestFirefoxVersionLinuxSigURL returns the URL to the latest Firefox version detatched signature for the given os and lang
func (t *FFDownloader) GetLatestFirefoxVersionLinuxSigURL(os, lang string) string {
	return t.GetLatestFirefoxVersionURL(os, lang) + ".asc"
}

// GetFirefoxUpdater gets the updater URL for the t.Lang. It returns
// the URL, a detatched sig if available for the platform, or an error
func (t *FFDownloader) GetFirefoxUpdater() (string, string, error) {
	return t.GetLatestFirefoxVersionURL(t.OS, t.Lang), t.GetLatestFirefoxVersionLinuxSigURL(t.OS, t.Lang), nil
}

// GetFirefoxUpdaterForLang gets the updater URL for the given language, overriding
// the t.Lang. It returns the URL, a detatched sig if available for the platform, or an error
func (t *FFDownloader) GetFirefoxUpdaterForLang(ietf string) (string, string, error) {
	return t.GetLatestFirefoxVersionURL(t.OS, ietf), t.GetLatestFirefoxVersionLinuxSigURL(t.OS, ietf), nil
}

// SendFirefoxVersionHEADRequest sends a HEAD request to the Firefox version URL
func (t *FFDownloader) SendFirefoxVersionHEADRequest() (string, error) {
	url := t.GetLatestFirefoxVersionURL(t.OS, t.Lang)
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", fmt.Errorf("t.SendFirefoxVersionHEADRequest: %s", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:68.0) Gecko/20100101 Firefox/68.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("t.SendFirefoxVersionHEADRequest: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("t.SendFirefoxVersionHEADRequest: %s", resp.Status)
	}
	return resp.Header.Get("Location"), nil
}

// ExtractFirefoxVersion extracts the Firefox version from the updater URL
func (t *FFDownloader) ExtractFirefoxVersion() (string, error) {
	url, err := t.SendFirefoxVersionHEADRequest()
	if err != nil {
		return "", fmt.Errorf("t.ExtractFirefoxVersion: %s", err)
	}
	// get the last element of the URL
	url = strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	// remove all file extensions
	url = strings.Replace(url, ".tar.xz", "", -1)
	url = strings.Replace(url, ".tar.bz2", "", -1)
	url = strings.Replace(url, ".tar.gz", "", -1)
	url = strings.Replace(url, ".zip", "", -1)
	url = strings.Replace(url, ".exe", "", -1)
	url = strings.Replace(url, ".msi", "", -1)
	url = strings.Replace(url, ".dmg", "", -1)
	return url, nil
}

// NamePerPlatformFirefox returns the name of the Firefox package per platform.
func (t *FFDownloader) NamePerPlatformFirefox(ietf string) string {
	extension := "tar.bz2"
	windowsonly := ""
	switch t.OS {
	case "osx":
		extension = "dmg"
	case "win":
		windowsonly = "-setup"
		extension = "exe"
	}
	return fmt.Sprintf("firefox%s-%s-%s.%s", windowsonly, t.GetRuntimePair(), ietf, extension)
}

// FirefoxBrowserDir returns the path to the directory where the Firefox browser is installed.
func (t *FFDownloader) FirefoxBrowserDir() string {
	return filepath.Join(t.UnpackPath, "firefox_"+t.Lang)
}

func (t *FFDownloader) Log(function, message string) {
	if t.Verbose {
		log.Println(fmt.Sprintf("%s: %s", function, message))
	}
}

// UnpackFirefox unpacks the Firefox package to the t.FirefoxBrowserDir()
func (t *FFDownloader) UnpackFirefox(binpath string) (string, error) {
	t.Log("UnpackFirefox()", fmt.Sprintf("Unpacking %s", binpath))
	if t.OS == "win" {
		installPath := t.FirefoxBrowserDir()
		if !FileExists(installPath) {
			t.Log("UnpackFirefox()", "Windows updater, running silent NSIS installer")
			t.Log("UnpackFirefox()", fmt.Sprintf("Running %s %s %s", binpath, "/S", "/D="+installPath))
			cmd := exec.Command(binpath, "/S", "/D="+installPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return "", fmt.Errorf("UnpackFirefox: windows exec fail %s", err)
			}
		}
		return installPath, nil
	}
	if t.OS == "osx" {
		cmd := exec.Command("open", "-W", "-n", "-a", "\""+t.UnpackPath+"\"", "\""+binpath+"\"")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("UnpackFirefox: osx open/mount fail %s", err)
		}
		//TODO: this might just need to be a hardcoded app path
		return t.UnpackPath, nil
	}
	if FileExists(t.FirefoxBrowserDir()) {
		return t.FirefoxBrowserDir(), nil
	}
	fmt.Printf("Unpacking %s %s\n", binpath, t.UnpackPath)
	os.MkdirAll(t.UnpackPath, 0755)
	UNPACK_DIRECTORY, err := os.Open(t.UnpackPath)
	if err != nil {
		return "", fmt.Errorf("UnpackFirefox: directory error %s", err)
	}
	defer UNPACK_DIRECTORY.Close()
	bzfile, err := os.Open(binpath)
	if err != nil {
		return "", fmt.Errorf("UnpackFirefox: BZFile error %s", err)
	}
	defer bzfile.Close()
	bzReader := bzip2.NewReader(bzfile)
	tarReader := tar.NewReader(bzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("UnpackFirefox: Tar looper Error %s", err)
		}
		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(filepath.Join(UNPACK_DIRECTORY.Name(), header.Name), 0755)
			continue
		}
		filename := filepath.Join(UNPACK_DIRECTORY.Name(), header.Name)
		file, err := os.Create(filename)
		if err != nil {
			//return "",
			fmt.Printf("UnpackFirefox: Tar unpacker error %s", err)
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
	return t.FirefoxBrowserDir(), nil
}

// DownloadFirefoxUpdater downloads the updater for the t.Lang. It returns
// the path to the downloaded updater and the downloaded detatched signature,
// or an error if one is encountered.
func (t *FFDownloader) DownloadFirefoxUpdater() (string, string, error) {
	return t.DownloadFirefoxUpdaterForLang(t.Lang)
}

func (t FFDownloader) SingleFileDownload(dl, name string, rangebottom int64) (string, error) {
	tbd := TBDownloader(t)
	return tbd.SingleFileDownload(dl, name, rangebottom)
}

// DownloadFirefoxUpdaterForLang downloads the updater for the given language, overriding
// t.Lang. It returns the path to the downloaded updater and the downloaded
// detatched signature, or an error if one is encountered.
func (t *FFDownloader) DownloadFirefoxUpdaterForLang(ietf string) (string, string, error) {
	binary, sig, err := t.GetFirefoxUpdaterForLang(ietf)
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdater: %s", err)
	}
	sigpath := ""
	if t.OS == "linux" {
		sigpath, err = t.SingleFileDownload(sig, t.NamePerPlatformFirefox(ietf)+".asc", 0)
		if err != nil {
			return "", "", fmt.Errorf("DownloadUpdater: %s", err)
		}
	}
	binpath, err := t.SingleFileDownload(binary, t.NamePerPlatformFirefox(ietf), 0)
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdater: %s", err)
	}
	return binpath, sigpath, nil
}

// CheckFirefoxSignature checks the signature of the updater.
// it returns an error if one is encountered. If not, it
// runs the updater and returns an error if one is encountered.
func (t *FFDownloader) CheckFirefoxSignature(binpath, sigpath string) (string, error) {
	if t.OS == "linux" {
		var err error
		pk := filepath.Join(t.DownloadPath, "TPO-signing-key.pub")
		if err = Verify(pk, sigpath, binpath); err == nil {
			t.Log("CheckFirefoxSignature: signature", "verified successfully")
			return t.UnpackFirefox(binpath)
		}
		return "", fmt.Errorf("CheckSignature: %s", err)
	}
	return "", nil
}

// BoolCheckFirefoxSignature turns CheckFirefoxSignature into a bool.
func (t *FFDownloader) BoolCheckFirefoxSignature(binpath, sigpath string) bool {
	_, err := t.CheckFirefoxSignature(binpath, sigpath)
	return err == nil
}

func (t FFDownloader) MakeTBDirectory() {
	tbd := TBDownloader(t)
	tbd.MakeTBDirectory()
}
