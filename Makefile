VERSION=0.0.4
CGO_ENABLED=0
#export CGO_ENABLED=0

GOOS?=$(shell uname -s | tr A-Z a-z)
GOARCH?="amd64"

ARG=-v -tags netgo -ldflags '-w -extldflags "-static"'

BINARY=i2p.plugins.tor-manager
SIGNER=hankhill19580@gmail.com
CONSOLEPOSTNAME=Tor Binary Manager
USER_GH=eyedeekay
PLUGIN=$(HOME)/.i2p/plugins/$(BINARY)-$(GOOS)-$(GOARCH)

PREFIX?=/usr/local

binary:
	go build $(ARG) -tags="netgo" -o $(BINARY)-$(GOOS)-$(GOARCH) .

lint:
	golint supervise/*.go
	golint get/*.go
	golint serve/*.go

install-binary: binary
	cp -v $(BINARY)-$(GOOS)-$(GOARCH) $(PLUGIN)/lib

install:
	mkdir -p /var/lib/i2pbrowser/icons
	install -m755 -v $(BINARY)-$(GOOS)-$(GOARCH) $(PREFIX)/bin/$(BINARY)-$(GOOS)-$(GOARCH)
	ln -sf $(PREFIX)/bin/$(BINARY)-$(GOOS)-$(GOARCH) $(PREFIX)/bin/i2pbrowser
	ln -sf $(PREFIX)/bin/$(BINARY)-$(GOOS)-$(GOARCH) $(PREFIX)/bin/torbrowser
	install i2ptorbrowser.desktop /usr/share/applications/i2ptorbrowser.desktop
	install torbrowser.desktop /usr/share/applications/torbrowser.desktop
	install garliconion.png /var/lib/i2pbrowser/icons/garliconion.png
	install onion.png /var/lib/i2pbrowser/icons/onion.png

build: dep binary
	
p: dep binary su3

clean:
	rm -f $(BINARY)-plugin plugin $(BINARY)-*zip -r $(BINARY)-$(GOOS)-$(GOARCH) $(BINARY)-$(GOOS)-$(GOARCH).exe tmp tor-browser/torbrowser-*.* $(BINARY) $(BINARY).exe
	rm -f *.su3 *.zip $(BINARY)-$(GOOS)-$(GOARCH) $(BINARY)-*
	git clean -df

all: windows linux osx bsd

windows:
	GOOS=windows GOARCH=amd64 make build su3
	GOOS=windows GOARCH=386 make build su3

linux:
	GOOS=linux GOARCH=amd64 make build su3
	GOOS=linux GOARCH=arm64 make build su3
	GOOS=linux GOARCH=386 make build su3

osx:
	GOOS=darwin GOARCH=amd64 make build su3
	GOOS=darwin GOARCH=arm64 make build su3

bsd:
#	GOOS=freebsd GOARCH=amd64 make build su3
#	GOOS=openbsd GOARCH=amd64 make build su3

dep:
	cp "$(HOME)/build/shellservice.jar" tor-browser/lib/shellservice.jar -v

su3:
	i2p.plugin.native -name=$(BINARY)-$(GOOS)-$(GOARCH) \
		-signer=$(SIGNER) \
		-version "$(VERSION)" \
		-author=$(SIGNER) \
		-autostart=true \
		-clientname=$(BINARY) \
		-consolename="$(BINARY) - $(CONSOLEPOSTNAME)" \
		-delaystart="1" \
		-desc="`cat desc`" \
		-exename=$(BINARY)-$(GOOS)-$(GOARCH) \
		-icondata=icon/icon.png \
		-consoleurl="http://127.0.0.1:7695" \
		-updateurl="http://idk.i2p/$(BINARY)/$(BINARY)-$(GOOS)-$(GOARCH).su3" \
		-website="http://idk.i2p/$(BINARY)/" \
		-command="$(BINARY)-$(GOOS)-$(GOARCH)" \
		-license=MIT \
		-res=tor-browser/
	unzip -o $(BINARY)-$(GOOS)-$(GOARCH).zip -d $(BINARY)-$(GOOS)-$(GOARCH)-zip

sum:
	sha256sum $(BINARY)-$(GOOS)-$(GOARCH).su3

version:
	gothub release -p -u eyedeekay -r $(BINARY) -t "$(VERSION)" -d "`cat desc`"; true

upload:
	gothub upload -R -u eyedeekay -r $(BINARY) -t "$(VERSION)" -f $(BINARY)-$(GOOS)-$(GOARCH).su3 -n $(BINARY)-$(GOOS)-$(GOARCH).su3 -l "`sha256sum $(BINARY)-$(GOOS)-$(GOARCH).su3`"
	gothub upload -R -u eyedeekay -r $(BINARY) -t "$(VERSION)" -f $(BINARY)-$(GOOS)-$(GOARCH) -n $(BINARY)-$(GOOS)-$(GOARCH) -l "`sha256sum $(BINARY)-$(GOOS)-$(GOARCH)`"

upload-windows:
	GOOS=windows GOARCH=amd64 make upload
	GOOS=windows GOARCH=386 make upload

upload-linux:
	GOOS=linux GOARCH=amd64 make upload
	GOOS=linux GOARCH=arm64 make upload
	GOOS=linux GOARCH=386 make upload

upload-osx:
	GOOS=darwin GOARCH=amd64 make upload
	GOOS=darwin GOARCH=arm64 make upload

upload-bsd:
#	GOOS=freebsd GOARCH=amd64 make upload
#	GOOS=openbsd GOARCH=amd64 make upload

upload-all: upload-windows upload-linux upload-osx upload-bsd

download-su3s:
	GOOS=windows GOARCH=amd64 make download-single-su3
	GOOS=windows GOARCH=386 make download-single-su3
	GOOS=linux GOARCH=amd64 make download-single-su3
	GOOS=linux GOARCH=arm64 make download-single-su3
	GOOS=linux GOARCH=386 make download-single-su3
	GOOS=darwin GOARCH=amd64 make download-single-su3
	GOOS=darwin GOARCH=arm64 make download-single-su3
#	GOOS=freebsd GOARCH=amd64 make download-single-su3
#	GOOS=openbsd GOARCH=amd64 make download-single-su3

download-single-su3:
	wget -N -c "https://github.com/$(USER_GH)/$(BINARY)/releases/download/$(VERSION)/$(BINARY)-$(GOOS)-$(GOARCH).su3"

early-release: clean linux windows version upload-linux upload-windows

release: clean all version upload-all

index: index-clearnet index-offline index-usage
	@echo "<!DOCTYPE html>" > index.html
	@echo "<html>" >> index.html
	@echo "<head>" >> index.html
	@echo "  <title>$(BINARY) - $(CONSOLEPOSTNAME)</title>" >> index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/style.css\" />" >> index.html
	@echo "</head>" >> index.html
	@echo "<body>" >> index.html
	sed 's|https://github.com/eyedeekay/i2p.plugins.tor-manager/releases/download/||g' README.md | \
		sed "s|$(VERSION)||g" | pandoc >> index.html
	@echo "</body>" >> index.html
	@echo "</html>" >> index.html

index-clearnet:
	@echo "<!DOCTYPE html>" > firefox.html
	@echo "<html>" >> firefox.html
	@echo "<head>" >> firefox.html
	@echo "  <title>$(BINARY) - $(CONSOLEPOSTNAME)</title>" >> firefox.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/style.css\" />" >> firefox.html
	@echo "</head>" >> firefox.html
	@echo "<body>" >> firefox.html
	pandoc FIREFOX.md >> firefox.html
	@echo "</body>" >> firefox.html
	@echo "</html>" >> firefox.html

index-offline:
	@echo "<!DOCTYPE html>" > offline.html
	@echo "<html>" >> offline.html
	@echo "<head>" >> offline.html
	@echo "  <title>$(BINARY) - $(CONSOLEPOSTNAME)</title>" >> offline.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/style.css\" />" >> offline.html
	@echo "</head>" >> offline.html
	@echo "<body>" >> offline.html
	pandoc OFFLINE.md >> offline.html
	@echo "</body>" >> offline.html
	@echo "</html>" >> offline.html

tor-browser/unpack/i2p.firefox:
	@echo "TODO"

tor-browser/unpack/i2p.firefox.config:
	@echo "TODO"

refresh-tor-keys: clean-tor-keys tor-browser/TPO-signing-key.pub

tor-keys: tor-browser/TPO-signing-key.pub

clean-tor-keys:
	rm -f tor-browser/TPO-signing-key.pub

tor-browser/TPO-signing-key.pub:
	gpg --armor --output ./tor-browser/TPO-signing-key.pub --export 0xEF6E286DDA85EA2A4BA7DE684E2C6E8793298290

deb: clean
	mv "hankhill19580_at_gmail.com.crl" ../; true
	mv "hankhill19580_at_gmail.com.crt" ../; true
	mv "hankhill19580_at_gmail.com.pem" ../; true
	rm ../i2p.plugins.tor-manager_$(VERSION).orig.tar.gz -f
	tar --exclude=".git" \
		--exclude="hankhill19580_at_gmail.com.crl" \
		--exclude="hankhill19580_at_gmail.com.crt" \
		--exclude="hankhill19580_at_gmail.com.pem" \
		--exclude="i2p.plugins.tor-manager" \
		--exclude="i2p.plugins.tor-manager.exe" \
		--exclude="tmp" \
		-cvzf ../i2p.plugins.tor-manager_$(VERSION).orig.tar.gz	.
	dpkg-buildpackage -us -uc
	mv "../hankhill19580_at_gmail.com.crl" ./
	mv "../hankhill19580_at_gmail.com.crt" ./
	mv "../hankhill19580_at_gmail.com.pem" ./

debsrc: clean
	mv "hankhill19580_at_gmail.com.crl" ../; true
	mv "hankhill19580_at_gmail.com.crt" ../; true
	mv "hankhill19580_at_gmail.com.pem" ../; true
	rm ../i2p.plugins.tor-manager_$(VERSION).orig.tar.gz -f
	tar --exclude=".git" \
		--exclude="hankhill19580_at_gmail.com.crl" \
		--exclude="hankhill19580_at_gmail.com.crt" \
		--exclude="hankhill19580_at_gmail.com.pem" \
		--exclude="i2p.plugins.tor-manager" \
		--exclude="i2p.plugins.tor-manager.exe" \
		--exclude="tmp" \
		-cvzf ../i2p.plugins.tor-manager_$(VERSION).orig.tar.gz	.
	debuild -S
	mv "../hankhill19580_at_gmail.com.crl" ./
	mv "../hankhill19580_at_gmail.com.crt" ./
	mv "../hankhill19580_at_gmail.com.pem" ./

DATE=`date +%Y/%m/%d`

usage:
	./i2p.plugins.tor-manager --help

usagemd:
	@echo "Tor(And sometimes Firefox) Manager for I2P" | tee USAGE.md
	@echo "===========================================" | tee -a USAGE.md
	@echo "" | tee -a USAGE.md
	@echo "## Usage: $(BINARY) [options]" | tee -a USAGE.md
	@echo "" | tee -a USAGE.md
	@echo "### Options:" | tee -a USAGE.md
	@echo "" | tee -a USAGE.md
	@echo '```sh' | tee -a USAGE.md
	./i2p.plugins.tor-manager --help 2>&1 | grep -v $(DATE) | grep -v $(HOME) | tee -a USAGE.md
	@echo '```' | tee -a USAGE.md
	@echo "" | tee -a USAGE.md

index-usage:
	@echo "<!DOCTYPE html>" > usage.html
	@echo "<html>" >> usage.html
	@echo "<head>" >> usage.html
	@echo "  <title>$(BINARY) - $(CONSOLEPOSTNAME)</title>" >> usage.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/style.css\" />" >> usage.html
	@echo "</head>" >> usage.html
	@echo "<body>" >> usage.html
	pandoc USAGE.md >> usage.html
	@echo "</body>" >> usage.html
	@echo "</html>" >> usage.html



#GO111MODULE=off
#export GO111MODULE=off

i2pd_prerelease_version=c-wrapper-libi2pd-api
i2pd_release_version=2.40.0

export GOPATH=$(HOME)/go

export USE_STATIC=yes
USE_STATIC=yes

export LDFLAGS=-static
LDFLAGS=-static

GXXFLAGS=-static
export GXXFLAGS=-static

CXXFLAGS=-static
export CXXFLAGS=-static

CGO_GXXFLAGS=-static
export CGO_GXXFLAGS=-static

CGO_CFLAGS=-static
export CGO_CFLAGS=-static

CGO_CXXFLAGS=-static
export CGO_CXXFLAGS=-static

CGO_CPPFLAGS=-static
export CGO_CPPFLAGS=-static

#CGO_LDFLAGS=-static
#export CGO_LDFLAGS=-static


#Trying to achieve fully-static builds, this doesn't work yet.
FLAGS=/usr/lib/x86_64-linux-gnu/libboost_system.a /usr/lib/x86_64-linux-gnu/libboost_date_time.a /usr/lib/x86_64-linux-gnu/libboost_filesystem.a /usr/lib/x86_64-linux-gnu/libboost_program_options.a /usr/lib/x86_64-linux-gnu/libssl.a /usr/lib/x86_64-linux-gnu/libcrypto.a /usr/lib/x86_64-linux-gnu/libz.a

example:
	go build -x -v --tags=netgo \
		-ldflags '-w -linkmode=external -extldflags "-static -ldl $(FLAGS)"'