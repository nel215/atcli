package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	ContestUrl string
}

type ConfigStore struct{}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{}
}

func (cs *ConfigStore) Save(config *Config) error {
	jb, err := json.Marshal(config)
	if err != nil {
		return err
	}
	f, err := os.Create("./.config.json")
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(jb)
	return nil
}

func (cs *ConfigStore) Load() (*Config, error) {
	f, err := os.Open("./.config.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	config := &Config{}
	err = json.Unmarshal(bs, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
