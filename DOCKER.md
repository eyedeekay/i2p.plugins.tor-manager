Docker Usage
============

```sh
docker build -t eyedeekay/i2p.plugins.tor-manager .
```

```sh
docker run -it --rm \
    --net=host \
    --env="DISPLAY" \
    --volume="$HOME/.Xauthority:/root/.Xauthority:rw" \
    eyedeekay/i2p.plugins.tor-manager
```
