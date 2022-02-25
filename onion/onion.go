package i2pdotonion

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/cretz/bine/tor"
)

type I2POnionService struct {
	OnionService net.Listener
	ServeDir     string
}

func NewOnionService(dir string) *I2POnionService {
	return &I2POnionService{ServeDir: dir}
}

func (ios *I2POnionService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := path.Clean(r.URL.Path)
	if path == "/" {
		path = "/index.html"
	}
	path = filepath.Join(ios.ServeDir, path)

	http.ServeFile(w, r, path)
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
