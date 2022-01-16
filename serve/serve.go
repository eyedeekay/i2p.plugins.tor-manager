package tbserve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
)

func generateMirrorJSON(hostname string) (map[string]interface{}, error) {
	if !strings.HasSuffix(hostname, "/") {
		hostname += "/"
	}
	path := filepath.Join(tbget.DOWNLOAD_PATH, "downloads.json")
	preBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("GenerateMirrorJSON: %s", err)
	}
	binpath, _, err := tbget.GetUpdaterForLangFromJsonBytes(preBytes, "en-US")
	if err != nil {
		return nil, fmt.Errorf("GenerateMirrorJSON: %s", err)
	}
	urlparts := strings.Split(binpath, "/")
	replaceString := GenerateReplaceString(urlparts[:len(urlparts)-1])
	fmt.Printf("Replacing: %s with %s\n", replaceString, hostname)
	jsonBytes := []byte(strings.Replace(string(preBytes), replaceString, hostname, -1))
	var JSON map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &JSON); err != nil {
		panic(err)
	}
	return JSON, nil
}

func GenerateMirrorJSON(hostname, lang string) (string, error) {
	JSON, err := generateMirrorJSON(hostname)
	if err != nil {
		return "", err
	}
	path := filepath.Join(tbget.DOWNLOAD_PATH, "downloads.json")
	preBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("GenerateMirrorJSONBytes: %s", err)
	}
	binpath, _, err := tbget.GetUpdaterForLangFromJsonBytes(preBytes, "en-US")
	if err != nil {
		return "", fmt.Errorf("GenerateMirrorJSONBytes: %s", err)
	}
	urlparts := strings.Split(binpath, "/")
	replaceString := GenerateReplaceString(urlparts[:len(urlparts)-1])

	if platform, ok := JSON["downloads"]; ok {
		rtp := tbget.GetRuntimePair()
		for k, v := range platform.(map[string]interface{}) {
			if k != rtp {
				delete(platform.(map[string]interface{}), k)
			}
			for k2 := range v.(map[string]interface{}) {
				if k2 != lang {
					delete(v.(map[string]interface{}), k2)
				}

			}
		}
		bytes, err := json.MarshalIndent(JSON, "", "  ")
		if err != nil {
			return "", err
		}
		return strings.Replace(string(bytes), replaceString, hostname, -1), nil
	}
	return "", fmt.Errorf("GenerateMirrorJSONBytes: %s", "No downloads found")
}

func GenerateReplaceString(urlparts []string) string {
	replaceString := ""
	for _, val := range urlparts {
		if val == "https" {
			replaceString += val + "//"
		} else {
			replaceString += val + "/"
		}
	}
	if !strings.HasSuffix(replaceString, "/") {
		replaceString += "/"
	}
	return replaceString
}
