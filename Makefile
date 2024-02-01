VERSION=0.0.14
CGO_ENABLED=0
export CGO_ENABLED=0
export PKG_CONFIG_PATH=/usr/lib/$(uname -m)-linux-musl/pkgconfig

GOOS?=$(shell uname -s | tr A-Z a-z)
GOARCH?="amd64"

ARG=-v -tags "netgo osusergo " -ldflags '-w -extldflags "-static"'
#FLAGS=/usr/lib/x86_64-linux-gnu/libboost_system.a /usr/lib/x86_64-linux-gnu/libboost_date_time.a /usr/lib/x86_64-linux-gnu/libboost_filesystem.a /usr/lib/x86_64-linux-gnu/libboost_program_options.a /usr/lib/x86_64-linux-gnu/libssl.a /usr/lib/x86_64-linux-gnu/libcrypto.a /usr/lib/x86_64-linux-gnu/libz.a
STATIC=-v -tags "netgo osusergo " -ldflags '-w -extldflags "-static"'
OSXFLAGS=-v -tags "netgo osusergo systray" -ldflags '-w -extldflags "-static"' 
# -gcflags='-DDARWIN -x objective-c -fobjc-arc -ldflags=framework=Cocoa'
#NOSTATIC=-v -tags netgo -ldflags '-w -extldflags "-ldl $(FLAGS)"'
WINGUI=-ldflags '-H=windowsgui'

BINARY=i2p.plugins.tor-manager
SIGNER=hankhill19580@gmail.com
CONSOLEPOSTNAME=Tor Binary Manager
USER_GH=eyedeekay
PLUGIN=$(HOME)/.i2p/plugins/$(BINARY)-$(GOOS)-$(GOARCH)

PREFIX?=/usr/local

version-file:
	@echo "package main" | tee version.go
	@echo "" | tee -a version.go
	@echo 'import (' | tee -a version.go
	@echo '	"fmt"' | tee -a version.go
	@echo '	"os"' | tee -a version.go
	@echo ')' | tee -a version.go
	@echo "" | tee -a version.go
	@echo 'var VERSION string = "$(VERSION)"' | tee -a version.go
	@echo "" | tee -a version.go
	@echo "func printversion() {" | tee -a version.go
	@echo "	fmt.Fprintf(os.Stdout, VERSION)" | tee -a version.go
	@echo "}" | tee -a version.go
	#@echo "" | tee -a version.go

binary:
	go build $(ARG) -tags="netgo osusergo systray" -o $(BINARY)-$(GOOS)-$(GOARCH) .

winbinary:
	CC=/usr/bin/x86_64-w64-mingw32-gcc \
		CXX=/usr/bin/x86_64-w64-mingw32-g++ \
		GOOS=windows go build $(WINGUI) -tags="netgo osusergo systray" -o $(BINARY)-$(GOOS)-$(GOARCH) .

osxsystray:
	export CGO_ENABLED=1 && \
		export GOOS=darwin && \
		export GOARCH=amd64 && \
		make raw

raw:
	/usr/bin/go build $(OSXFLAGS) -o $(BINARY)-$(GOOS)-$(GOARCH)-static .

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

uninstall:
	rm -rf /var/lib/i2pbrowser/icons
	rm -vf $(PREFIX)/bin/$(BINARY)-$(GOOS)-$(GOARCH) \
		$(PREFIX)/bin/i2pbrowser \
		$(PREFIX)/bin/torbrowser \
		/usr/share/applications/i2ptorbrowser.desktop \
		/usr/share/applications/torbrowser.desktop \
		/var/lib/i2pbrowser/icons/garliconion.png \
		/var/lib/i2pbrowser/icons/onion.png

build: dep binary

winbuild: dep winbinary
	
p: dep binary su3

clean: clean-flatpak
	rm -f $(BINARY)-plugin plugin $(BINARY)-*zip -r $(BINARY)-$(GOOS)-$(GOARCH) $(BINARY)-$(GOOS)-$(GOARCH).exe tmp tor-browser/torbrowser-*.* $(BINARY) $(BINARY).exe tmp-i2pbrowser
	rm -f *.su3 *.zip $(BINARY)-$(GOOS)-$(GOARCH) $(BINARY)-*
	git clean -df

all: clean windows linux osx bsd portable.zip appimage flatpak

portable.zip:
	zip -r portable.zip browse.cmd README-PORTABLE.txt \
		$(BINARY)-linux-amd64 \
		$(BINARY)-windows-amd64 \
		#$(BINARY)-darwin-amd64 \
		#$(BINARY)-darwin-arm64 \

backup-embed:
#	mkdir -p ../../../github.com/eyedeekay/go-I2P-jpackage.bak
#	cp ../../../github.com/eyedeekay/go-I2P-jpackage/* ../../../github.com/eyedeekay/go-I2P-jpackage.bak -r;true
#	rm -f ../../../github.com/eyedeekay/go-I2P-jpackage/*.tar.xz
#	tar -cvJf ../../../github.com/eyedeekay/go-I2P-jpackage/build.windows.I2P.tar.xz README.md LICENSE
#	tar -cvJf ../../../github.com/eyedeekay/go-I2P-jpackage/build.linux.I2P.tar.xz README.md LICENSE

unbackup-embed:
#	cp ../../../github.com/eyedeekay/go-I2P-jpackage.bak/*.tar.xz ../../../github.com/eyedeekay/go-I2P-jpackage/

unembed-windows:
#	mv ../../../github.com/eyedeekay/go-I2P-jpackage/build.windows.I2P.tar.xz ../../../github.com/eyedeekay/
#	tar -cvJf ../../../github.com/eyedeekay/go-I2P-jpackage/build.windows.I2P.tar.xz README.md LICENSE

unembed-linux:
#	mv ../../../github.com/eyedeekay/go-I2P-jpackage/build.linux.I2P.tar.xz ../../../github.com/eyedeekay/
#	tar -cvJf ../../../github.com/eyedeekay/go-I2P-jpackage/build.linux.I2P.tar.xz README.md LICENSE

winplugin: 
	GOOS=windows make backup-embed build unbackup-embed

linplugin: 
	GOOS=linux make backup-embed build unbackup-embed

osxplugin:
	GOOS=darwin make backup-embed osxsystray unbackup-embed

windows:
	GOOS=windows GOARCH=amd64 make winplugin su3 unembed-linux build unbackup-embed
	GOOS=windows GOARCH=386 make winplugin su3 unembed-linux build unbackup-embed

linux:
	GOOS=linux GOARCH=amd64 make build su3
#	linplugin su3 unembed-windows build unbackup-embed
#	GOOS=linux GOARCH=arm64 make linplugin su3 unembed-windows build unbackup-embed
	PKG_CONFIG_PATH=/usr/lib/i386-linux-gnu/pkgconfig GOOS=linux GOARCH=386 make su3 unembed-windows  unbackup-embed

osx:
	GOOS=darwin GOARCH=amd64 make osxplugin su3 unembed-windows unembed-linux  unbackup-embed
	GOOS=darwin GOARCH=arm64 make osxplugin su3 unembed-windows unembed-linux  unbackup-embed

bsd:
#	GOOS=freebsd GOARCH=amd64 make build su3
#	GOOS=openbsd GOARCH=amd64 make build su3

dep:
#	#cp "$(HOME)/build/shellservice.jar" tor-browser/lib/shellservice.jar -v

SIGNER_DIR=$(HOME)/i2p-go-keys/

su3:
	i2p.plugin.native -name=$(BINARY)-$(GOOS)-$(GOARCH) \
		-signer=$(SIGNER) \
		-signer-dir=$(SIGNER_DIR) \
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
	#unzip -o $(BINARY)-$(GOOS)-$(GOARCH).zip -d $(BINARY)-$(GOOS)-$(GOARCH)-zip

su3-mirror:
	i2p.plugin.native -name=$(BINARY)-$(GOOS)-$(GOARCH)-Mirrorkit \
		-signer=$(SIGNER) \
		-signer-dir=$(SIGNER_DIR) \
		-version "$(VERSION)" \
		-author=$(SIGNER) \
		-autostart=true \
		-clientname=$(BINARY)-Mirrorkit \
		-consolename="$(BINARY) - $(CONSOLEPOSTNAME) - Mirrorkit" \
		-delaystart="1" \
		-desc="`cat desc` - this is the automatic Tor Browser mirror-generator" \
		-exename=$(BINARY)-$(GOOS)-$(GOARCH) \
		-icondata=icon/icon.png \
		-consoleurl="http://127.0.0.1:7695" \
		-updateurl="http://idk.i2p/$(BINARY)/$(BINARY)-$(GOOS)-$(GOARCH)-Mirrorkit.su3" \
		-website="http://idk.i2p/$(BINARY)/" \
		-command="$(BINARY)-$(GOOS)-$(GOARCH) -notor -nevertor -mirrorall -bemirror" \
		-license=MIT \
		-res=tor-browser/
	unzip -o $(BINARY)-$(GOOS)-$(GOARCH)-Mirrorkit.zip -d $(BINARY)-$(GOOS)-$(GOARCH)-Mirrorkit-zip

sum:
	sha256sum $(BINARY)-$(GOOS)-$(GOARCH).su3

version:
	gothub release -u eyedeekay -r $(BINARY) -t "$(VERSION)" -d "`cat desc`"; true
	sleep 2s

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

upload-all: upload-windows upload-linux upload-osx upload-bsd upload-portable-zip appimage-upload flatpak-upload

download-su3s:

early-release: clean linux windows version upload-linux upload-windows

release: clean all version upload-all

README: pluginslist
	cat top.md plugins.md bottom.md > README.md

pluginslist:
	@echo "" > plugins.md
	@echo "Plugin:" >> plugins.md
	@echo "-------" >> plugins.md
	@echo "" >> plugins.md
	@echo "Clearnet visitor? You'll need to use the [Github Releases Mirror](https://github.com/eyedeekay/i2p.plugins.tor-manager/releases/$(VERSION))." >> plugins.md
	@echo "" >> plugins.md
	@echo "- [i2p.plugins.tor-manager-linux-386](i2p.plugins.tor-manager-linux-386.su3)" >> plugins.md
	@echo "- [i2p.plugins.tor-manager-windows-amd64](i2p.plugins.tor-manager-windows-amd64.su3)" >> plugins.md
	@echo "- [i2p.plugins.tor-manager-darwin-arm64](i2p.plugins.tor-manager-darwin-arm64.su3)" >> plugins.md
	@echo "- [i2p.plugins.tor-manager-linux-amd64](i2p.plugins.tor-manager-linux-amd64.su3)" >> plugins.md
	@echo "- [i2p.plugins.tor-manager-windows-386](i2p.plugins.tor-manager-windows-386.su3)" >> plugins.md
	@echo "- [i2p.plugins.tor-manager-darwin-amd64](i2p.plugins.tor-manager-darwin-amd64.su3)" >> plugins.md
	@echo "" >> plugins.md


index: README index-clearnet index-offline index-usage index-onion
	@echo "<!DOCTYPE html>" > index.html
	@echo "<html>" >> index.html
	@echo "<head>" >> index.html
	@echo "  <title>$(BINARY) - $(CONSOLEPOSTNAME)</title>" >> index.html
	@echo "  <link rel=\"stylesheet\" type=\"text/css\" href =\"/style.css\" />" >> index.html
	@echo "</head>" >> index.html
	@echo "<body>" >> index.html
	sed 's|https://github.com/eyedeekay/i2p.plugins.tor-manager/releases/download/||g' README.md | pandoc >> index.html
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

deb: clean
	./changelog.sh
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

debsrc: clean
	rm ../i2p.plugins.tor-manager_$(VERSION).orig.tar.gz -f
	tar --exclude=".git" \
		--exclude="hankhill19580_at_gmail.com.crl" \
		--exclude="hankhill19580_at_gmail.com.crt" \
		--exclude="hankhill19580_at_gmail.com.pem" \
		--exclude="i2p.plugins.tor-manager" \
		--exclude="i2p.plugins.tor-manager.exe" \
		--exclude="tmp" \
		--exclude="repo" \
		--exclude="flatpak.repo.i2p.plugins.tor-manager" \
		-cvzf ../i2p.plugins.tor-manager_$(VERSION).orig.tar.gz	.
	debuild -S

DATE=`date +%Y/%m/%d`

usage:
	TOR_MANAGER_CLEARNET_MIRROR=true TOR_MANAGER_REQUIRE_PASSWORD=false ./i2p.plugins.tor-manager --p2p=false --help=true

usagemd:
	@echo "Tor(And sometimes Firefox) Manager for I2P" | tee USAGE.md
	@echo "===========================================" | tee -a USAGE.md
	@echo "" | tee -a USAGE.md
	@echo "## Usage: $(BINARY) [options]" | tee -a USAGE.md
	@echo "" | tee -a USAGE.md
	@echo "### Options:" | tee -a USAGE.md
	@echo "" | tee -a USAGE.md
	@echo '```sh' | tee -a USAGE.md
	TOR_MANAGER_CLEARNET_MIRROR=true TOR_MANAGER_REQUIRE_PASSWORD=false ./i2p.plugins.tor-manager --p2p=false --help=true | tee -a USAGE.md
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
	go build -x -v --tags=netgo, \
		-ldflags '-w -linkmode=external -extldflags "-static -ldl $(FLAGS)"'

xhost:
	xhost + local:docker

docker: clean xhost
	docker build -t eyedeekay/i2p.plugins.tor-manager .
	sudo rm -rfv $(PWD).docker-build
	cp -rv $(PWD) $(PWD).docker-build
	cp -v $(HOME)/go/bin/i2p.plugin.native ./i2p.plugin.native
	docker run -it --rm \
		--env GOOS=linux \
		--env GOARCH=$(GOARCH) \
		-v $(PWD).docker-build:/go/src/i2pgit.org/idk/i2p.plugins.tor-manager \
		-v $(GOPATH)/src/github.com/eyedeekay/go-I2P-jpackage:/go/src/github.com/eyedeekay/go-I2P-jpackage \
		eyedeekay/i2p.plugins.tor-manager
	cp -v $(PWD).docker-build/i2p.plugins.tor-manager* $(PWD)
	sudo chown $(USER):$(USER) $(PWD)/i2p.plugins.tor-manager*
		#-u `id -u $(USER)`:`id -g $(USER)` \
		#-e DISPLAY=unix$(DISPLAY) \
		#--publish 127.0.0.1:7695:7695 \
		#-v /tmp/.X11-unix:/tmp/.X11-unix \

#TBLANG?=en-US
TORRENT?=false

all-torrents:
	go build
	TOR_MANAGER_CLEARNET_MIRROR=true TOR_MANAGER_REQUIRE_PASSWORD=false TOR_MANAGER_NEVER_USE_TOR=true ./i2p.plugins.tor-manager -notor -mirrorall

distclean: clean

signer=70D2060738BEF80523ACAFF7D75C03B39B5E14E1

flatpak: clean-flatpak linplugin
	cp -v i2p.plugins.tor-manager-linux-amd64 i2p.plugins.tor-manager
	flatpak-builder --gpg-sign="$(signer)" --user --force-clean --disable-cache build-dir org.i2pgit.idk.i2p.plugins.tor-manager.yml
	flatpak-builder --gpg-sign="$(signer)" --user --install --force-clean build-dir org.i2pgit.idk.i2p.plugins.tor-manager.yml

clean-flatpak:
	rm -rf .flatpak-builder build-dir

flatpak-deps: shared-modules
	flatpak install flathub org.kde.Platform//5.15-21.08 org.kde.Sdk//5.15-21.08

shared-modules:
	git clone https://github.com/flathub/shared-modules.git && \
		cd shared-modules && \
		git checkout 8ce6437c269ef28c49984c11246d27be433c21d5

flatpak-repo: flatpak
	flatpak-builder --gpg-sign="$(signer)" --repo=repo --force-clean build-dir org.i2pgit.idk.i2p.plugins.tor-manager.yml

flatpak-add:
	flatpak --user remote-add --no-gpg-verify org.i2pgit.idk.i2p.plugins.tor-manager-dev repo; true

flatpak-remote-add:
	gpg --export D75C03B39B5E14E1 > org.i2pgit.idk.i2p.plugins.tor-manager/key.gpg
	flatpak --user remote-add --gpg-import=https://eyedeekay.github.io/flatpak.repo.i2p.plugins.tor-manager/key.gpg org.i2pgit.idk.i2p.plugins.tor-manager https://eyedeekay.github.io/flatpak.repo.i2p.plugins.tor-manager

flatpak-install: flatpak-repo flatpak-add
	flatpak --user install org.i2pgit.idk.i2p.plugins.tor-manager-dev org.i2pgit.idk.i2p.plugins.tor-manager

flatpak-remote-install:
	flatpak --user install org.i2pgit.idk.i2p.plugins.tor-manager org.i2pgit.idk.i2p.plugins.tor-manager

flatpak-update: flatpak-repo
	flatpak --user update org.i2pgit.idk.i2p.plugins.tor-manager

run-flatpak:
	flatpak run org.i2pgit.idk.i2p.plugins.tor-manager

flatpak.repo.i2p.plugins.tor-manager:
	git clone git@github.com:eyedeekay/flatpak.repo.i2p.plugins.tor-manager.git

flatpak-upload: flatpak-repo flatpak.repo.i2p.plugins.tor-manager
	cp -rv repo/* flatpak.repo.i2p.plugins.tor-manager/
	cd flatpak.repo.i2p.plugins.tor-manager && ./find.sh

clean-appimage:
	rm -rf AppDir

appimage: clean-appimage
	go build $(ARG)
	wget -c https://github.com/AppImage/AppImageKit/releases/download/continuous/runtime-x86_64
	mkdir -p AppDir/usr/bin AppDir/usr/lib AppDir/usr/share/applications AppDir/usr/share/icons AppDir/var/lib/i2pbrowser/icons
	cp -v i2p.plugins.tor-manager AppDir/usr/bin/i2p.plugins.tor-manager
	cp -v i2ptorbrowser.desktop AppDir/usr/share/applications/i2p.plugins.tor-manager.desktop
	cp -v i2ptorbrowser.desktop AppDir/i2p.plugins.tor-manager.desktop
	cp -v garliconion.png AppDir/var/lib/i2pbrowser/icons/garliconion.png
	cp -v garliconion.png AppDir/I2P_in-Tor-Browser-x86_64.png
	cp -v garliconion.png AppDir/.DirIcon
	cp -v i2p.plugins.tor-manager AppDir/AppRun
	find AppDir -name '*.desktop' -exec sed -i 's|garliconion.png|garliconion|g' {} \;
	chmod +x AppDir/AppRun
	~/Downloads/appimagetool-x86_64.AppImage --runtime-file runtime-x86_64 AppDir/

appimage-upload: appimage
	gothub upload -R -u eyedeekay -r $(BINARY) -t "$(VERSION)" -f "I2P_in_Tor_Browser-x86_64.AppImage" -n "I2P in Tor Browser x86_64 AppImage" -l "`sha256sum I2P in Tor Browser x86_64 AppImage`"

