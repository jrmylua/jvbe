package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Oauth struct {
	Domain       string `yaml:"domain"`
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	CallbackUrl  string `yaml:"callback_url"`
    LogoutRedirectUrl string `yaml:"logout_redirect_url"`
}

type Config struct {
	DbConn string `yaml:"db_conn"`
	Port   int    `yaml:"port"`
	Oauth  Oauth  `yaml:"oauth"`
}

func ReadFile(src string) (*Config, error) {
	b, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
