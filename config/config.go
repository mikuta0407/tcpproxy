package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

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
	if configFile != "" {
		log.Println("Using config file:", configFile)
		return os.ReadFile(configFile)
	}

	var err error
	configFile, err = getConfigFilePath()
	log.Println("Using config file:", configFile)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(configFile)
}

// getConfigFilePath は、実行されているOSに応じて設定ファイルのパスを返します。
//
//   - Linux: XDG Base Directory Specification に従ってパスを検索します。
//     (例: $XDG_CONFIG_HOME/tcpproxy/tcpproxy.yml または ~/.config/tcpproxy/tcpproxy.yml)
//   - macOS: ~/.config/tcpproxy/tcpproxy.yml
//   - Windows: %userprofile%/tcpproxy/tcpproxy.yml (例: C:\Users\YourUser\tcpproxy\tcpproxy.yml)
func getConfigFilePath() (string, error) {
	const subpath = "tcpproxy/tcpproxy.yml"

	switch runtime.GOOS {
	case "linux":
		// xdg.SearchConfigFile は XDG_CONFIG_DIRS と XDG_CONFIG_HOME を
		// チェックしてファイルを探します。見つからない場合はエラーを返します。
		// 一般的にはまず $XDG_CONFIG_HOME/tcpproxy/tcpproxy.yml を探し、
		// 見つからなければ /etc/xdg/tcpproxy/tcpproxy.yml などを探します。
		// ここでは最初に見つかった設定ファイルのパスを返します。
		return xdg.SearchConfigFile(subpath)

	case "darwin": // macOS
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("ユーザーのホームディレクトリを取得できませんでした: %w", err)
		}
		return filepath.Join(homeDir, ".config", "tcpproxy", "tcpproxy.yml"), nil

	case "windows":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("ユーザーのホームディレクトリ (%userprofile%) を取得できませんでした: %w", err)
		}
		// Windows の場合は直接ホームディレクトリ配下に置きます。
		return filepath.Join(homeDir, "tcpproxy", "tcpproxy.yml"), nil

	default:
		return "", fmt.Errorf("サポートされていないOSです: %s", runtime.GOOS)
	}
}
