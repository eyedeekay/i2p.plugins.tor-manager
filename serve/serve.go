package tbserve

import (
	"embed"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
	TBSupervise "i2pgit.org/idk/i2p.plugins.tor-manager/supervise"
)

type Client struct {
	hostname string
	TBD      *tbget.TBDownloader
	TBS      *TBSupervise.Supervisor
}

func NewClient(hostname string, lang string, os string, arch string, content *embed.FS) (*Client, error) {
	m := &Client{
		hostname: hostname,
		TBD:      tbget.NewTBDownloader(lang, os, arch, content),
	}
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

func (m *Client) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	path := rq.URL.Path
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
		default:
			b, _ := m.Page()
			rw.Header().Set("Content-Type", "text/html")
			rw.Write([]byte(b))
		}
	}

}

func (m *Client) Serve() error {
	//http.Handle("/", m)
	go m.TBS.RunTorWithLang()
	return http.ListenAndServe("127.0.0.1:7695", nosurf.New(m))
}
