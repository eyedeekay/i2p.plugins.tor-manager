#! /usr/bin/env sh
export TOR_MANAGER_CLEARNET_MIRROR=true 
export TOR_MANAGER_REQUIRE_PASSWORD=false
export TOR_MANAGER_NEVER_USE_TOR=true
nohup /app/bin/i2p.plugins.tor-manager -nevertor -p2p=false -directory ~/.var/app/org.i2pgit.idk.i2p.plugins.tor-manager 2> ~/.var/app/org.i2pgit.idk.i2p.plugins.tor-manager/tor-manager.err.log 1> ~/.var/app/org.i2pgit.idk.i2p.plugins.tor-manager/tor-manager.log &
/app/bin/i2p.plugins.tor-manager -nevertor -p2p=false -i2pconfig -directory ~/.var/app/org.i2pgit.idk.i2p.plugins.tor-manager