package tbserve

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// Page generates the HTML for the panel.
func (m *Client) Page() (string, error) {

	htmlbytes := htmlhead

	htmlbytes = append(htmlbytes, []byte(`<body>
	<label class="switch">
	  <input type="checkbox" onclick='handleClick(this);'>
	  <span class="slider round"></span>
	</label>`)...)
	htmlbytes = append(htmlbytes, []byte(`<script>
	function handleClick(cb) {
		var xmlHttp = new XMLHttpRequest();
		xmlHttp.open( "GET", "http://`)...)

	htmlbytes = append(htmlbytes, []byte([]byte(m.GetAddress()))...)

	htmlbytes = append(htmlbytes, []byte(`/switch-theme", false ); // false for synchronous request
		xmlHttp.send( null );
		location.reload();
		return xmlHttp.responseText;
	}
	</script>
	`)...)

	mdbytes := m.PageHTML()
	htmlbytes = append(htmlbytes, mdbytes...)

	if alive, ours := m.TBS.TorIsAlive(); alive {
		htmlbytes = append(htmlbytes, m.TorOnStatusHTML(ours)...)
	} else {
		htmlbytes = append(htmlbytes, m.TorOffStatusHTML(ours)...)
	}
	htmlbytes = append(htmlbytes, []byte(`</body>
	</html>`)...)
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
		if m.DarkMode {
			rw.Header().Set("Content-Type", "text/css")
			rw.Write(darkDefaultCSS)
			return
		}
		rw.Header().Set("Content-Type", "text/css")
		rw.Write(defaultCSS)
		return
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
		rw.Header().Set("Content-Type", "image/png")
		if strings.HasSuffix(rq.URL.Path, "garliconion.png") {
			if bytes, err := m.TBD.Profile.ReadFile("garliconion.png"); err == nil {
				rw.Write(bytes)
				return
			}
		}
		if strings.HasSuffix(rq.URL.Path, "onion.png") {
			if bytes, err := m.TBD.Profile.ReadFile("onion.png"); err == nil {
				rw.Write(bytes)
				return
			}
		}
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
