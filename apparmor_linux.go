package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

// Generates files corresponding to:
// https://github.com/micahflee/torbrowser-launcher/blob/develop/apparmor/tunables/torbrowser
func generateTunables() (string, error) {
	var returnable string
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if *verbose {
		log.Println("Generating tunables for apparmor")
	}
	installDir := strings.Replace(client.TBD.BrowserDir(), userHome, "@{HOME}", 1)
	returnable += "@{torbrowser_installation_dir} = " + installDir + "\n"
	returnable += "@{torbrowser_home_dir} = " + installDir + "/Browser\n"
	return returnable, nil
}

func generateTorProfile() (string, error) {
	var returnable string
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if *verbose {
		log.Println("Generating Tor profile for apparmor")
	}

	torBrowserExecutable := strings.Replace(filepath.Join(client.TBD.BrowserDir(), "Browser", "TorBrowser", "Tor", "tor"), userHome, "@{HOME}", 1)

	returnable += `#include <tunables/global>
#include <tunables/torbrowser>

@{torbrowser_tor_executable} = ` + torBrowserExecutable + `

profile torbrowser_tor @{torbrowser_tor_executable} {
	#include <abstractions/base>

	network netlink raw,
	network tcp,
	network udp,

	/etc/host.conf r,
	/etc/nsswitch.conf r,
	/etc/passwd r,
	/etc/resolv.conf r,
	owner @{torbrowser_home_dir}/TorBrowser/Tor/tor mr,
	owner @{torbrowser_home_dir}/TorBrowser/Data/Tor/ rw,
	owner @{torbrowser_home_dir}/TorBrowser/Data/Tor/** rw,
	owner @{torbrowser_home_dir}/TorBrowser/Data/Tor/lock rwk,
	owner @{torbrowser_home_dir}/TorBrowser/Tor/*.so mr,
	owner @{torbrowser_home_dir}/TorBrowser/Tor/*.so.* mr,

	# Support some of the included pluggable transports
	owner @{torbrowser_home_dir}/TorBrowser/Tor/PluggableTransports/** rix,
	@{PROC}/sys/net/core/somaxconn r,
	#include <abstractions/ssl_certs>

	# Silence file_inherit logs
	deny @{torbrowser_home_dir}/{browser/,}omni.ja r,
	deny @{torbrowser_home_dir}/{browser/,}features/*.xpi r,
	deny @{torbrowser_home_dir}/TorBrowser/Data/Browser/profile.default/.parentlock rw,
	deny @{torbrowser_home_dir}/TorBrowser/Data/Browser/profile.default/extensions/*.xpi r,
	deny @{torbrowser_home_dir}/TorBrowser/Data/Browser/profile.default/startupCache/* r,
	# Silence logs from included pluggable transports
	deny /etc/hosts r,
	deny /etc/services r,

	@{PROC}/sys/kernel/random/uuid r,
	/sys/devices/system/cpu/ r,

	# OnionShare compatibility
	/tmp/onionshare/** rw,

	#include <local/torbrowser.Tor.tor>
}`
	return returnable, nil
}

func generateFirefoxProfile() (string, error) {
	var returnable string
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if *verbose {
		log.Println("Generating Firefox profile for apparmor")
	}

	firefoxBrowserExecutable := strings.Replace(filepath.Join(client.TBD.BrowserDir(), "Browser", "firefox.real"), userHome, "@{HOME}", 1)

	returnable += `#include <tunables/global>
#include <tunables/torbrowser>

@{i2pbrowser_home_dir} = ` + strings.Replace(filepath.Join(tbget.UNPACK_PATH(), "i2p.firefox"), userHome, "@{HOME}", 1) + `
@{torbrowser_firefox_executable} = ` + firefoxBrowserExecutable + `

profile torbrowser_firefox @{torbrowser_firefox_executable} {
	#include <abstractions/audio>
	#include <abstractions/dri-enumerate>
	#include <abstractions/gnome>
	#include <abstractions/ibus>
	#include <abstractions/mesa>
	#include <abstractions/opencl>
	#include if exists <abstractions/vulkan>
  
	# Uncomment the following lines if you want to give the Tor Browser read-write
	# access to most of your personal files.
	# #include <abstractions/user-download>
	# @{HOME}/ r,
  
	# Audio support
	/{,usr/}bin/pulseaudio Pixr,
  
	#dbus,
	network netlink raw,
	network tcp,
  
	ptrace (trace) peer=@{profile_name},
	signal (receive, send) set=("term") peer=@{profile_name},
  
	deny /etc/host.conf r,
	deny /etc/hosts r,
	deny /etc/nsswitch.conf r,
	deny /etc/os-release r,
	deny /etc/resolv.conf r,
	deny /etc/passwd r,
	deny /etc/group r,
	deny /etc/mailcap r,
  
	/etc/machine-id r,
	/var/lib/dbus/machine-id r,
  
	/dev/ r,
	/dev/shm/ r,
  
	owner @{PROC}/@{pid}/cgroup r,
	owner @{PROC}/@{pid}/environ r,
	owner @{PROC}/@{pid}/fd/ r,
	owner @{PROC}/@{pid}/mountinfo r,
	owner @{PROC}/@{pid}/stat r,
	owner @{PROC}/@{pid}/status r,
	owner @{PROC}/@{pid}/task/*/stat r,
	@{PROC}/sys/kernel/random/uuid r,
  
	owner @{torbrowser_installation_dir}/ r,
	owner @{torbrowser_installation_dir}/* r,
	owner @{torbrowser_installation_dir}/.** rwk,
	owner @{torbrowser_installation_dir}/update.test/ rwk,
	owner @{torbrowser_home_dir}/.** rwk,
	owner @{torbrowser_home_dir}/ rw,
	owner @{torbrowser_home_dir}/** rwk,
	owner @{torbrowser_home_dir}.bak/ rwk,
	owner @{torbrowser_home_dir}.bak/** rwk,
	owner @{torbrowser_home_dir}/*.so mr,
	owner @{torbrowser_home_dir}/.cache/fontconfig/ rwk,
	owner @{torbrowser_home_dir}/.cache/fontconfig/** rwkl,
	owner @{torbrowser_home_dir}/browser/** r,
	owner @{torbrowser_home_dir}/{,browser/}components/*.so mr,
	owner @{torbrowser_home_dir}/Downloads/ rwk,
	owner @{torbrowser_home_dir}/Downloads/** rwk,
	owner @{torbrowser_home_dir}/firefox rix,
	owner @{torbrowser_home_dir}/{,TorBrowser/UpdateInfo/}updates/[0-9]*/* rw,
	owner @{torbrowser_home_dir}/{,TorBrowser/UpdateInfo/}updates/[0-9]*/{,MozUpdater/bgupdate/}updater ix,
	owner @{torbrowser_home_dir}/updater ix,
	owner @{torbrowser_home_dir}/TorBrowser/Data/Browser/.parentwritetest rw,
	owner @{torbrowser_home_dir}/TorBrowser/Data/Browser/profiles.ini r,
	owner @{torbrowser_home_dir}/TorBrowser/Data/Browser/profile.default/{,**} rwk,
	owner @{torbrowser_home_dir}/TorBrowser/Data/fontconfig/fonts.conf r,
	owner @{torbrowser_home_dir}/fonts/* l,
	owner @{torbrowser_home_dir}/TorBrowser/Tor/tor px,
	owner @{torbrowser_home_dir}/TorBrowser/Tor/ r,
	owner @{torbrowser_home_dir}/TorBrowser/Tor/*.so mr,
	owner @{torbrowser_home_dir}/TorBrowser/Tor/*.so.* mr,
	owner @{torbrowser_home_dir}/TorBrowser/Tor/libstdc++/*.so mr,
	owner @{torbrowser_home_dir}/TorBrowser/Tor/libstdc++/*.so.* mr,
	owner @{i2pbrowser_home_dir}/i2p.firefox/{,**} rwk,
  
	# parent Firefox process when restarting after upgrade, Web Content processes
	owner @{torbrowser_firefox_executable} pxmr -> torbrowser_firefox,
  
	/etc/mailcap r,
	/etc/mime.types r,
  
	/usr/share/ r,
	/usr/share/glib-2.0/schemas/gschemas.compiled r,
	/usr/share/mime/ r,
	/usr/share/themes/ r,
	/usr/share/applications/** rk,
	/usr/share/gnome/applications/ r,
	/usr/share/gnome/applications/kde4/ r,
	/usr/share/poppler/cMap/ r,
  
	# Distribution homepage
	/usr/share/homepage/ r,
	/usr/share/homepage/** r,
  
	/sys/bus/pci/devices/ r,
	@{sys}/devices/pci[0-9]*/**/irq r,
	/sys/devices/system/cpu/ r,
	/sys/devices/system/cpu/present r,
	/sys/devices/system/node/ r,
	/sys/devices/system/node/node[0-9]*/meminfo r,
	/sys/fs/cgroup/cpu,cpuacct/{,user.slice/}cpu.cfs_quota_us r,
	deny /sys/devices/virtual/block/*/uevent r,
  
	# Should use abstractions/gstreamer instead once merged upstream
	/etc/udev/udev.conf r,
	/run/udev/data/+pci:* r,
	/sys/devices/pci[0-9]*/**/uevent r,
	owner /{dev,run}/shm/shmfd-* rw,
  
	# Required for multiprocess Firefox (aka Electrolysis, i.e. e10s)
	owner /{dev,run}/shm/org.chromium.* rw,
	owner /dev/shm/org.mozilla.ipc.[0-9]*.[0-9]* rw, # for Chromium IPC
  
	# Required for Wayland display protocol support
	owner /dev/shm/wayland.mozilla.ipc.[0-9]* rw,
  
	# Silence denial logs about permissions we don't need
	deny @{HOME}/.cache/fontconfig/ rw,
	deny @{HOME}/.cache/fontconfig/** rw,
	deny @{HOME}/.config/gtk-2.0/ rw,
	deny @{HOME}/.config/gtk-2.0/** rw,
	deny @{PROC}/@{pid}/net/route r,
	deny /sys/devices/system/cpu/cpufreq/policy[0-9]*/cpuinfo_max_freq r,
	deny /sys/devices/system/cpu/*/cache/index[0-9]*/size r,
	deny /run/user/[0-9]*/dconf/user rw,
	deny /usr/bin/lsb_release x,
  
	# Silence denial logs about PulseAudio
	deny /etc/pulse/client.conf r,
	deny /usr/bin/pulseaudio x,
  
	# KDE 4
	owner @{HOME}/.kde/share/config/* r,
  
	# Xfce4
	/etc/xfce4/defaults.list r,
	/usr/share/xfce4/applications/ r,
  
	# u2f (tested with Yubikey 4)
	/sys/class/ r,
	/sys/bus/ r,
	/sys/class/hidraw/ r,
	/run/udev/data/c24{5,7,9}:* r,
	/dev/hidraw* rw,
	# Yubikey NEO also needs this:
	/sys/devices/**/hidraw/hidraw*/uevent r,
  
	# Needed for Firefox sandboxing via unprivileged user namespaces
	capability sys_admin,
	capability sys_chroot,
	owner @{PROC}/@{pid}/{gid,uid}_map w,
	owner @{PROC}/@{pid}/setgroups w,
  
	# Remove these rules once we can assume abstractions/vulkan is recent enough
	# to include them
	/etc/glvnd/egl_vendor.d/{*,.json} r,
	/usr/share/glvnd/egl_vendor.d/{,*.json} r,
  
	#include <local/torbrowser.Browser.firefox>
  }
  
 `
	return returnable, nil
}

func GenerateAppArmor() error {
	if tunables, err := generateTunables(); err == nil {
		path := filepath.Join(tbget.UNPACK_PATH(), "tunables.torbrowser.apparmor")
		if *verbose {
			fmt.Printf("Writing %s\n", path)
		}
		if err := ioutil.WriteFile(path, []byte(tunables), 0644); err != nil {
			return err
		}
		if *verbose {
			fmt.Printf("%s", tunables)
		}
	} else {
		return err
	}
	if torProfile, err := generateTorProfile(); err == nil {
		path := filepath.Join(tbget.UNPACK_PATH(), "torbrowser.Tor.tor.apparmor")
		if *verbose {
			fmt.Printf("Writing %s\n", path)
		}
		if err := ioutil.WriteFile(path, []byte(torProfile), 0644); err != nil {
			return err
		}
		if *verbose {
			fmt.Printf("%s", torProfile)
		}
	} else {
		return err
	}
	if firefoxProfile, err := generateFirefoxProfile(); err == nil {
		path := filepath.Join(tbget.UNPACK_PATH(), "torbrowser.Browser.firefox.apparmor")
		if *verbose {
			fmt.Printf("Writing %s\n", path)
		}
		if err := ioutil.WriteFile(path, []byte(firefoxProfile), 0644); err != nil {
			return err
		}
		if *verbose {
			fmt.Printf("%s", firefoxProfile)
		}
	} else {
		return err
	}
	return nil
}
