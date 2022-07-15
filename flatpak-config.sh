#! /usr/bin/env sh
export TOR_MANAGER_CLEARNET_MIRROR=true
export TOR_MANAGER_REQUIRE_PASSWORD=false
export TOR_MANAGER_NEVER_USE_TOR=true 
export APP_ID=org.i2pgit.idk.i2p.plugins.tor-manager
/app/bin/i2p.plugins.tor-manager \
    --nevertor \
    --p2p=false \
    --i2pconfig \
    --directory=~/.var/app/org.i2pgit.idk.i2p.plugins.tor-manager