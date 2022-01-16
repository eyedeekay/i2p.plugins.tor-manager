package tbserve

import (
	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
	"testing"
)

func TestServe(t *testing.T) {
	tbget.DOWNLOAD_PATH = "../tor-browser"
	bytes, err := GenerateMirrorJSON("http://localhost:8080", "en-US")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bytes))
}
