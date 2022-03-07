package tbserve

import (
	"io/ioutil"
	"path/filepath"

	"github.com/russross/blackfriday"
)

var dmd string = `
# Tor Binary Manager

This plugin manages the Tor Browser Bundle and a Tor binary
for you. Combined with a SOCKS5 plugin for I2P, it acts as
an alternative to a fixed outproxy by using Tor, and also
provides a way to run I2P in the Tor Browser without any other
configuration.

 - [![Launch I2P in Tor Browser](garliconion.png) - Launch I2P in Tor Browser](/launch-i2p-browser)
 - [![Launch Tor Browser](onion.png) - Launch Tor Browser](/launch-tor-browser)
 
## Tor Controls

`
var defaultmd []byte = []byte(dmd)

var hhd string = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Tor Binary Manager</title>
<link rel="stylesheet" href="/style.css">
</head>
`

var htmlhead []byte = []byte(hhd)

var tstart string = `
- [![Stop Tor](/stop-tor.png)](/stop-tor) 
`

var torstart []byte = []byte(tstart)

var tstop string = `
- [![Start Tor](/start-tor.png)](/start-tor)
`

var torstop []byte = []byte(tstop)

var trun string = `
- Tor is Running as a System Service
`

var torrunning []byte = []byte(trun)

var tstopped string = `
- Tor is Stopped and there is no System Service
`

var torstopped []byte = []byte(tstopped)

// PageHTML returns the HTML for the page heading
func (m *Client) PageHTML() []byte {
	dir := filepath.Dir(m.TBD.DownloadPath)
	mdpath := filepath.Join(dir, m.TBD.Lang, "index.md")
	mdbytes, err := ioutil.ReadFile(mdpath)
	markdown := blackfriday.New()
	if err != nil {
		htmlbytes := markdown.Parse(defaultmd)
		return []byte(htmlbytes.String())
	}
	htmlbytes := markdown.Parse(mdbytes)
	return []byte(htmlbytes.String())
}

// TorOnStatusHTML returns the HTML for "Tor Status" section the page
func (m *Client) TorOnStatusHTML(ours bool) []byte {
	dir := filepath.Dir(m.TBD.DownloadPath)
	markdown := blackfriday.New()
	if ours {
		mdpath := filepath.Join(dir, m.TBD.Lang, "stoptor.md")
		torbytes, err := ioutil.ReadFile(mdpath)
		if err != nil {
			htmlbytes := markdown.Parse(torstop)
			return []byte(htmlbytes.String())
		}
		htmlbytes := markdown.Parse(torbytes)
		return []byte(htmlbytes.String())
	}
	mdpath := filepath.Join(dir, m.TBD.Lang, "torstarted.md")
	toron, err := ioutil.ReadFile(mdpath)
	if err != nil {
		htmlbytes := markdown.Parse(torrunning)
		return []byte(htmlbytes.String())
	}
	htmlbytes := markdown.Parse(toron)
	return []byte(htmlbytes.String())
}

// TorOffStatusHTML returns the HTML for "Tor Status" section the page
func (m *Client) TorOffStatusHTML(ours bool) []byte {
	dir := filepath.Dir(m.TBD.DownloadPath)
	markdown := blackfriday.New()
	if ours {
		mdpath := filepath.Join(dir, m.TBD.Lang, "startor.md")
		torbytes, err := ioutil.ReadFile(mdpath)
		if err != nil {
			htmlbytes := markdown.Parse(torstart)
			return []byte(htmlbytes.String())
		}
		htmlbytes := markdown.Parse(torbytes)
		return []byte(htmlbytes.String())
	}
	mdpath := filepath.Join(dir, m.TBD.Lang, "torstopped.md")
	toroff, err := ioutil.ReadFile(mdpath)
	if err != nil {
		htmlbytes := markdown.Parse(torstopped)
		return []byte(htmlbytes.String())
	}
	htmlbytes := markdown.Parse(toroff)
	return []byte(htmlbytes.String())
}
