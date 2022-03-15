package tbserve

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/justinas/nosurf"
	cp "github.com/otiai10/copy"
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
	i2pdotonion "i2pgit.org/idk/i2p.plugins.tor-manager/onion"
	TBSupervise "i2pgit.org/idk/i2p.plugins.tor-manager/supervise"
)

// Client manages and supervises a Tor Browser instance.
type Client struct {
	hostname string
	TBD      *tbget.TBDownloader
	TBS      *TBSupervise.Supervisor
	Onion    *i2pdotonion.I2POnionService
	DarkMode bool
	Host     string
	Port     int
	server   http.Server
}

// NewClient creates a new Client.
func NewClient(verbose bool, lang, OS, arch, mirror string, content *embed.FS) (*Client, error) {
	m := &Client{
		TBD: tbget.NewTBDownloader(lang, OS, arch, content),
	}
	m.TBD.Mirror = mirror
	m.TBD.Verbose = verbose
	m.TBD.MakeTBDirectory()
	var err error
	m.Onion, err = i2pdotonion.NewOnionService(m.TBD.DownloadPath)
	if err != nil {
		return nil, err
	}
	tgz, sig, sums, err := m.TBD.DownloadUpdaterForLang(lang)
	if err != nil {
		panic(err)
	}
	sum := ""
	if sums != "" && runtime.GOOS == "linux" && runtime.GOARCH == "arm64" {
		b, err := ioutil.ReadFile(sums)
		if err != nil {
			log.Fatal(err)
		}
		// find the line containing the checksum of our tgz file
		for _, line := range strings.Split(string(b), "\n") {
			if strings.Contains(line, lang+".tar.xz") {
				sum = strings.Split(line, " ")[0]
				break
			}
		}
		log.Println("Checksum for ARM:" + sum)
		// compute the sha256sum of the downloaded tar.xz file
		f, err := os.Open(tgz)
		if err != nil {
			log.Fatal(err)
		}
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}
		f.Close()
		if sum != hex.EncodeToString(h.Sum(nil)) {
			log.Fatal("Checksum mismatch")
		}
		var home string
		if home, err = m.TBD.CheckSignature(sums, sig); err != nil {
			log.Fatal(err)
		} else {
			_, err = m.TBD.UnpackUpdater(tgz)
			if err != nil {
				return nil, fmt.Errorf("unpacking updater: %v", err)
			}
			log.Printf("Signature check passed: %s %s", tgz, sig)
		}
		m.TBS = TBSupervise.NewSupervisor(home, lang)
		go m.TBS.RunTorWithLang()
		return m, nil
	}
	var home string
	if home, err = m.TBD.CheckSignature(tgz, sig); err != nil {
		log.Fatal(err)
	} else {
		_, err = m.TBD.UnpackUpdater(tgz)
		if err != nil {
			return nil, fmt.Errorf("unpacking updater: %v", err)
		}
		log.Printf("Signature check passed: %s %s", tgz, sig)
	}
	m.TBS = TBSupervise.NewSupervisor(home, lang)
	go m.TBS.RunTorWithLang()
	return m, nil
}

// NewFirefoxClient creates a new Client.
func NewFirefoxClient(verbose bool, lang, os, arch, mirror string, content *embed.FS) (*Client, error) {
	m := &Client{
		TBD: tbget.NewFirefoxDownloader(lang, os, arch, content),
	}
	m.TBD.Mirror = mirror
	m.TBD.Verbose = verbose
	m.TBD.MakeTBDirectory()
	tgz, sig, err := m.TBD.DownloadFirefoxUpdaterForLang(lang)
	if err != nil {
		panic(err)
	}
	var home string
	if home, err = m.TBD.CheckFirefoxSignature(tgz, sig); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Signature check passed: %s %s", tgz, sig)
	}
	m.TBS = TBSupervise.NewSupervisor(home, lang)
	return m, nil
}

// GetHost returns the hostname of the client.
func (m *Client) GetHost() string {
	if m.Host == "" {
		m.Host = "127.0.0.1"
	}
	return m.Host
}

// GetPort returns the port of the client.
func (m *Client) GetPort() string {
	if m.Port == 0 {
		m.Port = 7695
	}
	return strconv.Itoa(m.Port)
}

// GetAddress returns the address of the client.
func (m *Client) GetAddress() string {
	return m.GetHost() + ":" + m.GetPort()
}

// ServeHTTP handles HTTP requests.
func (m *Client) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	path := path.Clean(rq.URL.Path)
	rq.URL.Path = path
	log.Printf("ServeHTTP: '%s'", path)
	fileextension := filepath.Ext(path)
	switch fileextension {
	case ".json":
		m.serveJSON(rw, rq)
		return
	case ".css":
		m.serveCSS(rw, rq)
		return
	case ".js":
		m.serveJS(rw, rq)
		return
	case ".png":
		m.servePNG(rw, rq)
		return
	case ".ico":
		m.serveICO(rw, rq)
		return
	case ".svg":
		m.serveSVG(rw, rq)
		return
	default:
		switch path {
		case "/launch-tor-browser":
			log.Println("Starting Tor Browser")
			go m.TBS.RunTBWithLang()
			http.Redirect(rw, rq, "/", http.StatusFound)
		case "/launch-i2p-browser":
			log.Println("Starting I2P Browser")
			go m.TBS.RunI2PBWithLang()
			http.Redirect(rw, rq, "/", http.StatusFound)
		case "/start-tor":
			log.Println("Starting Tor")
			go m.TBS.RunTorWithLang()
			http.Redirect(rw, rq, "/", http.StatusFound)
		case "/stop-tor":
			log.Println("Stopping Tor")
			go m.TBS.StopTor()
			http.Redirect(rw, rq, "/", http.StatusFound)
		case "/switch-theme":
			log.Println("Switching theme")
			m.DarkMode = !m.DarkMode
			http.Redirect(rw, rq, "/", http.StatusFound)
		default:
			b, _ := m.Page()
			rw.Header().Set("Content-Type", "text/html")
			rw.Write([]byte(b))
		}
	}

}

// Serve serve the control panel locally
func (m *Client) Serve() error {
	//http.Handle("/", m)
	mirrorjson, err := m.GenerateMirrorJSON()
	if err != nil {
		return err
	}
	m.server = http.Server{
		Addr:    m.GetAddress(),
		Handler: nosurf.New(m),
	}
	ioutil.WriteFile(filepath.Join(m.TBD.DownloadPath, "mirror.json"), []byte(mirrorjson), 0644)
	cp.Copy(m.TBS.I2PProfilePath(), filepath.Join(m.TBD.DownloadPath, "i2p.firefox"))
	return m.server.ListenAndServe() //http.ListenAndServe(m.GetAddress(), nosurf.New(m))
}

func (m *Client) Shutdown(ctx context.Context) error {
	m.TBS.StopTor()
	return m.server.Shutdown(ctx)
}
