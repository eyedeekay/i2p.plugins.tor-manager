FROM debian:stable-backports
RUN echo "deb http://deb.debian.org/debian oldstable main" >> /etc/apt/sources.list &&  \
    apt-get update && apt-get dist-upgrade -y && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        gnupg2 \
        software-properties-common \
        wget \
        lib*appindicator* \
        golang-1.17-go \
        gcc \
        make \
        git \
        xz-utils \
        tar \
        && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    ln -sf /usr/lib/go-1.17/bin/go /usr/bin/go
WORKDIR /go/src/i2pgit.org/idk/i2p.plugins.tor-manager
#CMD ls /go/src/i2pgit.org/idk/i2p.plugins.tor-manager /go/src/github.com/eyedeekay/go-I2P-jpackage
CMD make linux 
#&& /usr/lib/go-1.17/bin/go build