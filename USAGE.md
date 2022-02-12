Tor(And sometimes Firefox) Manager for I2P
===========================================

## Usage: i2p.plugins.tor-manager [options]

### Options:

```sh
Usage: i2p.plugins.tor-manager [options]

Downloads, verifies and unpacks Tor Browser. Manages the Tor Browser
system in environments where Tor is not in use.

Options:

Usage of ./i2p.plugins.tor-manager:
  -apparmor
    	Generate apparmor rules
  -arch string
    	OS/arch to download (default "64")
  -bemirror
    	Act as an in-I2P mirror when you're done downloading
  -clearnet
    	Use clearnet (no Tor or I2P)
  -directory string
    	Directory operate in
  -help
    	Print help
  -host string
    	Host to serve on (default "127.0.0.1")
  -i2pbrowser
    	Open I2P in Tor Browser
  -i2pconfig
    	Open I2P routerconsole in Tor Browser with javscript enabled and non-routerconsole sites disabled
  -lang string
    	Language to download
  -offline
    	Work offline. Differs from Firefox's offline mode in that cannot be disabled until the browser is closed.
  -os string
    	OS/arch to download (default "linux")
  -port int
    	Port to serve on (default 7695)
  -profile string
    	use a custom profile path, normally blank
  -shortcuts
    	Create desktop shortcuts
  -torbrowser
    	Open Tor Browser
  -verbose
    	Verbose output
  -watch-profiles string
    	Monitor and control these Firefox profiles. Temporarily Unused.
```

