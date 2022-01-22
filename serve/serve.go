package tbserve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	tbget "i2pgit.org/idk/i2p.plugins.tor-manager/get"
	TBSupervise "i2pgit.org/idk/i2p.plugins.tor-manager/supervise"
)

type Client struct {
	hostname string
	TBD      *tbget.TBDownloader
	TBS      *TBSupervise.Supervisor
}

func NewClient(hostname string, lang string, os string, arch string) (*Client, error) {
	m := &Client{
		hostname: hostname,
		TBD:      tbget.NewTBDownloader(lang, os, arch),
	}
	tgz, sig, err := m.TBD.DownloadUpdaterForLang(lang)
	if err != nil {
		panic(err)
	}
	var home string
	if home, err = m.TBD.CheckSignature(tgz, sig); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Signature check passed: %s %s", tgz, sig)
	}
	m.TBS = TBSupervise.NewSupervisor(home, lang)
	return m, nil
}

func (m *Client) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	path := rq.URL.Path
	switch path {
	case "/start-tor-browser":
		m.TBS.RunTBWithLang()
	case "/start-i2p-browser":
		m.TBS.RunI2PBWithLang()
	case "/start-tor":
		m.TBS.RunTorWithLang()
	case "/stop-tor":
		m.TBS.StopTor()
	default:
		b, e := m.Page()
		if e != nil {
			fmt.Fprintf(rw, "Error: %s", e)
		}
		rw.Write([]byte(b))
	}
	rw.Write([]byte("Hello, world!"))
}

func (m *Client) Serve() error {
	//http.Handle("/", m)
	return http.ListenAndServe("127.0.0.1:7695", m)
}

func (m *Client) generateMirrorJSON() (map[string]interface{}, error) {
	if !strings.HasSuffix(m.hostname, "/") {
		m.hostname += "/"
	}
	path := filepath.Join(tbget.DOWNLOAD_PATH, "downloads.json")
	preBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("GenerateMirrorJSON: %s", err)
	}
	binpath, _, err := m.TBD.GetUpdaterForLangFromJsonBytes(preBytes, "en-US")
	if err != nil {
		return nil, fmt.Errorf("GenerateMirrorJSON: %s", err)
	}
	urlparts := strings.Split(binpath, "/")
	replaceString := GenerateReplaceString(urlparts[:len(urlparts)-1])
	fmt.Printf("Replacing: %s with %s\n", replaceString, m.hostname)
	jsonBytes := []byte(strings.Replace(string(preBytes), replaceString, m.hostname, -1))
	var JSON map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &JSON); err != nil {
		panic(err)
	}
	return JSON, nil
}

func (m *Client) GenerateMirrorJSON(lang string) (string, error) {
	JSON, err := m.generateMirrorJSON()
	if err != nil {
		return "", err
	}
	path := filepath.Join(tbget.DOWNLOAD_PATH, "downloads.json")
	preBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("GenerateMirrorJSONBytes: %s", err)
	}
	binpath, _, err := m.TBD.GetUpdaterForLangFromJsonBytes(preBytes, "en-US")
	if err != nil {
		return "", fmt.Errorf("GenerateMirrorJSONBytes: %s", err)
	}
	urlparts := strings.Split(binpath, "/")
	replaceString := GenerateReplaceString(urlparts[:len(urlparts)-1])

	if platform, ok := JSON["downloads"]; ok {
		rtp := m.TBD.GetRuntimePair()
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
		return strings.Replace(string(bytes), replaceString, m.hostname, -1), nil
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
