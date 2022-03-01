package i2pdotonion

import (
	"context"
	"crypto"
	"crypto/ed25519"
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
	Keys         crypto.PrivateKey
}

func NewOnionService(dir string) (*I2POnionService, error) {
	ios := &I2POnionService{
		ServeDir: filepath.Join(dir, "www"),
	}
	if err := ios.UnpackSite(); err != nil {
		return nil, err
	}

	if file, err := os.Stat(ios.KeysPath()); err == nil && file.Mode().IsRegular() {
		ios.Keys, err = torKeys(ios.KeysPath())
		if err != nil {
			return nil, err
		}
	}
	return ios, nil
}

func torKeys(addr string) (crypto.PrivateKey, error) {
	//log.Infof("Starting and registering onion service, please wait a couple of minutes...")
	//t, err := tor.Start(nil, nil)
	//if err != nil {
	//	log.Fatalf("Unable to start Tor: %v", err)
	//}
	var keys *ed25519.PrivateKey
	if _, err := os.Stat(addr + ".tor.private"); os.IsNotExist(err) {
		_, tkeys, err := ed25519.GenerateKey(nil)
		if err != nil {
			log.Fatalf("Unable to generate onion service key, %s", err)
		}
		keys = &tkeys
		f, err := os.Create(addr + ".tor.private")
		if err != nil {
			log.Fatalf("Unable to create Tor keys file for writing, %s", err)
		}
		defer f.Close()
		_, err = f.Write(tkeys)
		if err != nil {
			log.Fatalf("Unable to write Tor keys to disk, %s", err)
		}
	} else if err == nil {
		tkeys, err := ioutil.ReadFile(addr + ".tor.private")
		if err != nil {
			log.Fatalf("Unable to read Tor keys from disk")
		}
		k := ed25519.NewKeyFromSeed(tkeys)
		keys = &k
	} else {
		log.Fatalf("Unable to set up Tor keys, %s", err)
	}
	return keys, nil
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
	t, err := tor.Start(nil, &tor.StartConf{
		DataDir: filepath.Join(filepath.Dir(ios.ServeDir), "tor"),
	})
	t.DeleteDataDirOnClose = true
	if err != nil {
		return nil, fmt.Errorf("Unable to start Tor: %v", err)
	}
	//var err error
	listenCtx := context.Background()
	// Create a v3 onion service to listen on any port but show as 6667
	ios.OnionService, err = t.Listen(
		listenCtx,
		&tor.ListenConf{
			Version3:    true,
			RemotePorts: []int{80},
			Key:         ios.Keys,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to listen on Tor: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("Unable to write Tor public key to disk, %s", err)
	}
	return ios.OnionService, nil
}

func (ios *I2POnionService) KeysPath() string {
	return filepath.Join(filepath.Dir(filepath.Dir(ios.ServeDir)), "service.tor.private")
}

func GenerateTorKeys(file string) (*ed25519.PrivateKey, error) {
	var keys *ed25519.PrivateKey
	if _, err := os.Stat(file); os.IsNotExist(err) {
		//tkeys
		_, tkeys, err := ed25519.GenerateKey(nil)
		if err != nil {
			log.Fatalf("Unable to generate onion service key, %s", err)
		}
		keys = &tkeys
		f, err := os.Create(file)
		if err != nil {
			log.Fatalf("Unable to create Tor keys file for writing, %s", err)
		}
		defer f.Close()
		_, err = f.Write(tkeys)
		if err != nil {
			log.Fatalf("Unable to write Tor keys to disk, %s", err)
		}
	}
	return keys, nil
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
		if d.IsDir() {
			os.MkdirAll(filepath.Join(fp, strings.Replace(embedpath, "www", "", -1)), 0755)
		} else {
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
