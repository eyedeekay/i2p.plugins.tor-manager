package i2pdotonion

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cretz/bine/tor"
)

//go:embded www/*
var content embed.FS

type I2POnionService struct {
	OnionService net.Listener
	ServeDir     string
}

func NewOnionService(dir string) (*I2POnionService, error) {
	ios := &I2POnionService{ServeDir: dir}
	if err := ios.UnpackSite(); err != nil {
		return nil, err
	}
	return ios, nil
}

func (ios *I2POnionService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := path.Clean(r.URL.Path)
	if path == "/" {
		path = "/index.html"
	}
	path = filepath.Join(ios.ServeDir, path)
	finfo, err := os.Stat(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if finfo.IsDir() {
		http.NotFound(w, r)
	}
	http.ServeFile(w, r, path)
}

func (ios *I2POnionService) StandardHTML() string {
	return ""
}

func (ios *I2POnionService) Listen(net, addr string) (net.Listener, error) {
	if ios.OnionService != nil {
		return ios.OnionService, nil
	}
	fmt.Println("Starting and registering onion service, please wait a couple of minutes...")
	tb, err := tor.Start(context.Background(), nil)
	if err != nil {
		log.Panicf("Unable to start Tor: %v", err)
	}
	defer tb.Close()
	// Wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()
	// Create a v3 onion service to listen on any port but show as 80
	ios.OnionService, err = tb.Listen(listenCtx, &tor.ListenConf{Version3: true, RemotePorts: []int{80}})
	if err != nil {
		log.Panicf("Unable to create onion service: %v", err)
	}
	return ios.OnionService, nil
}

func (ios *I2POnionService) Serve(l net.Listener) error {
	ios.OnionService = l
	return http.Serve(ios.OnionService, ios)
}

func (ios *I2POnionService) ListenAndServe() error {
	var err error
	ios.OnionService, err = ios.Listen("", "")
	if err != nil {
		return err
	}
	return http.Serve(ios.OnionService, ios)
}

func (ios *I2POnionService) UnpackSite() error {
	docroot := filepath.Join(ios.ServeDir, "www")
	os.MkdirAll(docroot, 0755)
	if dir, err := os.Stat(docroot); err == nil && dir.IsDir() {
		return nil
	}
	//unpack the contents to the docroot
	return fs.WalkDir(content, ".", func(embedpath string, d fs.DirEntry, err error) error {
		fp := filepath.Join(docroot)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(embedpath, filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox", "", -1)))
		if d.IsDir() {
			os.MkdirAll(filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox", "", -1)), 0755)
		} else {
			fullpath := path.Join(embedpath)
			bytes, err := content.ReadFile(fullpath)
			if err != nil {
				return err
			}
			unpack := filepath.Join(fp, strings.Replace(embedpath, "tor-browser/unpack/i2p.firefox", "", -1))
			if err := ioutil.WriteFile(unpack, bytes, 0644); err != nil {
				return err
			}
		}
		return nil
	})

}
