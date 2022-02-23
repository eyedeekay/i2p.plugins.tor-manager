package tbget

import (
	"fmt"
	"io/ioutil"
	"log"
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
			os.Remove(fp)
			log.Println("Generating torrent for", fp)
			meta, err := t.GenerateTorrent(af, nil)
			if err != nil {
				return err
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
		log.Println("Copying", af, "to", sf)
		cp.Copy(af, sf)
		log.Println("Copying", fp, "to", sfp)
		cp.Copy(fp, sfp)
	}
	return nil
}

func (t *TBDownloader) GenerateTorrent(file string, announces []string) (*metainfo.MetaInfo, error) {
	info, err := metainfo.NewInfoFromFilePath(file, 5120)
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
	var url *url.URL
	if t.listener != nil {
		url, err = url.Parse("http://" + t.listener.Addr().(i2pkeys.I2PAddr).Base32() + "/" + filepath.Base(file))
		if err != nil {
			return nil, fmt.Errorf("GenerateTorrent: %s", err)
		}
		if t.Mirror != "" {
			mi.URLList = []string{url.String()}
		}
	}

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
