package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

const CONFIG_FILE_NAME = "config.json"

type Configuration struct {
	Folders   []string `json:"folders"`
	YTDlpPath string   `json:"yt-dlp-path"`
	ApiKey    string   `json:"yt-api-key"`
}

var (
	config *Configuration
	once   sync.Once
)

func GetConfigFolders() []string {
	conf, err := GetConfiguration()
	if err != nil {
		panic(err)
	}
	return conf.Folders
}

func AddConfigFolder(path string) {
	conf, err := GetConfiguration()
	if err != nil {
		panic(err)
	}
	conf.Folders = append(conf.Folders, path)
	data, err := json.MarshalIndent(conf, "", " ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(GetConfigurationFilePath(), data, 0644)
	if err != nil {
		panic(err)
	}
}

func GetConfigurationFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".config", "oko", CONFIG_FILE_NAME)
}
