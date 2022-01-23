package tbserve

import (
	"io/ioutil"
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

func (m *Client) Page() (string, error) {
	dir := filepath.Dir(m.TBD.DownloadPath)
	mdpath := filepath.Join(dir, m.TBD.Lang, "index.md")
	mdbytes, err := ioutil.ReadFile(mdpath)
	if err != nil {
		return string(blackfriday.MarkdownCommon(defaultmd)), err
	}
	htmlbytes := htmlhead
	htmlbytes = append(htmlbytes, []byte("<body>")...)
	htmlbytes = append(htmlbytes, blackfriday.MarkdownCommon(mdbytes)...)
	if m.TBS.TorIsAlive() {
		htmlbytes = append(htmlbytes, []byte(`
- [Stop Tor](/stop-tor)
`)...)
	} else {
		htmlbytes = append(htmlbytes, []byte(`
- [Start Tor](/start-tor)
`)...)
	}
	htmlbytes = append(htmlbytes, []byte(`</body>
	<html>`)...)
	return string(htmlbytes), nil
}
