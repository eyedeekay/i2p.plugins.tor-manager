VERSION=0.0.3
CGO_ENABLED=0
export CGO_ENABLED=0

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
	i2p.plugin.native -name=$(BINARY) \
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

index:
	@echo "<!DOCTYPE html>" > index.html
	@echo "<html>" >> index.html
	@echo "<head>" >> index.html
	@echo "  <title>$(BINARY) - $(CONSOLEPOSTNAME)</title>" >> index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/style.css\" />" >> index.html
	@echo "</head>" >> index.html
	@echo "<body>" >> index.html
	pandoc README.md >> index.html
	@echo "</body>" >> index.html
	@echo "</html>" >> index.html

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