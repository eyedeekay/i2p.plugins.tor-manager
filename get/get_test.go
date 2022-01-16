package tbget

import (
	"testing"
)

func TestGet(t *testing.T) {
	binary, sig, err := GetUpdaterForLang("en-US")
	if err != nil {
		t.Error(err)
	}
	t.Log(binary, sig)
	binpath, sigpath, err := DownloadUpdater()
	if err != nil {
		t.Error(err)
	}
	t.Log(binpath, sigpath)
}
