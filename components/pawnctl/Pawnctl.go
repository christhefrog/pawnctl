package pawnctl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Ready bool `json:"-"`

	Compilers map[string]string `json:"compilers"`
}

var config Config

func LoadConfig() (*Config, error) {
	if config.Ready {
		return &config, nil
	}

	path := "pawnctl.json"

	path = fmt.Sprint(filepath.Dir(os.Args[0]), "\\", path)

	file, err := os.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(file, &config)
		if err != nil {
			return nil, err
		}
	}

	if config.Compilers == nil {
		config.Compilers = make(map[string]string)
	}

	config.Ready = true

	if err != nil {
		fmt.Println("Global config not found. Creating in", path)
		config.Save()
	}

	return &config, nil
}

func (c *Config) IsCompilerInstalled(name string) bool {
	return c.Compilers[name] != ""
}

func (c *Config) AddCompiler(name string, exec string) error {
	c.Compilers[name] = exec
	return nil
}

func (c *Config) GetLatestCompiler() (string, string) {
	latest := c.Compilers["latest"]
	return latest, c.Compilers[latest]
}

func (c *Config) ListCompilers() []string {
	arr := make([]string, 0)
	for k := range c.Compilers {
		arr = append(arr, k)
	}
	return arr
}

func (c *Config) Save() error {
	if !c.Ready {
		return errors.New("config handle isn't opened")
	}

	path := "pawnctl.json"
	path = fmt.Sprint(filepath.Dir(os.Args[0]), "\\", path)

	bytes, _ := json.MarshalIndent(c, "", "\t")

	os.WriteFile(path, bytes, 0664)

	return nil
}
