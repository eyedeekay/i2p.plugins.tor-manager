package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	i2phttpproxy "github.com/eyedeekay/httptunnel"
	i2pbrowserproxy "github.com/eyedeekay/httptunnel/multiproxy"
)

var (
	watchProfiles        = flag.String("watch-profiles", "", "Monitor and control these Firefox profiles. Temporarily Unused.")
	aggressiveIsolation  = false
	samHostString        = "127.0.0.1" //flag.String("bridge-host", "127.0.0.1", "host: of the SAM bridge")
	samPortString        = "7656"      //flag.String("bridge-port", "7656", ":port of the SAM bridge")
	destfile             = "invalid.tunkey"
	debugConnection      = false
	inboundTunnelLength  = 3
	outboundTunnelLength = 3
	inboundTunnels       = 4
	outboundTunnels      = 4
	inboundBackups       = 2
	outboundBackups      = 2
	inboundVariance      = 0
	outboundVariance     = 0
	dontPublishLease     = true
	reduceIdle           = false
	useCompression       = true
	reduceIdleTime       = 2000000
	reduceIdleQuantity   = 1
	runQuiet             = false
)

var addr string

func proxy() {
	ln, err := net.Listen("tcp", "127.0.0.1:4444")
	if err != nil {
		log.Fatal(err)
	}
	cln, err := net.Listen("tcp", "127.0.0.1:7696")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go proxyMain(ctx, ln, cln)
	<-ctx.Done()
	cancel()
}

func proxyMain(ctx context.Context, ln net.Listener, cln net.Listener) {
	flag.Parse()
	for {
		_, err := net.Listen("tcp", samHostString+":"+samPortString)
		if err != nil {
			break
		}
	}
	profiles := strings.Split(*watchProfiles, ",")

	srv := &http.Server{
		ReadTimeout:  600 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ln.Addr().String(),
	}
	var err error
	srv.Handler, err = i2pbrowserproxy.NewHttpProxy(
		i2pbrowserproxy.SetHost(samHostString),
		i2pbrowserproxy.SetPort(samPortString),
		i2pbrowserproxy.SetProxyAddr(ln.Addr().String()),
		i2pbrowserproxy.SetControlAddr(cln.Addr().String()),
		i2pbrowserproxy.SetDebug(debugConnection),
		i2pbrowserproxy.SetInLength(uint(inboundTunnelLength)),
		i2pbrowserproxy.SetOutLength(uint(outboundTunnelLength)),
		i2pbrowserproxy.SetInQuantity(uint(inboundTunnels)),
		i2pbrowserproxy.SetOutQuantity(uint(outboundTunnels)),
		i2pbrowserproxy.SetInBackups(uint(inboundBackups)),
		i2pbrowserproxy.SetOutBackups(uint(outboundBackups)),
		i2pbrowserproxy.SetInVariance(inboundVariance),
		i2pbrowserproxy.SetOutVariance(outboundVariance),
		i2pbrowserproxy.SetUnpublished(dontPublishLease),
		i2pbrowserproxy.SetReduceIdle(reduceIdle),
		i2pbrowserproxy.SetCompression(useCompression),
		i2pbrowserproxy.SetReduceIdleTime(uint(reduceIdleTime)),
		i2pbrowserproxy.SetReduceIdleQuantity(uint(reduceIdleQuantity)),
		i2pbrowserproxy.SetKeysPath(destfile),
		i2pbrowserproxy.SetProxyMode(aggressiveIsolation),
	)
	i2pbrowserproxy.Quiet = runQuiet
	if err != nil {
		log.Fatal(err)
	}

	ctrlsrv := &http.Server{
		ReadHeaderTimeout: 600 * time.Second,
		WriteTimeout:      600 * time.Second,
		Addr:              cln.Addr().String(),
	}
	ctrlsrv.Handler, err = i2phttpproxy.NewSAMHTTPController(profiles, srv)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				srv.Handler.(*i2pbrowserproxy.SAMMultiProxy).Close()
				srv.Shutdown(ctx)
				ctrlsrv.Shutdown(ctx)
			}
		}
	}()

	go func() {
		log.Println("Starting control server on", cln.Addr())
		if err := ctrlsrv.Serve(cln); err != nil {
			if err == http.ErrServerClosed {
				return
			}
			log.Fatal("Serve:", err)
		}
		log.Println("Stopping control server on", cln.Addr())
	}()

	go func() {
		log.Println("Starting proxy server on", ln.Addr())
		if err := srv.Serve(ln); err != nil {
			if err == http.ErrServerClosed {
				return
			}
			log.Fatal("Serve:", err)
		}
		log.Println("Stopping proxy server on", ln.Addr())
	}()

	counter()

	<-ctx.Done()
}

func counter() {
	var x int
	for {
		log.Println("Identity is", x, "minute(s) old")
		time.Sleep(1 * time.Minute)
		x++
	}
}
