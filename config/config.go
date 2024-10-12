package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
)

type ProxyConfig struct {
	Proxies []struct {
		Name        string `yaml:"name"`
		Source      string `yaml:"source"`
		Destination string `yaml:"destination"`
	} `yaml:"proxies"`
}

func LoadConfigFile(configFile string) (proxyConfig ProxyConfig, err error) {

	yamlData, err := openConfigFile(configFile)
	if err != nil {
		log.Fatal("Error reading YAML file:", err)
	}
	err = yaml.Unmarshal(yamlData, &proxyConfig)
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
