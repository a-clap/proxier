package config

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Config struct {
	Settings Settings
	Files    Files
}
type Settings struct {
	Settings map[string]interface{}
	values   map[string]string
}

type Files struct {
	Files []File
}

type File struct {
	Name   string   `json:"name"`
	Append []string `json:"append"`
	Remove []string `json:"remove"`
}

func New(buf []byte) (*Config, error) {
	var jsonMap map[string]json.RawMessage
	err := json.Unmarshal(buf, &jsonMap)
	if err != nil {
		return nil, err
	}

	c := &Config{
		Settings{
			Settings: make(map[string]interface{}),
			values:   make(map[string]string),
		},
		Files{Files: nil},
	}
	if err = json.Unmarshal(jsonMap["settings"], &c.Settings.Settings); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(jsonMap["files"], &c.Files); err != nil {
		return nil, err
	}

	c.getValues()
	c.parse()

	return c, nil
}

func (c *Config) getValues() {
	for key, value := range c.Settings.Settings {
		if value, ok := value.(string); ok {
			c.Settings.values[key] = value
		}
	}
}

func (c *Config) Get(key string) (string, error) {
	if v, ok := c.Settings.values[key]; ok {
		return v, nil
	} else {
		return "", fmt.Errorf("%s doesn't exist", key)
	}
}

func (c *Config) parseSingleVariable(value string) (newValue string) {
	newValue = value

	r := regexp.MustCompile(`\${(.+?[^}])}`)
	s := r.FindAllStringSubmatchIndex(value, len(value))
	if s == nil {
		return
	}
	var valueKey [][]string
	for _, idx := range s {
		key := value[idx[2]:idx[3]]
		fullKey := value[idx[0]:idx[1]]
		keyValue, err := c.Get(key)

		if err != nil {
			keyValue = ""
		}
		valueKey = append(valueKey, []string{keyValue, fullKey})
	}
	for _, pair := range valueKey {
		newValue = strings.Replace(newValue, pair[1], pair[0], -1)
	}

	return
}

func (c *Config) parse() {
	var s [][]string

	for key, v := range c.Settings.values {
		s = append(s, []string{key, c.parseSingleVariable(v)})
	}

	for _, value := range s {
		c.Settings.values[value[0]] = value[1]
	}
}
