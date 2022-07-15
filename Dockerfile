FROM debian:stable-backports
ENV GOPATH /go
# set the locale to en_US.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8
ENV LC_ALL en_US.UTF-8
# set the timezone to UTC
ENV TZ UTC
ENV APP_ID org.i2pgit.idk.i2p.plugins.tor-manager
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
        make \
        dos2unix \
        curl \
        jq \ 
        openjdk-17-* \
        ant \
        debhelper \
        ant \
        debconf \
        libjetty9-java \
        libservlet3.1-java \
        libtaglibs-standard-jstlel-java \
        libtomcat9-java \
        dh-apparmor \
        bash-completion \
        gettext \
        libgetopt-java \
        libjson-simple-java \
        libgmp-dev \
        libservice-wrapper-java \
        po-debconf \
        geoip-database \
        gettext-base \
        libgetopt-java \
        libjson-simple-java \
        libjson-simple-java \
        libjetty9-java \
        libservlet3.1-java \
        libtaglibs-standard-jstlel-java \
        libtomcat9-java \
        famfamfam-flag-png \
        sensible-utils \
        unzip \
        p7zip-full \
        zenity \
        && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    ln -sf /usr/lib/go-1.17/bin/go /usr/bin/go 
RUN    git clone https://github.com/eyedeekay/go-I2P-jpackage /go/src/github.com/eyedeekay/go-I2P-jpackage
WORKDIR /go/src/github.com/eyedeekay/go-I2P-jpackage
RUN    make || ls i2p.firefox -lah
RUN    git clone https://i2pgit.org/idk/i2p.plugins.tor-manager /go/src/i2pgit.org/idk/i2p.plugins.tor-manager
WORKDIR /go/src/i2pgit.org/idk/i2p.plugins.tor-manager
RUN GOOS=linux GOARCH=amd64 make build
CMD ["/go/src/i2pgit.org/idk/i2p.plugins.tor-manager/i2p.plugins.tor-manager-linux-amd64"]