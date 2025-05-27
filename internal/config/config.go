package config

import (
	"fmt"
	"os"

	"encoding/json"
)

type Config struct {
	DbURL    string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

func (c *Config) SetUser(name string) error {
	c.UserName = name
	jsonData, err := json.Marshal(c)

	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%s/.gatorconfig.json", home), jsonData, 644)
	if err != nil {
		return err
	}

	return nil

}

func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(fmt.Sprintf("%s/.gatorconfig.json", home))

	var conf Config

	json.Unmarshal(data, &conf)

	return conf, nil
}
