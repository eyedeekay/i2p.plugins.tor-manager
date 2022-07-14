Docker Usage
============

It is also possible to run this from within a Docker container to avoid
interacting with your base system(for the most part).

```sh
git clone https://i2pgit.org/idk/i2p.plugins.tor-manager.git i2p.plugins.tor-manager && i2p.plugins.tor-manager
docker build -t idk/i2p.plugins.tor-manager .
```

This example shares the X Server with the host, and shares the network with
the host, for the sake of simplicity. A more complex setup might extend this
container with a VNC server and allow the user to connect to it as if it were
a remote server and avoid directly connecting to the host's X Server or some
other means of isolating the browser from the host.

```sh
docker run -it --rm \
    --env="DISPLAY" \
    --volume="$HOME/.Xauthority:/root/.Xauthority:rw" \
    idk/i2p.plugins.tor-manager
```

If you use this without passing the `--net=host` flag, it will always start
a new I2P router inside of the Docker container. If you would like to use an I2P
router running on the host system, then you will need to pass the `--net=host`
flag.
