package config

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Settings Settings
}
type Settings struct {
	User      *string `json:"user"`
	Password  *string `json:"password"`
	HttpProxy *string `json:"http_proxy_server"`
	Port      *string `json:"port"`
}

type variables struct {
	variables map[string]interface{}
}

func New(buf []byte) (*Config, error) {
	c := &Config{}
	err := json.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}
	if err = c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) validate() error {
	if c.Settings.User == nil {
		return fmt.Errorf("field \"user\" not found")
	}
	if c.Settings.Password == nil {
		return fmt.Errorf("field \"password\" not found")
	}
	if c.Settings.HttpProxy == nil {
		return fmt.Errorf("field \"http_proxy_server\" not found")
	}
	if c.Settings.Port == nil {
		return fmt.Errorf("field \"port\" not found")
	}

	return nil
}
