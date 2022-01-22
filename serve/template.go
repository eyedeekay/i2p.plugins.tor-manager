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
 - [Start Tor](/start-tor)
 - [Stop Tor](/stop-tor)

`)

func (m *Client) Page() (string, error) {
	dir := filepath.Dir(m.TBD.DownloadPath)
	mdpath := filepath.Join(dir, m.TBD.Lang, "index.md")
	mdbytes, err := ioutil.ReadFile(mdpath)
	if err != nil {
		return string(blackfriday.MarkdownCommon(defaultmd)), err
	}
	htmlbytes := blackfriday.MarkdownCommon(mdbytes)
	return string(htmlbytes), nil
}
