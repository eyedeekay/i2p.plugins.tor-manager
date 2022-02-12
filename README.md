# i2p.plugins.tor-updater

A Tor package updater and runner as an I2P Plugin. This plugin is
usable on Windows, Linux, and OSX, as is the freestanding binary.
This also functions as a freestanding update for the Tor Browser
Bundle and is capable of configuring Tor Browser from the terminal
and updating it without running it, should the user choose to operate
this way.

Usage:
------

See [Usage](USAGE.md) for command-line usage.

[HTML version](usage.html)

Plugin:
-------

- [i2p.plugins.tor-manager-linux-386](i2p.plugins.tor-manager-linux-386.su3)
- [i2p.plugins.tor-manager-windows-amd64](i2p.plugins.tor-manager-windows-amd64.su3)
- [i2p.plugins.tor-manager-darwin-arm64](i2p.plugins.tor-manager-darwin-arm64.su3)
- [i2p.plugins.tor-manager-linux-amd64](i2p.plugins.tor-manager-linux-amd64.su3)
- [i2p.plugins.tor-manager-windows-386](i2p.plugins.tor-manager-windows-386.su3)
- [i2p.plugins.tor-manager-darwin-amd64](i2p.plugins.tor-manager-darwin-amd64.su3)

Status:
-------

![Screenshot 2](screenshot-console.png)

Linux: Usable, everything implemented works.
Windows: Usable, everything implemented works.
OSX: Usable, everything implemented works.

Other systems are not targeted and should use a Tor binary built from source,
provided by TPO or, their prefered package management system and not this plugin.
The plugin will not start a Tor instance if a SOCKS proxy is open on port 9050.

![Screenshot](screenshot-i2pbrowser.png)

### Primary Goals

1. Ship known-good public keys, download a current Tor for the platform in the background, authenticate it, and launch it only if necessary.
 - Works on Windows, Linux, OSX
2. Supervise Tor as a ShellService plugin to I2P
 - Works on Linux, Windows, OSX
3. Keep Tor up-to-date
 - Works on Windows, Linux, OSX
4. Work as an I2P Plugin OR as a freestanding app to be compatible with all I2P distributions
 - Works on Linux, Windows, OSX
5. Download Tor Browser from an in-I2P mirror(or one of a network of in-I2P mirrors)
 - Not done

### Secondary Goals:

1. Launch Tor Browser
 - Works on Linux, Windows, OSX
2. Configure and launch Tor browser for use with I2P
 - Works on Linux, Windows, OSX

#### Optional Features I might add if there is interest

1. Mirror the files which it downloads to an I2P Site
 - Works on Windows, Linux, OSX
2. Mirror the files which it downloads to I2P torrents
 - Not done
3. Set up an onion site which announces an I2P mirror exists
 - Not done
4. Use Bittorrent-over-I2P to download the Tor Browser software
 - Not Done
5. Import libi2pd and offer the use of an embedded i2pd router.
 - Not done.
6. Option to use BRB in a thread as an in-I2P replacement for `mibbit` IRC client.
 - Not done.

### Usage as a Library

[More information at the GoDoc](https://pkg.go.dev/i2pgit.org/idk/i2p.plugins.tor-manager)

This is also useful as a library for downloading a Tor Browser Bundle. This API
isn't really stable, more "stabilizing." Feel free to use it, but it may still
change a little.

To create a new instance, use:

``` Go
client, err = tbserve.NewClient(*verbose, *lang, *system, *arch, &content)
```

Customize the client using the exposed variables and methods:

``` Go
client.Host = *host
client.Port = *port
client.TBS.Profile = &content
client.TBS.PassThroughArgs = flag.Args()
```

And serve the controller:

``` Go
if err := client.Serve(); err != nil {
  log.Fatal(err)
}
```

### Similar Projects:

- https://github.com/micahflee/torbrowser-launcher
- https://github.com/whonix/tb-updater

<a href="https://www.flaticon.com/free-icons/garlic" title="garlic icons">Garlic icons created by Icongeek26 - Flaticon</a>
<a href="https://www.flaticon.com/free-icons/onion" title="onion icons">Onion icons created by Freepik - Flaticon</a>

### More Screenshots:

- ![Screenshot](screenshot.png)
- ![Screenshot](screenshot-dark.png)
