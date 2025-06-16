package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func VerifyConfigurationExistence() {
	home, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	configDir := filepath.Join(home, ".config", "oko")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("OS ERROR: failed to create a directory")
		panic(err)
	}

	configFile := filepath.Join(configDir, CONFIG_FILE_NAME)

	//checking if configuration exists
	if !isFileExists(configFile) {
		_, err = os.Create(configFile)

		if err != nil {
			fmt.Println("OS ERROR: failed to create configuration file")
			panic(err)
		}

		config := Configuration{
			Folders: []string{},
		}
		data, _ := json.MarshalIndent(config, "", " ")
		if err := os.WriteFile(configFile, data, 0644); err != nil {
			panic("OS ERROR: failed to write data in the configuration file")
		}
	} else {
		_, err = GetConfiguration()
		if err != nil {
			panic(err)
		}
	}

}

func GetConfiguration() (*Configuration, error) {
	var err error = nil
	once.Do(func() {
		var data []byte
		data, err = os.ReadFile(GetConfigurationFilePath())
		if err != nil {
			err = fmt.Errorf("failed to read config file\n%s", err.Error())
			return
		}
		var c Configuration
		err = json.Unmarshal(data, &c)
		if err != nil {
			err = fmt.Errorf("incorrect formatting of configuration data\n%s", err.Error())
			return
		}
		config = &c
	})
	return config, err
}

func isFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
