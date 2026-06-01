package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}
	return filepath.Join(path, configFileName), nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("could not get config file path: %w", err)
	}
	config := Config{}
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("could not open config file: %w", err)
	}
	defer log.Fatal(file.Close())
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("could not parse config file: %w", err)
	}
	return config, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	return write(*c)
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("could not create config file: %w", err)
	}
	defer log.Fatal(file.Close())

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
