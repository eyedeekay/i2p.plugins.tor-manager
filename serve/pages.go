package tbserve

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func (m *Client) Page() (string, error) {

	htmlbytes := htmlhead
	htmlbytes = append(htmlbytes, []byte("<body>")...)

	mdbytes := m.PageHTML()
	htmlbytes = append(htmlbytes, mdbytes...)

	if alive, ours := m.TBS.TorIsAlive(); alive {
		htmlbytes = append(htmlbytes, m.TorOnStatusHTML(ours)...)
	} else {
		htmlbytes = append(htmlbytes, m.TorOffStatusHTML(ours)...)
	}
	htmlbytes = append(htmlbytes, []byte(`</body>
	<html>`)...)
	return string(htmlbytes), nil
}

func (m *Client) serveJSON(rw http.ResponseWriter, rq *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	mirrorjson, err := m.GenerateMirrorJSON()
	if err != nil {
		rw.Write([]byte(`{
			"version": "latest",
			"torbrowser": {
				"version": "latest",
				"download": "https://www.torproject.org/download/download-easy.html.en"
			},
			"i2p": {
				"version": "latest",
				"download": "https://geti2p.net/en/download"
			}
		}`))
		return
	}
	rw.Write([]byte(mirrorjson))
}

func (m *Client) serveCSS(rw http.ResponseWriter, rq *http.Request) {
	cssbytes, err := ioutil.ReadFile(filepath.Join(m.TBD.DownloadPath, rq.URL.Path))
	if err != nil {
		rw.Header().Set("Content-Type", "text/css")
		rw.Write(defaultCSS)
	}
	rw.Header().Set("Content-Type", "text/css")
	rw.Write(cssbytes)
}

func (m *Client) serveJS(rw http.ResponseWriter, rq *http.Request) {
	jsbytes, err := ioutil.ReadFile(filepath.Join(m.TBD.DownloadPath, rq.URL.Path))
	if err != nil {
		return
	}
	rw.Header().Set("Content-Type", "application/javascript")
	rw.Write(jsbytes)
}

func (m *Client) servePNG(rw http.ResponseWriter, rq *http.Request) {
	pngbytes, err := ioutil.ReadFile(filepath.Join(m.TBD.DownloadPath, rq.URL.Path))
	if err != nil {
		return
	}
	rw.Header().Set("Content-Type", "image/png")
	rw.Write(pngbytes)
}

func (m *Client) serveICO(rw http.ResponseWriter, rq *http.Request) {
	pngbytes, err := ioutil.ReadFile(filepath.Join(m.TBD.DownloadPath, rq.URL.Path))
	if err != nil {
		return
	}
	rw.Header().Set("Content-Type", "image/x-icon")
	rw.Write(pngbytes)
}

func (m *Client) serveSVG(rw http.ResponseWriter, rq *http.Request) {
	pngbytes, err := ioutil.ReadFile(filepath.Join(m.TBD.DownloadPath, rq.URL.Path))
	if err != nil {
		return
	}
	rw.Header().Set("Content-Type", "image/svg+xml")
	rw.Write(pngbytes)
}
