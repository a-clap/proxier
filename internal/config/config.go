package config

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type Config struct {
	Settings  Settings
	Variables Variables
	Files     Files `json:"files"`
}
type Settings struct {
	Settings map[string]interface{}
}

type Variables struct {
	Variables map[string]interface{}
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
		},
		Variables{
			Variables: make(map[string]interface{}),
		},
		Files{Files: nil},
	}
	if err = json.Unmarshal(jsonMap["settings"], &c.Settings.Settings); err != nil {
		//	log sth
	}

	if err = json.Unmarshal(jsonMap["variables"], &c.Variables.Variables); err != nil {

	} else {
		c.parseVariables()
	}
	return c, nil
}

func (c *Config) Get(key string) (string, error) {
	if v, ok := c.Settings.Settings[key]; ok {
		if s, isString := v.(string); isString {
			return s, nil
		}
	}

	if v, ok := c.Variables.Variables[key]; ok {
		if s, isString := v.(string); isString {
			return s, nil
		}
	}
	return "", fmt.Errorf("%s doesn't exist", key)
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

func (c *Config) parseVariables() {
	var s [][]string
	for key, variable := range c.Variables.Variables {
		if v, ok := variable.(string); ok {
			s = append(s, []string{key, c.parseSingleVariable(v)})
		}
	}
	for _, value := range s {
		c.Variables.Variables[value[0]] = value[1]
	}
}
