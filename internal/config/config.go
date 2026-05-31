package config

import (
	"encoding/json"
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
		log.Fatal(err)
		return "", err
	}
	return filepath.Join(path, configFileName), nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
		return Config{}, err
	}
	config := Config{}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return Config{}, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatal(err)
		return Config{}, err
	}
	log.Print(config)
	return config, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	return write(*c)
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
