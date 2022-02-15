module i2pgit.org/idk/i2p.plugins.tor-manager

go 1.17

require (
	github.com/cloudfoundry/jibber_jabber v0.0.0-20151120183258-bcc4c8345a21
	github.com/eyedeekay/go-i2pd v0.0.0-20220213070306-9807541b2dfc
	github.com/getlantern/systray v1.1.0
	github.com/go-ole/go-ole v1.2.6
	github.com/jchavannes/go-pgp v0.0.0-20200131171414-e5978e6d02b4
	github.com/magisterquis/connectproxy v0.0.0-20200725203833-3582e84f0c9b
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
)

require (
	github.com/boreq/friendlyhash v0.0.0-20190522010448-1ca64b3ca69e // indirect
	github.com/eyedeekay/go-i2cp v0.0.0-20190716135428-6d41bed718b0 // indirect
	github.com/eyedeekay/goSam v0.32.31-0.20210415231611-c6d9c0e340b8 // indirect
	github.com/eyedeekay/sam-forwarder v0.0.0-20190908210105-71ca8cd65fda // indirect
	github.com/getlantern/context v0.0.0-20190109183933-c447772a6520 // indirect
	github.com/getlantern/errors v1.0.1 // indirect
	github.com/getlantern/golog v0.0.0-20201105130739-9586b8bde3a9 // indirect
	github.com/getlantern/hex v0.0.0-20190417191902-c6586a6fe0b7 // indirect
	github.com/getlantern/hidden v0.0.0-20190325191715-f02dbb02be55 // indirect
	github.com/getlantern/ops v0.0.0-20200403153110-8476b16edcd6 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69 // indirect
	github.com/mwitkow/go-http-dialer v0.0.0-20161116154839-378f744fb2b8 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2 // indirect
	github.com/ybbus/jsonrpc/v2 v2.1.6 // indirect
	github.com/zieckey/goini v0.0.0-20180118150432-0da17d361d26 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
)

require (
	github.com/eyedeekay/checki2cp v0.0.21
	github.com/eyedeekay/go-i2pcontrol v0.0.0-20201227222421-6e9f31a29a91
	github.com/eyedeekay/httptunnel v0.0.0-20210508193128-6e9606d6eb24
	github.com/eyedeekay/sam3 v0.32.32
	github.com/itchio/damage v0.0.0-20190703135837-76df725fc766
	github.com/itchio/headway v0.0.0-20200301160421-e15721f23905
	github.com/justinas/nosurf v1.1.1
	github.com/mitchellh/go-ps v1.0.0
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.17.0 // indirect
	github.com/otiai10/copy v1.7.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/russross/blackfriday v1.6.0
	github.com/ulikunitz/xz v0.5.10
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	howett.net/plist v1.0.0 // indirect
)

replace github.com/eyedeekay/go-i2pd v0.0.0-20220213070306-9807541b2dfc => ./go-i2pd

replace github.com/getlantern/systray v1.1.0 => github.com/eyedeekay/systray v1.1.1-0.20220213191004-800d7458fdfd
