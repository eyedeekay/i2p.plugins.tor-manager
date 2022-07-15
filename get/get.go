package tbget

import (
	"archive/tar"
	"context"
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
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/cretz/bine/tor"
	"github.com/dustin/go-humanize"
	sam "github.com/eyedeekay/sam3/helper"
	"github.com/itchio/damage"
	"github.com/itchio/damage/hdiutil"
	"github.com/itchio/headway/state"
	"github.com/magisterquis/connectproxy"
	cp "github.com/otiai10/copy"
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

func Languages() []string {
	jsonText, err := http.Get(TOR_UPDATES_URL)
	if err != nil {
		return []string{}
	}
	defer jsonText.Body.Close()
	jsonBytes, err := ioutil.ReadAll(jsonText.Body)
	if err != nil {
		return []string{}
	}
	var updates map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &updates); err != nil {
		return []string{}
	}
	var languages []string
	//updates[]
	for i := range updates["downloads"].(map[string]interface{})["win64"].(map[string]interface{}) {
		languages = append(languages, i)
	}
	return languages
}

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
	NoUnpack     bool
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
		t.ARCH = "64"
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

	tpk := "TPO-signing-key.pub"
	if t.OS == "linux" && runtime.GOARCH == "arm64" {
		tpk = "NOT-TPO-signing-key.pub"
	}

	empath := path.Join("tor-browser", tpk)
	opath := filepath.Join(t.DownloadPath, tpk)
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
	log.Println("MirrorIze()", "Replacing", replaceStr, t.Mirror)
	if t.OS == "linux" && runtime.GOARCH == "arm64" {
		replaceStr = strings.Replace(replaceStr, "linux64", "linux-arm64", -1)
		if strings.HasSuffix(replaceStr, ".tar.xz.asc") {
			//sha256sums-unsigned-build.txt.asc
			lastElement := filepath.Base(
				strings.Replace(replaceStr, "https://", strings.Replace(replaceStr, "http://", "", 1), 1),
			)
			replaceStr = strings.Replace(replaceStr, lastElement, "sha256sums-unsigned-build.txt.asc", -1)
		}
	}
	if strings.Contains(t.Mirror, "i2psnark") {
		replaceStr = strings.Replace(replaceStr, "https://dist.torproject.org/torbrowser/", t.Mirror, 1)
		dpath := filepath.Base(replaceStr)
		replaceStr = strings.Replace(replaceStr, "http://", "", 1)
		replaceStr = filepath.Dir(replaceStr)
		replaceStr = filepath.Dir(replaceStr)
		newurl := "http://" + filepath.Join(replaceStr, dpath)
		log.Println("MirrorIze()", "Final URL", newurl)
		return newurl
	}
	if t.Mirror != "" {
		return strings.Replace(replaceStr, "https://dist.torproject.org/torbrowser/", t.Mirror, 1)
	}
	log.Println("MirrorIze()", "Final URL", replaceStr)
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
	fmt.Fprintf(os.Stderr, "\r%s", strings.Repeat(" ", 35))
	fmt.Fprintf(os.Stderr, "\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func (t *TBDownloader) StartConf() *tor.StartConf {
	return StartConf(t.TorPath())
}

func StartConf(tp string) *tor.StartConf {
	paths := []string{
		"/bin/tor",
		"/usr/bin/tor",
		"/usr/sbin/tor",
		"/usr/local/bin/tor",
		"/usr/bin/tor",
	}
	path := strings.Split(os.Getenv("PATH"), ":")
	for _, p := range path {
		p := filepath.Join(p, "tor")
		paths = append(paths, p)
	}
	for _, path := range paths {
		if FileExists(path) {
			return &tor.StartConf{
				ExePath:           path,
				RetainTempDataDir: false,
			}
		}
	}
	if FileExists(tp) {
		return &tor.StartConf{
			ExePath:           tp,
			RetainTempDataDir: false,
		}
	}
	return nil
}

// SetupProxy sets up the proxy for the given URL
func (t *TBDownloader) SetupProxy() error {
	return SetupProxy(t.Mirror, t.TorPath())
}

func unSetupProxy() {
	http.DefaultClient.Transport = nil
}

var t *tor.Tor

func SetupProxy(mirror, tp string) error {
	var d proxy.Dialer
	http.DefaultClient.Transport = nil
	defer unSetupProxy()
	if MirrorIsI2P(mirror) {
		log.Println("Using I2P mirror, setting up proxy")
		var err error
		proxyURL, err := url.Parse("http://127.0.0.1:4444")
		if err != nil {
			return err
		}
		d, err = connectproxy.New(proxyURL, proxy.Direct)
		if nil != err {
			return err
		}
		tr := &http.Transport{
			Dial: d.Dial,
		}
		http.DefaultClient.Transport = tr
	} else {
		nut := os.Getenv("TOR_MANAGER_NEVER_USE_TOR")
		if nut != "true" {
			if !strings.Contains(mirror, "127.0.0.1") && !strings.Contains(mirror, "localhost") {
				if tmp, torerr := net.Listen("tcp", "127.0.0.1:9050"); torerr != nil {
					log.Println("System Tor is running, downloading over that because obviously.")
					is_flatpak := os.Getenv("APP_ID") != ""
					if is_flatpak {
						log.Println("Flatpak detected, using Tor without bine")
						url_i := url.URL{}
						url_proxy, err := url_i.Parse("socks5://127.0.0.1:9050")
						if err != nil {
							return err
						}

						tr := &http.Transport{}
						tr.Proxy = http.ProxyURL(url_proxy) // set proxy
						http.DefaultClient.Transport = tr
						return nil
					}
					var err error
					if t == nil {
						t, err = tor.Start(context.Background(), StartConf(tp))
						if err != nil {
							if t == nil {
								return err
							}
						}
					}
					//defer t.Close()
					// Wait at most a minute to start network and get
					dialCtx, _ := context.WithTimeout(context.Background(), time.Minute)
					//defer dialCancel()
					// Make connection
					dialer, err := t.Dialer(dialCtx, nil)
					if err != nil {
						return err
					}
					tr := &http.Transport{DialContext: dialer.DialContext}
					http.DefaultClient.Transport = tr
				} else {
					tmp.Close()
				}
			}
		}
	}
	return nil
}

// SingleFileDownload downloads a single file from the given URL to the given path.
// it returns the path to the downloaded file, or an error if one is encountered.
func (t *TBDownloader) SingleFileDownload(dl, name string, rangebottom int64) (string, error) {
	t.MakeTBDirectory()
	path := filepath.Join(t.DownloadPath, name)
	if filepath.IsAbs(name) {
		path = name
	}

	t.Log("SingleFileDownload()", fmt.Sprintf("Checking for updates %s to %s", dl, path))
	if !t.BotherToDownload(dl, name) {
		t.Log("SingleFileDownload()", "File already exists, skipping download")
		return path, nil
	}
	err := t.SetupProxy()
	if err != nil {
		return "", err
	}
	dlurl, err := url.Parse(dl)
	if err != nil {
		return "", err
	}
	if FileExists(path) {
		size, err := os.Stat(path)
		if err != nil {
			return "", err
		}
		rangebottom = size.Size()
		t.Log("SingleFileDownload()", fmt.Sprintf("Resuming download from %d", rangebottom))
	}
	req := http.Request{
		Method: "GET",
		URL:    dlurl,
		Header: http.Header{
			"Range": []string{fmt.Sprintf("bytes=%d-", rangebottom)},
		},
	}
	t.Log("SingleFileDownload()", "Downloading file "+dl)
	//file, err := http.Get(dl)
	file, err := http.DefaultClient.Do(&req)
	//Do(&req, nil)
	if err != nil {
		return "", fmt.Errorf("SingleFileDownload: Request Error %s", err)
	}
	defer file.Body.Close()
	outFile, err := Create(path)
	if err != nil {
		return "", fmt.Errorf("SingleFileDownload: Write Error %s", err)
	}
	defer outFile.Close()
	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{
		Total: uint64(rangebottom),
	}
	if rangebottom, err := io.Copy(outFile, io.TeeReader(file.Body, counter)); err != nil {
		return t.SingleFileDownload(dl, name, rangebottom)
		//"", err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")
	//io.Copy(outFile, file.Body)
	t.Log("SingleFileDownload()", "Downloading file complete")
	return path, nil
}

func Create(path string) (*os.File, error) {
	if FileExists(path) {
		stat, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		return os.OpenFile(path, os.O_APPEND|os.O_WRONLY, stat.Mode().Perm())
	}
	// Create the file
	outFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return outFile, nil
}

// FileExists returns true if the given file exists. It will return true if used on an existing directory.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (t *TBDownloader) FetchContentLength(dl, name string) (int64, error) {
	t.MakeTBDirectory()
	return FetchContentLength(dl, name)
}
func FetchContentLength(dl, name string) (int64, error) {
	log.Println("FetchContentLength():", fmt.Sprintf("Checking for updates %s to %s", dl, name))
	err := SetupProxy(dl, "")
	if err != nil {
		return 0, err
	}
	dlurl, err := url.Parse(dl)
	if err != nil {
		return 0, err
	}
	req := http.Request{
		Method: "HEAD",
		URL:    dlurl,
	}
	log.Println("FetchContentLength()", "Downloading file "+dl)
	//file, err := http.Get(dl)
	file, err := http.DefaultClient.Do(&req)
	//Do(&req, nil)
	if err != nil {
		return 0, fmt.Errorf("FetchContentLength: Request Error %s", err)
	}
	file.Body.Close()
	log.Println("Content-Length:", file.ContentLength)
	return file.ContentLength, nil
}

// BotherToDownload returns true if we need to download a file because we don't have an up-to-date
// version yet.
func (t *TBDownloader) BotherToDownload(dl, name string) bool {
	path := filepath.Join(t.DownloadPath, name)
	if !FileExists(path) {
		return true
	}
	stat, err := os.Stat(path)
	if err != nil {
		return true
	}
	// 86 MB
	if !strings.Contains(name, ".asc") {
		contentLength, err := t.FetchContentLength(dl, name)
		if err != nil {
			return true
		}

		l := 4
		if len(strconv.Itoa(int(contentLength))) < 4 {
			return true
		}
		lenString := strconv.Itoa(int(contentLength))[:l]
		lenSize := strconv.Itoa(int(stat.Size()))[:l]
		fmt.Fprintf(os.Stderr, "comparing sizes: %v %v", lenString, lenSize)

		if stat.Size() == contentLength {
			//if lenString != lenSize {
			//	return true
			//} else {
			fmt.Fprintf(os.Stderr, "BotherToDownload(): %s is fully downloaded\n", name)
			return false
			//}
		}
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
func (t *TBDownloader) NamePerPlatform(ietf, version string) string {
	extension := "tar.xz"
	windowsonly := ""
	switch t.OS {
	case "osx":
		extension = "dmg"
	case "win":
		windowsonly = "-installer"
		extension = "exe"
	}
	//version, err := t.Get
	return fmt.Sprintf("tor-browser%s-%s-%s_%s.%s", windowsonly, t.GetRuntimePair(), version, ietf, extension)
}

func (t *TBDownloader) GetVersion() string {
	binary, _, err := t.GetUpdaterForLang(t.Lang)
	if err != nil {
		return ""
	}
	version := strings.Split(binary, "/")[len(strings.Split(binary, "/"))-2]
	return version
}

func (t *TBDownloader) GetName() string {
	return t.NamePerPlatform(t.Lang, t.GetVersion())
}

// DownloadUpdater downloads the updater for the t.Lang. It returns
// the path to the downloaded updater and the downloaded detatched signature,
// or an error if one is encountered.
func (t *TBDownloader) DownloadUpdater() (string, string, string, error) {
	return t.DownloadUpdaterForLang(t.Lang)
}

// DownloadUpdaterForLang downloads the updater for the given language, overriding
// t.Lang. It returns the path to the downloaded updater and the downloaded
// detatched signature, or an error if one is encountered.
func (t *TBDownloader) DownloadUpdaterForLang(ietf string) (string, string, string, error) {
	binary, sig, err := t.GetUpdaterForLang(ietf)
	if err != nil {
		return "", "", "", fmt.Errorf("DownloadUpdaterForLang: %s", err)
	}
	version := t.GetVersion()
	if strings.Contains(t.Mirror, "i2psnark") {
		if !TorrentDownloaded(ietf, t.GetRuntimePair()) {
			t.Log("DownloadUpdaterForLang()", "Downloading torrent")
			SetupProxy("http://idk.i2p/", "")
			//Download the torrent files from their static locations.
			i2psnark, err := FindSnarkDirectory()
			if err != nil {
				return "", "", "", err
			}
			log.Println("Downloading torrent from", i2psnark)
			asctorrent := filepath.Join(t.NamePerPlatform(ietf, version) + ".asc" + ".torrent")
			fmt.Println("Downloading", asctorrent)
			_, err = t.SingleFileDownload("http://idk.i2p/torbrowser/"+asctorrent, filepath.Join(i2psnark, asctorrent), 0)
			if err != nil {
				return "", "", "", fmt.Errorf("DownloadUpdaterForLang: %s", err)
			}
			bintorrent := filepath.Join(t.NamePerPlatform(ietf, version) + ".torrent")
			fmt.Println("Downloading", bintorrent)
			_, err = t.SingleFileDownload("http://idk.i2p/torbrowser/"+bintorrent, filepath.Join(i2psnark, bintorrent), 0)
			if err != nil {
				return "", "", "", fmt.Errorf("DownloadUpdaterForLang: %s", err)
			}
		}
		for !TorrentDownloaded(ietf, t.GetRuntimePair()) {
			log.Println("DownloadUpdaterForLang:", "Waiting for torrent to download")
			time.Sleep(time.Second * 10)
		}
		time.Sleep(time.Second * 10)
	}

	sigpath, err := t.SingleFileDownload(sig, t.NamePerPlatform(ietf, version)+".asc", 0)
	if err != nil {
		return "", "", "", fmt.Errorf("DownloadUpdaterForLang: %s", err)
	}
	binpath, err := t.SingleFileDownload(binary, t.NamePerPlatform(ietf, version), 0)
	if err != nil {
		return "", sigpath, "", fmt.Errorf("DownloadUpdaterForLang: %s", err)
	}
	var sumpath string
	if t.OS == "linux" && runtime.GOARCH == "arm64" {
		sumpath, err = t.SingleFileDownload("https://sourceforge.net/projects/tor-browser-ports/files/11.0.6/sha256sums-unsigned-build.txt/download", t.NamePerPlatform(ietf, version)+".sha256sums", 0)
		if err != nil {
			return "", sigpath, sumpath, fmt.Errorf("DownloadUpdaterForLang: %s", err)
		}
	}
	return binpath, sigpath, sumpath, nil
}

// BrowserDir returns the path to the directory where the browser is installed.
func (t *TBDownloader) BrowserDir() string {
	return filepath.Join(t.UnpackPath, "tor-browser_"+t.Lang)
}

func (t *TBDownloader) I2PBrowserDir() string {
	return filepath.Join(t.UnpackPath, "i2p-browser_"+t.Lang)
}

// UnpackUpdater unpacks the updater to the given path.
// it returns the path or an erorr if one is encountered.
func (t *TBDownloader) UnpackUpdater(binpath string) (string, error) {
	if t.NoUnpack {
		return binpath, nil
	}
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
			if err := cp.Copy(t.BrowserDir(), t.I2PBrowserDir()); err != nil {
				return "", fmt.Errorf("UnpackUpdater: copy fail %s", err)
			}
		}
		// copy BrowserDir() to I2PBrowserDir()

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
			if _, err := damage.Mount(host, binpath, t.BrowserDir()); err != nil {
				return "", fmt.Errorf("UnpackUpdater: osx open/mount fail %s", err)
			}
		}
		if !FileExists(t.I2PBrowserDir()) {
			if _, err := damage.Mount(host, binpath, t.I2PBrowserDir()); err != nil {
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
		if !FileExists(t.I2PBrowserDir()) {
			if err := cp.Copy(t.BrowserDir(), t.I2PBrowserDir()); err != nil {
				return "", fmt.Errorf("UnpackUpdater: copy fail %s", err)
			}
		}
		return t.BrowserDir(), nil
	}
	fmt.Fprintf(os.Stderr, "Unpacking %s %s\n", binpath, t.UnpackPath)
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
			fmt.Fprintf(os.Stderr, "Unpacked %s\n", header.Name)
		}
	}
	if !FileExists(t.I2PBrowserDir()) {
		if err := cp.Copy(t.BrowserDir(), t.I2PBrowserDir()); err != nil {
			return "", fmt.Errorf("UnpackUpdater: copy fail %s", err)
		}
	}
	return t.BrowserDir(), nil
}

// TorPath returns the path to the Tor executable
func (s *TBDownloader) TorPath() string {
	if s.OS == "osx" {
		return filepath.Join(s.UnpackPath, "Tor Browser.app", "Contents", "Resources", "TorBrowser", "Tor", "tor")
	}
	return filepath.Join(s.UnpackPath, "Browser", "TorBrowser", "Tor", "tor")
}

// CheckSignature checks the signature of the updater.
// it returns an error if one is encountered. If not, it
// runs the updater and returns an error if one is encountered.
func (t *TBDownloader) CheckSignature(binpath, sigpath string) (string, error) {
	pk := filepath.Join(t.DownloadPath, "TPO-signing-key.pub")
	if t.OS == "linux" && runtime.GOARCH == "arm64" {
		pk = filepath.Join(t.DownloadPath, "NOT-TPO-signing-key.pub")
	}
	var err error
	if err = Verify(pk, sigpath, binpath); err == nil {
		log.Println("CheckSignature: signature", "verified successfully")
		if !t.NoUnpack {
			return t.UnpackUpdater(binpath)
		}
		log.Printf("CheckSignature: %s", "NoUnpack set, skipping unpack")
		return t.BrowserDir(), nil
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
	return MirrorIsI2P(t.Mirror)
}

func MirrorIsI2P(mirror string) bool {
	// check if hostname is an I2P hostname
	url, err := url.Parse(mirror)
	if err != nil {
		return false
	}
	log.Println("Checking if", url.Hostname(), "is an I2P hostname")

	return strings.Contains(url.Hostname(), ".i2p")
}
