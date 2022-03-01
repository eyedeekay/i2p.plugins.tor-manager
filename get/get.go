package tbget

import (
	"archive/tar"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/dustin/go-humanize"
	sam "github.com/eyedeekay/sam3/helper"
	"github.com/itchio/damage"
	"github.com/itchio/damage/hdiutil"
	"github.com/itchio/headway/state"
	"github.com/magisterquis/connectproxy"
	"github.com/ulikunitz/xz"

	"golang.org/x/net/proxy"
)

// WORKING_DIR is the working directory for the application.
var WORKING_DIR = ""

// DefaultDir returns the default directory for the application.
func DefaultDir() string {
	if WORKING_DIR == "" {
		WORKING_DIR, _ = os.Getwd()
	}
	if !FileExists(WORKING_DIR) {
		os.MkdirAll(WORKING_DIR, 0755)
	}
	wd, err := filepath.Abs(WORKING_DIR)
	if err != nil {
		log.Fatal(err)
	}
	return wd
}

// UNPACK_PATH returns the path to the unpacked files.
func UNPACK_PATH() string {
	var UNPACK_PATH = filepath.Join(DefaultDir(), "unpack")
	return UNPACK_PATH
}

// DOWNLOAD_PATH returns the path to the downloads.
func DOWNLOAD_PATH() string {
	var DOWNLOAD_PATH = filepath.Join(DefaultDir(), "tor-browser")
	return DOWNLOAD_PATH
}

// TOR_UPDATES_URL is the URL of the Tor Browser update list.
const TOR_UPDATES_URL string = "https://aus1.torproject.org/torbrowser/update_3/release/downloads.json"

var (
	// DefaultIETFLang is the default language for the TBDownloader.
	DefaultIETFLang, _ = jibber_jabber.DetectIETF()
)

// TBDownloader is a struct which manages browser updates
type TBDownloader struct {
	UnpackPath   string
	DownloadPath string
	Lang         string
	OS, ARCH     string
	Mirror       string
	Verbose      bool
	Profile      *embed.FS
	listener     net.Listener
}

// OS is the operating system of the TBDownloader.
var OS = "linux"

// ARCH is the architecture of the TBDownloader.
var ARCH = "64"

// NewTBDownloader returns a new TBDownloader with the given language, using the TBDownloader's OS/ARCH pair
func NewTBDownloader(lang string, os, arch string, content *embed.FS) *TBDownloader {
	OS = os
	ARCH = arch
	return &TBDownloader{
		Lang:         lang,
		DownloadPath: DOWNLOAD_PATH(),
		UnpackPath:   UNPACK_PATH(),
		OS:           os,
		ARCH:         arch,
		Verbose:      false,
		Profile:      content,
	}
}

// ServeHTTP serves the DOWNLOAD_PATH as a mirror
func (t *TBDownloader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = path.Clean(r.URL.Path)
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

// Serve runs ServeHTTP on an I2P listener
func (t *TBDownloader) Serve() {
	var err error
	t.listener, err = sam.I2PListener("torbrowser-mirror", "127.0.0.1:7656", filepath.Join(t.UnpackPath, "torbrowser-mirror"))
	if err != nil {
		log.Fatal(err)
	}
	defer t.listener.Close()
	http.Serve(t.listener, t)
}

// GetRuntimePair returns the runtime.GOOS and runtime.GOARCH pair.
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
	if t.OS != "osx" {
		return fmt.Sprintf("%s%s", t.OS, t.ARCH)
	}
	return t.OS
}

// GetUpdater returns the updater for the given language, using the TBDownloader's OS/ARCH pair
// and only the defaults. It returns the URL of the updater and the detatched signature, or an error if one is not found.
func (t *TBDownloader) GetUpdater() (string, string, error) {
	return t.GetUpdaterForLang(t.Lang)
}

// GetUpdaterForLang returns the updater for the given language, using the TBDownloader's OS/ARCH pair
// it expects ietf to be a language. It returns the URL of the updater and the detatched signature, or an error if one is not found.
func (t *TBDownloader) GetUpdaterForLang(ietf string) (string, string, error) {
	jsonText, err := http.Get(TOR_UPDATES_URL)
	if err != nil {
		return "", "", fmt.Errorf("t.GetUpdaterForLang: %s", err)
	}
	defer jsonText.Body.Close()
	return t.GetUpdaterForLangFromJSON(jsonText.Body, ietf)
}

// GetUpdaterForLangFromJSON returns the updater for the given language, using the TBDownloader's OS/ARCH pair
// it expects body to be a valid json reader and ietf to be a language. It returns the URL of the updater and
// the detatched signature, or an error if one is not found.
func (t *TBDownloader) GetUpdaterForLangFromJSON(body io.ReadCloser, ietf string) (string, string, error) {
	jsonBytes, err := io.ReadAll(body)
	if err != nil {
		return "", "", fmt.Errorf("t.GetUpdaterForLangFromJSON: %s", err)
	}
	t.MakeTBDirectory()
	if err = ioutil.WriteFile(filepath.Join(t.DownloadPath, "downloads.json"), jsonBytes, 0644); err != nil {
		return "", "", fmt.Errorf("t.GetUpdaterForLangFromJSON: %s", err)
	}
	return t.GetUpdaterForLangFromJSONBytes(jsonBytes, ietf)
}

// Log logs things if Verbose is true.
func (t *TBDownloader) Log(function, message string) {
	if t.Verbose {
		log.Println(fmt.Sprintf("%s: %s", function, message))
	}
}

// MakeTBDirectory creates the tor-browser directory if it doesn't exist. It also unpacks a local copy of the TPO signing key.
func (t *TBDownloader) MakeTBDirectory() {
	os.MkdirAll(t.DownloadPath, 0755)

	empath := path.Join("tor-browser", "TPO-signing-key.pub")
	opath := filepath.Join(t.DownloadPath, "TPO-signing-key.pub")
	if !FileExists(opath) {
		t.Log("MakeTBDirectory()", "Initial TPO signing key not found, using the one embedded in the executable")
		bytes, err := t.Profile.ReadFile(empath)
		if err != nil {
			log.Fatal(err)
		}
		t.Log("MakeTBDirectory()", "Writing TPO signing key to disk")
		err = ioutil.WriteFile(opath, bytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
		t.Log("MakeTBDirectory()", "Writing TPO signing key to disk complete")
	}
	empath = path.Join("tor-browser", "unpack", "awo@eyedeekay.github.io.xpi")
	dpath := filepath.Join(t.DownloadPath, "awo@eyedeekay.github.io.xpi")
	opath = filepath.Join(t.UnpackPath, "awo@eyedeekay.github.io.xpi")
	if !FileExists(opath) {
		t.Log("MakeTBDirectory()", "Initial TAWO XPI not found, using the one embedded in the executable")
		bytes, err := t.Profile.ReadFile(empath)
		if err != nil {
			log.Fatal(err)
		}
		os.MkdirAll(filepath.Dir(dpath), 0755)
		os.MkdirAll(filepath.Dir(opath), 0755)
		t.Log("MakeTBDirectory()", "Writing AWO XPI to disk")
		err = ioutil.WriteFile(opath, bytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(dpath, bytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
		t.Log("MakeTBDirectory()", "Writing AWO XPI disk complete")
	}
}

// GetUpdaterForLangFromJSONBytes returns the updater for the given language, using the TBDownloader's OS/ARCH pair
// it expects jsonBytes to be a valid json string and ietf to be a language. It returns the URL of the updater and
// the detatched signature, or an error if one is not found.
func (t *TBDownloader) GetUpdaterForLangFromJSONBytes(jsonBytes []byte, ietf string) (string, string, error) {
	t.MakeTBDirectory()
	var dat map[string]interface{}
	t.Log("GetUpdaterForLangFromJSONBytes()", "Parsing JSON")
	if err := json.Unmarshal(jsonBytes, &dat); err != nil {
		return "", "", fmt.Errorf("func (t *TBDownloader)Name: %s", err)
	}
	t.Log("GetUpdaterForLangFromJSONBytes()", "Parsing JSON complete")
	if platform, ok := dat["downloads"]; ok {
		rtp := t.GetRuntimePair()
		if updater, ok := platform.(map[string]interface{})[rtp]; ok {
			if langUpdater, ok := updater.(map[string]interface{})[ietf]; ok {
				t.Log("GetUpdaterForLangFromJSONBytes()", "Found updater for language")
				bin := langUpdater.(map[string]interface{})["binary"].(string)
				sig := langUpdater.(map[string]interface{})["sig"].(string)
				return t.MirrorIze(bin), t.MirrorIze(sig), nil
			}
			// If we didn't find the language, try splitting at the hyphen
			lang := strings.Split(ietf, "-")[0]
			if langUpdater, ok := updater.(map[string]interface{})[lang]; ok {
				t.Log("GetUpdaterForLangFromJSONBytes()", "Found updater for backup language")
				bin := langUpdater.(map[string]interface{})["binary"].(string)
				sig := langUpdater.(map[string]interface{})["sig"].(string)
				return t.MirrorIze(bin), t.MirrorIze(sig), nil
			}
			// If we didn't find the language after splitting at the hyphen, try the default
			t.Log("GetUpdaterForLangFromJSONBytes()", "Last attempt, trying default language")
			return t.GetUpdaterForLangFromJSONBytes(jsonBytes, t.Lang)
		}
		return "", "", fmt.Errorf("t.GetUpdaterForLangFromJSONBytes: no updater for platform %s", rtp)
	}
	return "", "", fmt.Errorf("t.GetUpdaterForLangFromJSONBytes: %s", ietf)
}

func (t *TBDownloader) MirrorIze(replaceStr string) string {
	if t.Mirror != "" {
		return strings.Replace(replaceStr, "https://dist.torproject.org/torbrowser/", t.Mirror, 1)
	}
	return replaceStr
}

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// SingleFileDownload downloads a single file from the given URL to the given path.
// it returns the path to the downloaded file, or an error if one is encountered.
func (t *TBDownloader) SingleFileDownload(dl, name string) (string, error) {
	t.MakeTBDirectory()
	path := filepath.Join(t.DownloadPath, name)
	t.Log("SingleFileDownload()", fmt.Sprintf("Checking for updates %s to %s", dl, path))
	if !t.BotherToDownload(dl, name) {
		t.Log("SingleFileDownload()", "File already exists, skipping download")
		return path, nil
	}
	var d proxy.Dialer
	if t.MirrorIsI2P() {
		log.Println("Using I2P mirror, setting up proxy")
		var err error
		proxyURL, err := url.Parse("http://127.0.0.1:4444")
		if err != nil {
			panic(err)
		}
		d, err = connectproxy.New(proxyURL, proxy.Direct)
		if nil != err {
			panic(err)
		}
		tr := &http.Transport{
			Dial: d.Dial,
		}
		http.DefaultClient.Transport = tr
	}
	t.Log("SingleFileDownload()", "Downloading file")
	file, err := http.Get(dl)
	if err != nil {
		return "", fmt.Errorf("SingleFileDownload: %s", err)
	}
	defer file.Body.Close()
	outFile, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("SingleFileDownload: %s", err)
	}
	defer outFile.Close()
	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(outFile, io.TeeReader(file.Body, counter)); err != nil {
		return "", err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")
	//io.Copy(outFile, file.Body)
	t.Log("SingleFileDownload()", "Downloading file complete")
	return path, nil
}

// FileExists returns true if the given file exists. It will return true if used on an existing directory.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// BotherToDownload returns true if we need to download a file because we don't have an up-to-date
// version yet.
func (t *TBDownloader) BotherToDownload(dl, name string) bool {
	path := filepath.Join(t.DownloadPath, name)
	if !FileExists(path) {
		return true
	}
	defer ioutil.WriteFile(filepath.Join(t.DownloadPath, name+".last-url"), []byte(dl), 0644)
	lastURL, err := ioutil.ReadFile(filepath.Join(t.DownloadPath, name+".last-url"))
	if err != nil {
		return true
	}
	if string(lastURL) == dl {
		return false
	}
	return true
}

// NamePerPlatform returns the name of the updater for the given platform with appropriate extensions.
func (t *TBDownloader) NamePerPlatform(ietf string) string {
	extension := "tar.xz"
	windowsonly := ""
	switch t.OS {
	case "osx":
		extension = "dmg"
	case "win":
		windowsonly = "-installer"
		extension = "exe"
	}
	return fmt.Sprintf("torbrowser%s-%s-%s.%s", windowsonly, t.GetRuntimePair(), ietf, extension)
}

// DownloadUpdater downloads the updater for the t.Lang. It returns
// the path to the downloaded updater and the downloaded detatched signature,
// or an error if one is encountered.
func (t *TBDownloader) DownloadUpdater() (string, string, error) {
	return t.DownloadUpdaterForLang(t.Lang)
}

// DownloadUpdaterForLang downloads the updater for the given language, overriding
// t.Lang. It returns the path to the downloaded updater and the downloaded
// detatched signature, or an error if one is encountered.
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

// BrowserDir returns the path to the directory where the browser is installed.
func (t *TBDownloader) BrowserDir() string {
	return filepath.Join(t.UnpackPath, "tor-browser_"+t.Lang)
}

// UnpackUpdater unpacks the updater to the given path.
// it returns the path or an erorr if one is encountered.
func (t *TBDownloader) UnpackUpdater(binpath string) (string, error) {
	t.Log("UnpackUpdater()", fmt.Sprintf("Unpacking %s", binpath))
	if t.OS == "win" {
		installPath := t.BrowserDir()
		if !FileExists(installPath) {
			t.Log("UnpackUpdater()", "Windows updater, running silent NSIS installer")
			t.Log("UnpackUpdater()", fmt.Sprintf("Running %s %s %s", binpath, "/S", "/D="+installPath))
			cmd := exec.Command(binpath, "/S", "/D="+installPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return "", fmt.Errorf("UnpackUpdater: windows exec fail %s", err)
			}
		}
		return installPath, nil
	}
	if t.OS == "osx" {
		binpath = "tor-browser/torbrowser-osx64-en-US.dmg"
		log.Println("hdiutil", "mount", "\""+binpath+"\"")
		//cmd := exec.Command("open", "-W", "-n", "-a", "\""+binpath+"\"")
		//cmd := exec.Command("hdiutil", "attach", "\""+binpath+"\"")
		consumer := &state.Consumer{
			OnMessage: func(lvl string, msg string) {
				log.Printf("[%s] %s", lvl, msg)
			},
		}
		host := hdiutil.NewHost(consumer)
		if !FileExists(t.BrowserDir()) {
			_, err := damage.Mount(host, binpath, t.BrowserDir())
			if err != nil {
				return "", fmt.Errorf("UnpackUpdater: osx open/mount fail %s", err)
			}
		}
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//err := cmd.Run()
		//TODO: this might just need to be a hardcoded app path
		return t.BrowserDir(), nil
	}
	if FileExists(t.BrowserDir()) {
		return t.BrowserDir(), nil
	}
	fmt.Printf("Unpacking %s %s\n", binpath, t.UnpackPath)
	os.MkdirAll(t.UnpackPath, 0755)
	UNPACK_DIRECTORY, err := os.Open(t.UnpackPath)
	if err != nil {
		return "", fmt.Errorf("UnpackUpdater: directory error %s", err)
	}
	defer UNPACK_DIRECTORY.Close()
	xzfile, err := os.Open(binpath)
	if err != nil {
		return "", fmt.Errorf("UnpackUpdater: XZFile error %s", err)
	}
	defer xzfile.Close()
	xzReader, err := xz.NewReader(xzfile)
	if err != nil {
		return "", fmt.Errorf("UnpackUpdater: XZReader error %s", err)
	}
	tarReader := tar.NewReader(xzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("UnpackUpdater: Tar looper Error %s", err)
		}
		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(filepath.Join(UNPACK_DIRECTORY.Name(), header.Name), 0755)
			continue
		}
		filename := filepath.Join(UNPACK_DIRECTORY.Name(), header.Name)
		file, err := os.Create(filename)
		if err != nil {
			return "", fmt.Errorf("UnpackUpdater: Tar unpacker error %s", err)
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
	return t.BrowserDir(), nil
}

// CheckSignature checks the signature of the updater.
// it returns an error if one is encountered. If not, it
// runs the updater and returns an error if one is encountered.
func (t *TBDownloader) CheckSignature(binpath, sigpath string) (string, error) {
	pk := filepath.Join(t.DownloadPath, "TPO-signing-key.pub")
	var err error
	if err = Verify(pk, sigpath, binpath); err == nil {
		t.Log("CheckSignature: signature", "verified successfully")
		return t.UnpackUpdater(binpath)
	}
	return "", fmt.Errorf("CheckSignature: %s", err)
}

// BoolCheckSignature turns CheckSignature into a bool.
func (t *TBDownloader) BoolCheckSignature(binpath, sigpath string) bool {
	_, err := t.CheckSignature(binpath, sigpath)
	return err == nil
}

// TestHTTPDefaultProxy returns true if the I2P proxy is up or blocks until it is.
func TestHTTPDefaultProxy() bool {
	return TestHTTPProxy("127.0.0.1", "4444")
}

// Seconds increments the seconds and displays the number of seconds every 10 seconds
func Seconds(now int) int {
	time.Sleep(time.Second)
	if now == 3 {
		return 0
	}
	return now + 1
}

// TestHTTPBackupProxy returns true if the I2P backup proxy is up or blocks until it is.
func TestHTTPBackupProxy() bool {
	now := 0
	limit := 0
	for {
		_, err := net.Listen("tcp", "127.0.0.1:4444")
		if err != nil {
			log.Println("SAM HTTP proxy is open", err)
			return true
		} else {
			if now == 0 {
				log.Println("Waiting for HTTP Proxy", (10 - limit), "remaining attempts")
				limit++
			}
			now = Seconds(now)
		}
		if limit == 10 {
			break
		}
	}
	return false
}

// TestHTTPProxy returns true if the proxy at host:port is up or blocks until it is.
func TestHTTPProxy(host, port string) bool {
	now := 0
	limit := 0
	for {
		proxy := hTTPProxy(host, port)
		if proxy {
			return true
		} else {
			if now == 0 {
				log.Println("Waiting for HTTP Proxy", (10 - limit), "remaining attempts")
				limit++
			}
			now = Seconds(now)
		}
		if limit == 10 {
			break
		}
	}
	return false
}

func hTTPProxy(host, port string) bool {
	proxyURL, err := url.Parse("http://" + host + ":" + port)
	if err != nil {
		log.Panic(err)
	}
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	resp, err := myClient.Get("http://proxy.i2p/")
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return strings.Contains(string(body), "I2P HTTP proxy OK")
		}
	}
	return false
}

func (t *TBDownloader) MirrorIsI2P() bool {
	// check if hostname is an I2P hostname
	url, err := url.Parse(t.Mirror)
	if err != nil {
		return false
	}
	log.Println("Checking if", url.Hostname(), "is an I2P hostname")

	return strings.Contains(url.Hostname(), ".i2p")
}
