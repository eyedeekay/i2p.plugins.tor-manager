Tor(And sometimes Firefox) Manager for I2P
===========================================

## Usage: i2p.plugins.tor-manager [options]

### Options:

```sh

Usage: i2p.plugins.tor-manager [options]

Downloads, verifies and unpacks Tor Browser. Manages the Tor Browser
system in environments where Tor is not in use. Monitors a long-running
Tor process and downloads updates when Tor is not available.

Options:

      --apparmor                         Generate apparmor rules
      --arch string                      OS/arch to download (default "64")
      --bemirror                         Act as an in-I2P mirror when you're done downloading
      --chat                             Open a WebChat client
      --clearnet                         Use clearnet (no Tor or I2P) in Tor Browser
      --destruct                         Destructively delete the working directory when finished
      --directory string                 Directory operate in (default "tmp-i2pbrowser")
      --help                             Print help and quit
      --host string                      Host to serve on (default "127.0.0.1")
      --i2pbrowser                       Open I2P in Tor Browser
      --i2pconfig                        Open I2P routerconsole in Tor Browser with javscript enabled and non-routerconsole sites disabled
      --i2peditor                        Open I2P Site Editor in Tor Browser
      --lang string                      Language to download (default "en-US")
      --license                          Print the license and exit
      --mirror string                    Mirror to use. I2P will be used if an I2P proxy is present, if system Tor is available, it will be downloaded over the Tor proxy. (default "https://dist.torproject.org/torbrowser/")
      --mirrorall                        Download and mirror every language and OS/arch combination
      --nevertor                         Never use Tor for downloading Tor Browser
      --notor                            Do not automatically start Tor
      --nounpack                         Do not unpack the Tor Browser
      --offline                          Work offline. Differs from Firefox's offline mode in that cannot be disabled until the browser is closed.
      --onion                            Serve an onion site which shows some I2P propaganda (default true)
      --os string                        OS/arch to download (default "linux")
      --p2p                              Use bittorrent over I2P to download the initial copy of Tor Browser (default true)
      --password string                  Password to encrypt the working directory with. Implies -destruct, only the encrypted container will be saved.
      --port int                         Port to serve on (default 7695)
      --profile string                   use a custom profile path, normally blank
      --shortcuts                        Create desktop shortcuts
      --snowflake                        Offer a Snowflake to other Tor Browser users
      --snowflake-broker string          broker URL (default "https://snowflake-broker.torproject.net/")
      --snowflake-capacity uint          maximum concurrent clients
      --snowflake-directory string       directory with a page to serve locally for your snowflake. If empty, no local page is served.
      --snowflake-keep-local-addresses   keep local LAN address ICE candidates
      --snowflake-log string             log filename
      --snowflake-port string            port to serve info page(directory) on (default "7676")
      --snowflake-relay string           websocket relay URL (default "wss://snowflake.bamsoftware.com/")
      --snowflake-stun string            broker URL (default "stun:stun.stunprotocol.org:3478")
      --snowflake-unsafe-logging         prevent logs from being scrubbed
      --systray                          Create a systray icon
      --torbrowser                       Open Tor Browser
      --torrent                          Create a torrent of the downloaded files and seed it over I2P using an Open Tracker (default true)
      --torversion                       Print the version of Tor Browser that will be downloaded and exit
      --verbose                          Verbose output
      --watch-profiles string            Monitor and control these Firefox profiles. Temporarily Unused.

Available Languages:

  - sv-SE
  - el
  - fr
  - he
  - ro
  - zh-CN
  - ca
  - es-ES
  - is
  - ms
  - mk
  - nb-NO
  - tr
  - de
  - hu
  - it
  - ja
  - en-US
  - es-AR
  - lt
  - th
  - zh-TW
  - ar
  - da
  - pt-BR
  - ru
  - ga-IE
  - id
  - vi
  - fa
  - my
  - nl
  - cs
  - ka
  - ko
  - pl

