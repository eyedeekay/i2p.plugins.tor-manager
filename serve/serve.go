package tbserve

import (
	"embed"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/justinas/nosurf"
	cp "github.com/otiai10/copy"
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
	TBSupervise "i2pgit.org/idk/i2p.plugins.tor-manager/supervise"
)

// Client manages and supervises a Tor Browser instance.
type Client struct {
	hostname string
	TBD      *tbget.TBDownloader
	TBS      *TBSupervise.Supervisor
	DarkMode bool
	Host     string
	Port     int
}

// NewClient creates a new Client.
func NewClient(verbose bool, lang string, os string, arch string, content *embed.FS) (*Client, error) {
	m := &Client{
		TBD: tbget.NewTBDownloader(lang, os, arch, content),
	}
	m.TBD.Verbose = verbose
	m.TBD.MakeTBDirectory()
	tgz, sig, err := m.TBD.DownloadUpdaterForLang(lang)
	if err != nil {
		panic(err)
	}
	var home string
	if home, err = m.TBD.CheckSignature(tgz, sig); err != nil {
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
	path := strings.Replace(rq.URL.Path, "..", "", -1)
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
	ioutil.WriteFile(filepath.Join(m.TBD.DownloadPath, "mirror.json"), []byte(mirrorjson), 0644)
	cp.Copy(m.TBS.I2PProfilePath(), filepath.Join(m.TBD.DownloadPath, "i2p.firefox"))
	go m.TBS.RunTorWithLang()
	return http.ListenAndServe(m.GetAddress(), nosurf.New(m))
}
