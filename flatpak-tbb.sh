#! /usr/bin/env sh
export TOR_MANAGER_CLEARNET_MIRROR=true
export TOR_MANAGER_REQUIRE_PASSWORD=false 
export APP_ID=org.i2pgit.idk.i2p.plugins.tor-manager
/app/bin/i2p.plugins.tor-manager \
    --p2p=false \
    --torbrowser \
    --directory=~/.var/app/org.i2pgit.idk.i2p.plugins.tor-manager