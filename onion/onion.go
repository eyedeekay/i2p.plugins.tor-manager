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

	"github.com/cretz/bine/tor"
)

//go:embed www/*
var content embed.FS

type I2POnionService struct {
	OnionService net.Listener
	ServeDir     string
}

func NewOnionService(dir string) (*I2POnionService, error) {
	ios := &I2POnionService{
		ServeDir: filepath.Join(dir, "www"),
	}
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
	fmt.Println("ServeHTTP:", path)
	path = filepath.Join(ios.ServeDir, path)
	fmt.Println("ServeHTTP:", path)
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
	//defer tb.Close()
	//Wait at most a few minutes to publish the service
	//listenCtx, listenCancel
	listenCtx := context.Background() //context.WithTimeout(context.Background(), 3*time.Minute)
	//defer listenCancel()
	//Create a v3 onion service to listen on any port but show as 80
	ios.OnionService, err = tb.Listen(listenCtx, &tor.ListenConf{Version3: true, RemotePorts: []int{80}})
	if err != nil {
		log.Panicf("Unable to create onion service: %v", err)
	}
	return ios.OnionService, nil
}

func (ios *I2POnionService) Serve(l net.Listener) error {
	ios.OnionService = l
	log.Printf("Serve: %s", ios.OnionService.Addr())
	return http.Serve(ios.OnionService, ios)
}

func (ios *I2POnionService) ListenAndServe() error {
	var err error
	ios.OnionService, err = ios.Listen("", "")
	if err != nil {
		return err
	}
	log.Printf("ListenAndServe: %s", ios.OnionService.Addr())
	return http.Serve(ios.OnionService, ios)
}

func (ios *I2POnionService) UnpackSite() error {
	docroot := ios.ServeDir
	fmt.Println("UnpackSite: ", docroot)
	if dir, err := os.Stat(docroot); err == nil && dir.IsDir() {
		return nil
	}
	os.MkdirAll(docroot, 0755)
	//unpack the contents to the docroot
	return fs.WalkDir(content, ".", func(embedpath string, d fs.DirEntry, err error) error {

		fp := filepath.Join(docroot)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("UnpackSite: ", embedpath, fp, embedpath)
		if d.IsDir() {
			log.Println("UnpackSite: mkdir", filepath.Join(fp, strings.Replace(embedpath, "www", "", -1)))
			os.MkdirAll(filepath.Join(fp, strings.Replace(embedpath, "www", "", -1)), 0755)
		} else {
			log.Println("UnpackSite: copy", filepath.Join(fp, strings.Replace(embedpath, "www", "", -1)))
			fullpath := path.Join(embedpath)
			bytes, err := content.ReadFile(fullpath)
			if err != nil {
				return err
			}
			unpack := filepath.Join(fp, strings.Replace(embedpath, "www", "", -1))
			if err := ioutil.WriteFile(unpack, bytes, 0644); err != nil {
				return err
			}
		}
		return nil
	})
}
