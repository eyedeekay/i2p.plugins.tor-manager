# i2p.plugins.tor-updater

A Tor package updater and runner as an I2P Plugin. This plugin is
still being changed rapidly but it should be usable on most Linux
distributions as of 23 Jan, 2022.

Usage:
------

```sh
Usage of ./i2p.plugins.tor-manager-linux-amd64:
  -arch string
    	OS/arch to download (default "64")
  -i2pbrowser
    	Open I2P in Tor Browser
  -lang string
    	Language to download
  -os string
    	OS/arch to download (default "linux")
  -torbrowser
    	Open Tor Browser
```

### Primary Goals


- Ship known-good public keys, download a current Tor for the platform in the background, authenticate it, and launch it only if necessary.
 - Works on Windows, Linux, probably also OSX
- Supervise Tor as a ShellService plugin to I2P
 - Works on Linux
- Keep Tor up-to-date
 - Works on Windows, Linux, probably also OSX
- Work as an I2P Plugin OR as a freestanding app to be compatible with all I2P distributions
 - Works on Linux
- Download Tor Browser from an in-I2P mirror(or one of a network of in-I2P mirrors)
 - Not done

### Secondary Goals:

- Launch Tor Browser
 - Works on Linux
- Configure and launch Tor browser for use with I2P
 - Works on Linux

#### Optional Features I might add if there is interest

- Mirror the files which it downloads to an I2P Site
- Mirror the files which it downloads to I2P torrents
- Set up an onion site which announces an I2P mirror exists
- Use Bittorrent-over-I2P to download the Tor Browser software

### Similar Projects:

- https://github.com/micahflee/torbrowser-launcher
- https://github.com/whonix/tb-updater
