package tbget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/eyedeekay/sam3/i2pkeys"
	cp "github.com/otiai10/copy"
	"github.com/xgfone/bt/bencode"
	"github.com/xgfone/bt/metainfo"
)

func (t *TBDownloader) DownloadedFilesList() ([]string, error) {
	files, err := ioutil.ReadDir(t.DownloadPath)
	if err != nil {
		return nil, fmt.Errorf("DownloadedFilesList: %s", err)
	}
	var list []string
	for _, f := range files {
		list = append(list, f.Name())
	}
	return list, nil

}

func (t *TBDownloader) GenerateMissingTorrents() error {
	files, err := t.DownloadedFilesList()
	if err != nil {
		return err
	}
	for _, f := range files {
		fp := filepath.Join(t.DownloadPath, f+".torrent")
		af := filepath.Join(t.DownloadPath, f)
		if !strings.HasSuffix(af, ".torrent") {
			//os.Remove(fp)
			if !FileExists(fp) {
				log.Println("Generating torrent for", fp)
				meta, err := t.GenerateTorrent(af, nil)
				if err != nil {
					//return err
					log.Println("GenerateMissingTorrents:", err)
					continue
				}
				file, err := os.Create(fp)
				if err != nil {
					return err
				}
				meta.Write(file)
				file.Close()
			}
			snark, err := FindSnarkDirectory()
			if err != nil {
				return err
			}
			sf := filepath.Join(snark, f)
			sfp := filepath.Join(snark, f+".torrent")
			if !FileExists(sf) {
				log.Println("Copying", af, "to", sf)
				cp.Copy(af, sf)
			}
			if !FileExists(sfp) {
				log.Println("Copying", fp, "to", sfp)
				cp.Copy(fp, sfp)
			}
		}
	}
	return nil
}

func (t *TBDownloader) GenerateTorrent(file string, announces []string) (*metainfo.MetaInfo, error) {

	//info, err := metainfo.NewInfoFromFilePath(file, 5120)
	info, err := metainfo.NewInfoFromFilePath(file, 10240)
	if err != nil {
		return nil, fmt.Errorf("GenerateTorrent: %s", err)
	}
	info.Name = filepath.Base(file)

	var mi metainfo.MetaInfo
	mi.InfoBytes, err = bencode.EncodeBytes(info)
	if err != nil {
		return nil, fmt.Errorf("GenerateTorrent: %s", err)
	}

	switch len(announces) {
	case 0:
		mi.Announce = "http://mb5ir7klpc2tj6ha3xhmrs3mseqvanauciuoiamx2mmzujvg67uq.b32.i2p/a"
	case 1:
		mi.Announce = announces[0]
	default:
		mi.AnnounceList = metainfo.AnnounceList{announces}
	}
	url, err := url.Parse("http://idk.i2p/torbrowser/" + filepath.Base(file))
	if err != nil {
		return nil, fmt.Errorf("GenerateTorrent: %s", err)
	}
	mi.URLList = []string{url.String()}
	if t.listener != nil {
		url, err := url.Parse("http://" + t.listener.Addr().(i2pkeys.I2PAddr).Base32() + "/" + filepath.Base(file))
		if err != nil {
			return nil, fmt.Errorf("GenerateTorrent: %s", err)
		}
		if t.Mirror != "" {
			mi.URLList = append(mi.URLList, url.String())
		}
	}
	clearurl, err := url.Parse("https://eyedeekay.github.io/torbrowser/" + filepath.Base(file))
	if err != nil {
		return nil, fmt.Errorf("GenerateTorrent: %s", err)
	}
	mi.URLList = append(mi.URLList, clearurl.String())
	return &mi, nil
}

func FindSnarkDirectory() (string, error) {
	// Snark could be at:
	// or: $I2P_CONFIG/i2psnark/
	// or: $I2P/i2psnark/
	// or: $HOME/.i2p/i2psnark/
	// or: /var/lib/i2p/i2p-config/i2psnark/
	// or: %LOCALAPPDATA\i2p\i2psnark\
	// or: %APPDATA\i2p\i2psnark\

	SNARK_CONFIG := os.Getenv("SNARK_CONFIG")
	if SNARK_CONFIG != "" {
		checkfori2pcustom := filepath.Join(SNARK_CONFIG)
		if FileExists(checkfori2pcustom) {
			return checkfori2pcustom, nil
		}
	}

	I2P_CONFIG := os.Getenv("I2P_CONFIG")
	if I2P_CONFIG != "" {
		checkfori2pcustom := filepath.Join(I2P_CONFIG, "i2psnark")
		if FileExists(checkfori2pcustom) {
			return checkfori2pcustom, nil
		}
	}

	I2P := os.Getenv("I2P")
	if I2P != "" {
		checkfori2p := filepath.Join(I2P, "i2psnark")
		if FileExists(checkfori2p) {
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
		checkfori2plocal := filepath.Join(home, "AppData", "Local", "i2p", "i2psnark")
		if FileExists(checkfori2plocal) {
			return checkfori2plocal, nil
		}
		checkfori2proaming := filepath.Join(home, "AppData", "Roaming", "i2p", "i2psnark")
		if FileExists(checkfori2proaming) {
			return checkfori2proaming, nil
		}
	case "linux":
		checkfori2phome := filepath.Join(home, ".i2p", "i2psnark")
		if FileExists(checkfori2phome) {
			return checkfori2phome, nil
		}
		checkfori2pservice := filepath.Join("/var/lib/i2p/i2p-config", "i2psnark")
		if FileExists(checkfori2pservice) {
			return checkfori2pservice, nil
		}
	case "darwin":
		return "", fmt.Errorf("FindSnarkDirectory: Automatic torrent generation is not supported on MacOS, for now copy the files manually")
	}
	return "", fmt.Errorf("FindSnarkDirectory: Unable to find snark directory")
}

func TorrentReady() bool {
	if _, err := FindSnarkDirectory(); err != nil {
		return false
	}
	return true
}

func TorrentPath() (string, string) {
	extension := "tar.xz"
	windowsonly := ""
	switch runtime.GOOS {
	case "darwin":
		extension = "dmg"
	case "windows":
		windowsonly = "-installer"
		extension = "exe"
	}
	//version, err := t.Get
	return fmt.Sprintf("tor-browser%s", windowsonly), extension
}

func GetTorBrowserVersionFromUpdateURL() (string, error) {
	// download the json file from TOR_UPDATES_URL
	// parse the json file to get the latest version
	// return the latest version
	err := SetupProxy(TOR_UPDATES_URL, "")
	if err != nil {
		return "", err
	}
	resp, err := http.Get(TOR_UPDATES_URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var updates map[string]interface{}
	err = json.Unmarshal(body, &updates)
	if err != nil {
		return "", err
	}
	latest := updates["downloads"].(map[string]interface{})["linux64"].(map[string]interface{})["en-US"].(map[string]interface{})

	for key, value := range latest {
		if key == "binary" {
			log.Printf("%s: %s\n", key, value)
			url, err := url.Parse(value.(string))
			if err != nil {
				return "", err
			}
			spl := strings.Split(url.Path, "/")
			return spl[len(spl)-1], nil
		}
	}

	return "11.0.9", nil
}

func TorrentDownloaded() bool {
	version, err := GetTorBrowserVersionFromUpdateURL()
	if err != nil {
		panic(err)
	}
	log.Println("Tor Browser Version", version)

	cmpsize := 8661000
	found := false
	if dir, err := FindSnarkDirectory(); err == nil {
		err := filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				prefix, suffix := TorrentPath()
				path = filepath.Base(path)
				if strings.HasPrefix(path, prefix) && strings.Contains(path, version) && strings.HasSuffix(path, suffix) {
					if !strings.HasSuffix(path, ".torrent") {
						if info.Size() > int64(cmpsize) {
							fmt.Println("TorrentDownloaded: Torrent Download found:", path, info.Size(), int64(cmpsize))
							found = true
							return nil
						} else {
							return fmt.Errorf("TorrentDownloaded: Torrent Download found but size is too small: %s", path)
						}
					}
				}
				return nil
			})
		if found {
			return err == nil
		}
		return false
	}
	return false
}

func Torrent() bool {
	if !TorrentReady() {
		return false
	}
	if !TorrentDownloaded() {
		return false
	}
	return true
}
