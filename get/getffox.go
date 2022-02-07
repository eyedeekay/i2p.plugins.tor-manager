package tbget

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jchavannes/go-pgp/pgp"
	"github.com/ulikunitz/xz"
	"golang.org/x/crypto/openpgp"
)

// FFOX_UPDATES_URL is the URL to the Firefox updates page
const FFOX_UPDATES_URL string = "https://download.mozilla.org/?product=firefox-latest&os=%s&lang=%s"

// GetLatestFirefoxVersionURL returns the URL to the latest Firefox version for the given os and lang
func (t *TBDownloader) GetLatestFirefoxVersionURL(os, lang string) string {
	return fmt.Sprintf(FFOX_UPDATES_URL, t.GetRuntimePair(), lang)
}

// GetLatestFirefoxVersionLinuxSigURL returns the URL to the latest Firefox version detatched signature for the given os and lang
func (t *TBDownloader) GetLatestFirefoxVersionLinuxSigURL(os, lang string) string {
	return t.GetLatestFirefoxVersionURL(os, lang) + ".asc"
}

// GetFirefoxUpdater gets the updater URL for the t.Lang. It returns
// the URL, a detatched sig if available for the platform, or an error
func (t *TBDownloader) GetFirefoxUpdater() (string, string, error) {
	return t.GetLatestFirefoxVersionURL(t.OS, t.Lang), t.GetLatestFirefoxVersionLinuxSigURL(t.OS, t.Lang), nil
}

// GetFirefoxUpdaterForLang gets the updater URL for the given language, overriding
// the t.Lang. It returns the URL, a detatched sig if available for the platform, or an error
func (t *TBDownloader) GetFirefoxUpdaterForLang(ietf string) (string, string, error) {
	return t.GetLatestFirefoxVersionURL(t.OS, ietf), t.GetLatestFirefoxVersionLinuxSigURL(t.OS, ietf), nil
}

// SendFirefoxVersionHEADRequest sends a HEAD request to the Firefox version URL
func (t *TBDownloader) SendFirefoxVersionHEADRequest() (string, error) {
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

// ExtactFirefoxVersion extracts the Firefox version from the updater URL
func (t *TBDownloader) ExtractFirefoxVersion() (string, error) {
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
func (t *TBDownloader) NamePerPlatformFirefox(ietf string) string {
	extension := "tar.xz"
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

// FirefoxBrowserDirectory returns the path to the directory where the Firefox browser is installed.
func (t *TBDownloader) FirefoxBrowserDir() string {
	return filepath.Join(t.UnpackPath, "firefox_"+t.Lang)
}

// UnpackFirefox unpacks the Firefox package to the t.FirefoxBrowserDir()
func (t *TBDownloader) UnpackFirefox(binpath string) (string, error) {
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
	xzfile, err := os.Open(binpath)
	if err != nil {
		return "", fmt.Errorf("UnpackFirefox: XZFile error %s", err)
	}
	defer xzfile.Close()
	xzReader, err := xz.NewReader(xzfile)
	if err != nil {
		return "", fmt.Errorf("UnpackFirefox: XZReader error %s", err)
	}
	tarReader := tar.NewReader(xzReader)
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
			return "", fmt.Errorf("UnpackFirefox: Tar unpacker error %s", err)
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
func (t *TBDownloader) DownloadFirefoxUpdater() (string, string, error) {
	return t.DownloadFirefoxUpdaterForLang(t.Lang)
}

// DownloadFirefoxUpdaterForLang downloads the updater for the given language, overriding
// t.Lang. It returns the path to the downloaded updater and the downloaded
// detatched signature, or an error if one is encountered.
func (t *TBDownloader) DownloadFirefoxUpdaterForLang(ietf string) (string, string, error) {
	binary, sig, err := t.GetFirefoxUpdaterForLang(ietf)
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdater: %s", err)
	}
	sigpath := ""
	if t.OS == "linux" {
		sigpath, err = t.SingleFileDownload(sig, t.NamePerPlatformFirefox(ietf)+".asc")
		if err != nil {
			return "", "", fmt.Errorf("DownloadUpdater: %s", err)
		}
	}
	binpath, err := t.SingleFileDownload(binary, t.NamePerPlatformFirefox(ietf))
	if err != nil {
		return "", "", fmt.Errorf("DownloadUpdater: %s", err)
	}
	return binpath, sigpath, nil
}

// CheckSignature checks the signature of the updater.
// it returns an error if one is encountered. If not, it
// runs the updater and returns an error if one is encountered.
func (t *TBDownloader) CheckFirefoxSignature(binpath, sigpath string) (string, error) {
	if t.OS == "linux" {
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
	}
	err := fmt.Errorf("signature check failed")
	return "", fmt.Errorf("CheckSignature: %s", err)
}

// BoolCheckFirefoxSignature turns CheckFirefoxSignature into a bool.
func (t *TBDownloader) BoolCheckFirefoxSignature(binpath, sigpath string) bool {
	_, err := t.CheckFirefoxSignature(binpath, sigpath)
	return err == nil
}