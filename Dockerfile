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
        make \
        nsis* \
        dos2unix \
        curl \
        jq \ 
        openjdk-17-* \
        ant \
        debhelper \
        ant \
        debconf \
        default-jdk \
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
        && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    ln -sf /usr/lib/go-1.17/bin/go /usr/bin/go && \
    wget https://nsis.sourceforge.io/mediawiki/images/c/c7/ShellExecAsUser.zip && \
    wget https://nsis.sourceforge.io/mediawiki/images/1/1d/ShellExecAsUserUnicodeUpdate.zip && \
    wget https://nsis.sourceforge.io/mediawiki/images/6/68/ShellExecAsUser_amd64-Unicode.7z && \
    # unzip them to /usr/share/nsis/Plugins/
    unzip ShellExecAsUser.zip -d /usr/share/nsis/Plugins/x86-ansi && \
    unzip ShellExecAsUserUnicodeUpdate.zip -d /usr/share/nsis/Plugins/x86-unicode && \
    7zr x ShellExecAsUser_amd64-Unicode.7z -o/usr/share/nsis/Plugins/amd64-unicode
RUN git clone https://github.com/eyedeekay/go-I2P-jpackage /go/src/github.com/eyedeekay/go-I2P-jpackage && \
    cd /go/src/github.com/eyedeekay/go-I2P-jpackage && \
    touch /go/src/github.com/eyedeekay/go-I2P-jpackage/build.windows.I2P.tar.xz && \
    touch /go/src/github.com/eyedeekay/go-I2P-jpackage/build.linux.I2P.tar.xz && \
    go generate && \
    git clone https://i2pgit.org/idk/i2p.plugins.tor-manager /go/src/i2pgit.org/idk/i2p.plugins.tor-manager
WORKDIR /go/src/i2pgit.org/idk/i2p.plugins.tor-manager
#CMD ls /go/src/i2pgit.org/idk/i2p.plugins.tor-manager /go/src/github.com/eyedeekay/go-I2P-jpackage
RUN GOOS=linux GOARCH=amd64 make build
#&& /usr/lib/go-1.17/bin/go build
CMD ["/go/src/i2pgit.org/idk/i2p.plugins.tor-manager/i2p.plugins.tor-manager-linux-amd64"]