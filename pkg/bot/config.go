package bot

import (
	"encoding/json"
	"os"
)

type config struct {
	Scope       string `json:"scope"`
	RedirectURI string `json:"redirectURI"`
	App         []struct {
		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
		State        string `json:"state"`
	} `json:"app"`
	APIInvokeInterval int    `json:"apiInvokeInterval"`
	DataFileName      string `json:"dataFileName"`

	// http server
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
}

func loadConfig() (*config, error) {
	bs, err := os.ReadFile("./appsetting.json")
	if err != nil {
		return nil, err
	}

	var jo config
	if err = json.Unmarshal(bs, &jo); err != nil {
		return nil, err
	}

	return &jo, nil
}
