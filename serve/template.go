package tbserve

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/russross/blackfriday"
)

var defaultmd []byte = []byte(`
# Tor Binary Manager

This plugin manages the Tor Browser Bundle and a Tor binary
for you. Combined with a SOCKS5 plugin for I2P, it acts as
an alternative to a fixed outproxy by using Tor, and also
provides a way to run I2P in the Tor Browser without any other
configuration.

 - [Launch I2P in Tor Browser](/launch-i2p-browser)
 - [Launch Tor Browser](/launch-tor-browser)
 
## Tor Controls

`)

var htmlhead []byte = []byte(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Tor Binary Manager</title>
<link rel="stylesheet" href="/style.css">
</head>
`)

var torstart []byte = []byte(`
- [![Stop Tor](/stop-tor.png)](/stop-tor) 
`)

var torstop []byte = []byte(`
- [![Start Tor](/start-tor.png)](/start-tor)
`)

var torrunning []byte = []byte(`
- Tor is Running as a System Service
`)

var torstopped []byte = []byte(`
- Tor is Stopped and there is no System Service
`)

func (m *Client) PageHTML() []byte {
	dir := filepath.Dir(m.TBD.DownloadPath)
	mdpath := filepath.Join(dir, m.TBD.Lang, "index.md")
	mdbytes, err := ioutil.ReadFile(mdpath)
	if err != nil {
		htmlbytes := blackfriday.MarkdownCommon(defaultmd)
		return htmlbytes
	}
	htmlbytes := blackfriday.MarkdownCommon(mdbytes)
	return htmlbytes
}

func (m *Client) TorOnStatusHTML(ours bool) []byte {
	dir := filepath.Dir(m.TBD.DownloadPath)
	if ours {
		mdpath := filepath.Join(dir, m.TBD.Lang, "stoptor.md")
		torbytes, err := ioutil.ReadFile(mdpath)
		if err != nil {
			htmlbytes := blackfriday.MarkdownCommon(torstop)
			return htmlbytes
		} else {
			htmlbytes := blackfriday.MarkdownCommon(torbytes)
			return htmlbytes
		}
	} else {
		mdpath := filepath.Join(dir, m.TBD.Lang, "toron.md")
		toron, err := ioutil.ReadFile(mdpath)
		if err != nil {
			htmlbytes := blackfriday.MarkdownCommon(torrunning)
			return htmlbytes
		} else {
			htmlbytes := blackfriday.MarkdownCommon(toron)
			return htmlbytes
		}
	}
}

func (m *Client) TorOffStatusHTML(ours bool) []byte {
	dir := filepath.Dir(m.TBD.DownloadPath)
	if ours {
		mdpath := filepath.Join(dir, m.TBD.Lang, "stoptor.md")
		torbytes, err := ioutil.ReadFile(mdpath)
		if err != nil {
			htmlbytes := blackfriday.MarkdownCommon(torstart)
			return htmlbytes
		} else {
			htmlbytes := blackfriday.MarkdownCommon(torbytes)
			return htmlbytes
		}
	} else {
		mdpath := filepath.Join(dir, m.TBD.Lang, "toron.md")
		toroff, err := ioutil.ReadFile(mdpath)
		if err != nil {
			htmlbytes := blackfriday.MarkdownCommon(torstopped)
			return htmlbytes
		} else {
			htmlbytes := blackfriday.MarkdownCommon(toroff)
			return htmlbytes
		}
	}
}

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
