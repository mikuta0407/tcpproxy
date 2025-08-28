package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
)

// Proxy は単一のポートフォワーディングルールを定義します。
type Proxy struct {
	Name        string `yaml:"name"`
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

// Config は設定ファイル全体を定義します。
type Config struct {
	Proxies []Proxy `yaml:"proxies"`
}

func LoadConfigFile(configFile string) (proxyConfig *Config, err error) {

	yamlData, err := openConfigFile(configFile)
	if err != nil {
		log.Fatal("Error reading YAML file:", err)
	}
	err = yaml.Unmarshal(yamlData, proxyConfig)
	if err != nil {
		log.Fatal("Error unmarshaling YAML:", err)
	}
	return
}

func openConfigFile(configFile string) ([]byte, error) {
	var err error
	if configFile != "" {
		yamlData, err := os.ReadFile(configFile)
		if err == nil {
			return yamlData, err
		}
	}

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	cfp, err := xdg.SearchConfigFile("tcpproxy/tcpproxy.yml")
	if err != nil {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		yamlData, err := os.ReadFile(filepath.Join(home, "tcpproxy.yml"))
		return yamlData, err
	}

	return os.ReadFile(cfp)
}
