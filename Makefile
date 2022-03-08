VERSION=0.0.5
#CGO_ENABLED=0
#export CGO_ENABLED=0
export PKG_CONFIG_PATH=/usr/lib/$(uname -m)-linux-musl/pkgconfig

GOOS?=$(shell uname -s | tr A-Z a-z)
GOARCH?="amd64"

ARG=-v -tags netgo -ldflags '-w' # -extldflags "-static"'
#FLAGS=/usr/lib/x86_64-linux-gnu/libboost_system.a /usr/lib/x86_64-linux-gnu/libboost_date_time.a /usr/lib/x86_64-linux-gnu/libboost_filesystem.a /usr/lib/x86_64-linux-gnu/libboost_program_options.a /usr/lib/x86_64-linux-gnu/libssl.a /usr/lib/x86_64-linux-gnu/libcrypto.a /usr/lib/x86_64-linux-gnu/libz.a
#ARG=-ldflags '-w -linkmode=external -extldflags "-static -ldl $(FLAGS)"'
#NOSTATIC=-v -tags netgo -ldflags '-w -extldflags "-ldl $(FLAGS)"'
WINGUI=-ldflags '-H=windowsgui'

BINARY=i2p.plugins.tor-manager
SIGNER=hankhill19580@gmail.com
CONSOLEPOSTNAME=Tor Binary Manager
USER_GH=eyedeekay
PLUGIN=$(HOME)/.i2p/plugins/$(BINARY)-$(GOOS)-$(GOARCH)

PREFIX?=/usr/local

binary:
	go build $(ARG) -tags="netgo osusergo systray" -o $(BINARY)-$(GOOS)-$(GOARCH) .

winbinary:
	CC=/usr/bin/x86_64-w64-mingw32-gcc \
		CXX=/usr/bin/x86_64-w64-mingw32-g++ \
		GOOS=windows go build $(WINGUI) -tags="netgo osusergo systray" -o $(BINARY)-$(GOOS)-$(GOARCH) .

nosystray:
	go build $(NOSTATIC) -tags="netgo osusergo nosystray" -o $(BINARY)-$(GOOS)-$(GOARCH) .

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

winbuild: dep winbinary
	
p: dep binary su3

clean:
	rm -f $(BINARY)-plugin plugin $(BINARY)-*zip -r $(BINARY)-$(GOOS)-$(GOARCH) $(BINARY)-$(GOOS)-$(GOARCH).exe tmp tor-browser/torbrowser-*.* $(BINARY) $(BINARY).exe
	rm -f *.su3 *.zip $(BINARY)-$(GOOS)-$(GOARCH) $(BINARY)-*
	git clean -df

all: clean windows linux osx bsd portable.zip

portable.zip:
	zip -r portable.zip browse.cmd README-PORTABLE.txt \
		$(BINARY)-linux-amd64 \
		$(BINARY)-windows-amd64 \
		#$(BINARY)-darwin-amd64 \
		#$(BINARY)-darwin-arm64 \

backup-embed:
	mkdir -p ../../../github.com/eyedeekay/go-I2P-jpackage.bak
	cp ../../../github.com/eyedeekay/go-I2P-jpackage/* ../../../github.com/eyedeekay/go-I2P-jpackage.bak -r;true
	rm -f ../../../github.com/eyedeekay/go-I2P-jpackage/*.tar.xz
	tar -cvJf ../../../github.com/eyedeekay/go-I2P-jpackage/build.windows.I2P.tar.xz README.md LICENSE
	tar -cvJf ../../../github.com/eyedeekay/go-I2P-jpackage/build.linux.I2P.tar.xz README.md LICENSE

unbackup-embed:
	cp ../../../github.com/eyedeekay/go-I2P-jpackage.bak/*.tar.xz ../../../github.com/eyedeekay/go-I2P-jpackage/

unembed-windows:
	mv ../../../github.com/eyedeekay/go-I2P-jpackage/build.windows.I2P.tar.xz ../../../github.com/eyedeekay/
	tar -cvJf ../../../github.com/eyedeekay/go-I2P-jpackage/build.windows.I2P.tar.xz README.md LICENSE

unembed-linux:
	mv ../../../github.com/eyedeekay/go-I2P-jpackage/build.linux.I2P.tar.xz ../../../github.com/eyedeekay/
	tar -cvJf ../../../github.com/eyedeekay/go-I2P-jpackage/build.linux.I2P.tar.xz README.md LICENSE

winplugin: 
	GOOS=windows make backup-embed build unbackup-embed

linplugin: 
	GOOS=linux make backup-embed build unbackup-embed

linplugin-nosystray: 
	GOOS=linux make backup-embed nosystray unbackup-embed

osxplugin:
	GOOS=darwin make backup-embed nosystray unbackup-embed

windows:
	GOOS=windows GOARCH=amd64 make winplugin su3 unembed-linux build unbackup-embed
	GOOS=windows GOARCH=386 make winplugin su3 unembed-linux build unbackup-embed

linux:
	GOOS=linux GOARCH=amd64 make linplugin su3 unembed-windows build unbackup-embed
#	GOOS=linux GOARCH=arm64 make linplugin su3 unembed-windows build unbackup-embed
	PKG_CONFIG_PATH=/usr/lib/i386-linux-gnu/pkgconfig GOOS=linux GOARCH=386 make linplugin-nosystray su3 unembed-windows nosystray unbackup-embed

osx:
	GOOS=darwin GOARCH=amd64 make osxplugin su3 unembed-windows unembed-linux nosystray unbackup-embed
	GOOS=darwin GOARCH=arm64 make osxplugin su3 unembed-windows unembed-linux nosystray unbackup-embed

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
	gothub release -u eyedeekay -r $(BINARY) -t "$(VERSION)" -d "`cat desc`"; true

upload:
	gothub upload -R -u eyedeekay -r $(BINARY) -t "$(VERSION)" -f $(BINARY)-$(GOOS)-$(GOARCH).su3 -n $(BINARY)-$(GOOS)-$(GOARCH).su3 -l "`sha256sum $(BINARY)-$(GOOS)-$(GOARCH).su3`"
	gothub upload -R -u eyedeekay -r $(BINARY) -t "$(VERSION)" -f $(BINARY)-$(GOOS)-$(GOARCH) -n $(BINARY)-$(GOOS)-$(GOARCH) -l "`sha256sum $(BINARY)-$(GOOS)-$(GOARCH)`"

upload-portable.zip:
	gothub upload -R -u eyedeekay -r $(BINARY) -t "$(VERSION)" -f portable.zip -n i2pbrowser.portable.zip -l "`sha256sum portable.zip`"

upload-windows:
	GOOS=windows GOARCH=amd64 make upload
	GOOS=windows GOARCH=386 make upload

upload-linux:
	GOOS=linux GOARCH=amd64 make upload
#	GOOS=linux GOARCH=arm64 make upload
	GOOS=linux GOARCH=386 make upload

upload-osx:
	GOOS=darwin GOARCH=amd64 make upload
	GOOS=darwin GOARCH=arm64 make upload

upload-bsd:
#	GOOS=freebsd GOARCH=amd64 make upload
#	GOOS=openbsd GOARCH=amd64 make upload

upload-all: upload-windows upload-linux upload-osx upload-bsd upload-portable-zip

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

index: index-clearnet index-offline index-usage index-onion
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

index-onion:
	@echo "<!DOCTYPE html>" > onion/www/index.html
	@echo "<html>" >> onion/www/index.html
	@echo "<head>" >> onion/www/index.html
	@echo "  <title>$(BINARY) - $(CONSOLEPOSTNAME)</title>" >> onion/www/index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/style.css\" />" >> onion/www/index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/default.css\" />" >> onion/www/index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/desktop.css\" />" >> onion/www/index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/mobile.css\" />" >> onion/www/index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/syntax.css\" />" >> onion/www/index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/widescreen.rtl.css\" />" >> onion/www/index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/default.rtl.css\" />" >> onion/www/index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/desktop.rtl.css\" />" >> onion/www/index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/reset.css\" />" >> onion/www/index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/widescreen.css\" />" >> onion/www/index.html
	@echo "</head>" >> onion/www/index.html
	@echo "<body>" >> onion/www/index.html
	pandoc ONION.md >> onion/www/index.html
	@echo "</body>" >> onion/www/index.html
	@echo "</html>" >> onion/www/index.html

tor-browser/unpack/i2p.firefox:
	@echo "TODO"

tor-browser/unpack/i2p.firefox.config:
	@echo "TODO"

refresh-tor-keys: clean-tor-keys tor-browser/TPO-signing-key.pub

tor-keys: tor-browser/TPO-signing-key.pub

clean-tor-keys:
	rm -f tor-browser/TPO-signing-key.pub

tor-browser/TPO-signing-key.pub:
	#gpg --output ./tor-browser/TPO-signing-key.pub --export -r torbrowser@torproject.org
	#gpg --armor --output ./tor-browser/TPO-signing-key.pub --export -r torbrowser@torproject.org
	#gpg -r 0xEF6E286DDA85EA2A4BA7DE684E2C6E8793298290 --output ./tor-browser/TPO-signing-key.pub --export 
	#gpg -r 0xEF6E286DDA85EA2A4BA7DE684E2C6E8793298290 --armor --output ./tor-browser/TPO-signing-key.pub --export 

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

example:
	go build -x -v --tags=netgo,nosystray \
		-ldflags '-w -linkmode=external -extldflags "-static -ldl $(FLAGS)"'

xhost:
	xhost + local:docker

docker: xhost
	docker build -t eyedeekay/i2p.plugins.tor-manager .
	docker run it --volume $(PWD)/build:/go/src/eyedeekay/i2p.plugins.tor-manager/build \
		--env-file $(HOME)/i2p.plugins.tor-manager/env \
		--publish 127.0.0.1:7695:7695 \
		-e DISPLAY=unix$(DISPLAY) \
		-v /tmp/.X11-unix:/tmp/.X11-unix \
		--rm eyedeekay/i2p.plugins.tor-manager
