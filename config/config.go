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

func GetYoutubeApiKey() string {
	conf, err := GetConfiguration()
	if err != nil {
		panic(err)
	}
	return conf.ApiKey
}

func GetYTDlpPath() string {
	conf, err := GetConfiguration()
	if err != nil {
		panic(err)
	}

	if conf.YTDlpPath == "" {
		return "yt-dlp"
	}

	return conf.YTDlpPath
}

func GetCacheSongsDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(cacheDir, "oko")
}

func GetSongFolders() []string {
	configFolders := GetConfigFolders()
	folders := make([]string, 0, len(configFolders)+1)
	seen := make(map[string]struct{}, len(configFolders)+1)

	for _, folder := range configFolders {
		if _, exists := seen[folder]; exists {
			continue
		}

		seen[folder] = struct{}{}
		folders = append(folders, folder)
	}

	cacheSongsDir := GetCacheSongsDir()
	if _, exists := seen[cacheSongsDir]; !exists {
		folders = append(folders, cacheSongsDir)
	}

	return folders
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
