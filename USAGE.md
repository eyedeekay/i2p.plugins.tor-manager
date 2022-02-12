Tor(And sometimes Firefox) Manager for I2P
===========================================

## Usage: i2p.plugins.tor-manager [options]

### Options:

```sh
Unpacking /home/idk/go/src/i2pgit.org/idk/i2p.plugins.tor-manager/tor-browser/torbrowser-linux64-en-US.tar.xz /home/idk/go/src/i2pgit.org/idk/i2p.plugins.tor-manager/unpack
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
Usage: ./firefox.real [ options ... ] [URL]
       where options include:

X11 options
  --display=DISPLAY  X display to use
  --sync             Make X calls synchronous
  --g-fatal-warnings Make all warnings fatal

Firefox options
  -h or --help       Print this message.
  -v or --version    Print Firefox version.
  --full-version     Print Firefox version, build and platform build ids.
  -P <profile>       Start with <profile>.
  --profile <path>   Start with profile at <path>.
  --migration        Start with migration wizard.
  --ProfileManager   Start with ProfileManager.
  --no-remote        (default) Do not accept or send remote commands; implies
                     --new-instance.
  --allow-remote     Accept and send remote commands.
  --new-instance     Open new instance, not a new window in running instance.
  --safe-mode        Disables extensions and themes for this session.
  --MOZ_LOG=<modules> Treated as MOZ_LOG=<modules> environment variable,
                     overrides it.
  --MOZ_LOG_FILE=<file> Treated as MOZ_LOG_FILE=<file> environment variable,
                     overrides it. If MOZ_LOG_FILE is not specified as an
                     argument or as an environment variable, logging will be
                     written to stdout.
  --headless         Run without a GUI.
  --browser          Open a browser window.
  --new-window <url> Open <url> in a new window.
  --new-tab <url>    Open <url> in a new tab.
  --private-window <url> Open <url> in a new private window.
  --preferences      Open Preferences dialog.
  --screenshot [<path>] Save screenshot to <path> or in working directory.
  --window-size width[,height] Width and optionally height of screenshot.
  --search <term>    Search <term> with your default search engine.
  --setDefaultBrowser Set this app as the default browser.
  --first-startup    Run post-install actions before opening a new window.
  --kiosk Start the browser in kiosk mode.
  --jsconsole        Open the Browser Console.
  --jsdebugger [<path>] Open the Browser Toolbox. Defaults to the local build
                     but can be overridden by a firefox path.
  --wait-for-jsdebugger Spin event loop until JS debugger connects.
                     Enables debugging (some) application startup code paths.
                     Only has an effect when `--jsdebugger` is also supplied.
  --devtools         Open DevTools on initial load.
  --start-debugger-server [ws:][ <port> | <path> ] Start the devtools server on
                     a TCP port or Unix domain socket path. Defaults to TCP port
                     6000. Use WebSocket protocol if ws: prefix is specified.
  --recording <file> Record drawing for a given URL.
  --recording-output <file> Specify destination file for a drawing recording.
  --remote-debugging-port [<port>] Start the Firefox remote agent,
                     which is a low-level debugging interface based on the CDP protocol.
                     Defaults to listen on localhost:9222.

Tor Browser Script Options
  --verbose         Display Tor and Firefox output in the terminal
  --log [file]      Record Tor and Firefox output in file (default: tor-browser.log)
  --detach          Detach from terminal and run Tor Browser in the background.
  --register-app    Register Tor Browser as a desktop app for this user
  --unregister-app  Unregister Tor Browser as a desktop app for this user
```

