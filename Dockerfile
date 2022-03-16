FROM debian:stable-backports
RUN apt-get update && apt-get dist-upgrade -y && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        gnupg2 \
        software-properties-common \
        wget \
        lib*appindicator*dev \
        golang-1.16-go \
        make \
        && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN which go