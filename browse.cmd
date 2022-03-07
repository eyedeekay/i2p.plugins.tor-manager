:; # This is a script which is valid on Linux, OSX, and Windows.
:; # It will automatically pick the correct path for the correct OS.
:; # It is useful if you want to put copies of this software for
:; # different OSes in the same directory on the same storage device.
:; # Each copy will use the same configuration and working directories,
:; # effectively making it a single, portable installation.
:; #
:; # Such an installation will also be freaking huge, roughly ~400MB,
:; # containing JVM's for every platform and multiple baseline copies of
:; # the static embedded resources.
:; #
:; # It enforces password-protected configuration/working directories by
:; # default.

:; TOR_MANAGER_REQUIRE_PASSWORD=${TOR_MANAGER_REQUIRE_PASSWORD:-true}
:; OS=$(uname -s | tr '[:upper:]' '[:lower:]')
:; ARCH=$(uname -m | sed 's/x86_64/amd64/' | sed 's/i[3-6]86/386/')
:; ./i2p.plugins.tor-manager-"$OS-$ARCH"; exit $?
@ECHO OFF
set TOR_MANAGER_REQUIRE_PASSWORD=true
.\i2p.plugins.tor-manager-windows-"%PROCESSOR_ARCHITECTURE%"
