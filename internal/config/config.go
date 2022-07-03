package config

import (
	"encoding/json"
	"fmt"
	"log"
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

func Template() []byte {
	type template struct {
		Settings map[string]string `json:"settings"`
		Files    []File            `json:"files"`
	}

	s := template{
		Settings: make(map[string]string),
		Files:    make([]File, 2),
	}
	s.Settings["user"] = "user"
	s.Settings["password"] = "password"
	s.Settings["server"] = "192.168.0.100"
	s.Settings["port"] = "80"
	s.Settings["http_proxy"] = "\"http://${user}:${password}@${server}:${port}\""
	s.Settings["https_proxy"] = "\"https://${user}:${password}@${server}:${port}\""

	s.Files[0] = File{
		Name:   "/etc/environment",
		Append: []string{"HTTP_PROXY=${http_proxy}"},
		Remove: []string{"HTTP_PROXY"},
	}

	s.Files[1] = File{
		Name: "/etc/apt/apt.conf.d/proxy.conf",
		Append: []string{
			"Acquire::http::proxy ${http_proxy}",
			"Acquire::https::proxy ${http_proxy}",
		},
		Remove: []string{"Acquire"},
	}

	buf, err := json.Marshal(s)
	if err != nil {
		log.Fatalln("error parsing template config, shouldn't happen")
	}
	return buf
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
	if err = json.Unmarshal(jsonMap["files"], &c.Files.Files); err != nil {
		return nil, err
	}

	c.getValues()
	c.parse()

	return c, nil
}

func (c *Config) Get(key string) (string, error) {
	if v, ok := c.Settings.values[key]; ok {
		return v, nil
	} else {
		return "", fmt.Errorf("%s doesn't exist", key)
	}
}

func (c *Config) GetFiles() []File {
	var files []File = nil

	for _, f := range c.Files.Files {
		singleFile := File{
			Name: f.Name,
		}
		for _, line := range f.Append {
			singleFile.Append = append(singleFile.Append, c.parseSingleVariable(line))
		}
		for _, line := range f.Remove {
			singleFile.Remove = append(singleFile.Remove, c.parseSingleVariable(line))
		}
		files = append(files, singleFile)
	}

	return files
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

func (c *Config) getValues() {
	for key, value := range c.Settings.Settings {
		if value, ok := value.(string); ok {
			c.Settings.values[key] = value
		}
	}
}
