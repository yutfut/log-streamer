package tools

import "github.com/BurntSushi/toml"

type Config struct {
	Files   []string `toml:"files"`
}

func NewConfig() *Config {
	return &Config{}
}

func ReadConfigFile(configPath string, dst interface{}, ) error {
	_, err := toml.DecodeFile(configPath, dst)
	return err
}