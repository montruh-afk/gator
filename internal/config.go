package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Url       string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func getConfigPath() (string, error) {
	//get home dir
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	//get the full path
	configPath := filepath.Join(home, ".gatorconfig.json")
	return configPath, nil
}

func Read() (Config, error) {
	file := Config{}
	

	configPath, err := getConfigPath()
	if err != nil {
		return file, err
	}

	// Read from path after error checking
	data, err := os.ReadFile(configPath)
	if err != nil {
		return file, err
	}

	//Unmarshal data into Config struct
	if err := json.Unmarshal(data, &file); err != nil {
		return file, err
	}

	return file, nil
}

func Write(cfg Config) error {
	filePath, err := getConfigPath()

	if err != nil {
		return err
	}
	//Prepare data to be written to the configuration file
	save, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	//Write to config located a filePath with permission "0666" which provides read/write permissions to everyone
	if err := os.WriteFile(filePath, save, 0666); err != nil {
		return err
	}
	return nil
}

func (c *Config) SetUser(user string) error {
	c.Current_user_name = user
	return Write(*c)
}