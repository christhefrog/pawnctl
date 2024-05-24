package sampman

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Compilers []CompilerInfo `json:"compilers"`
}

type CompilerInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Exec string `json:"exec"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	path = fmt.Sprint(filepath.Dir(os.Args[0]), "\\", path)

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Config not found. Creating in", path)
		config.Save("sampman.json")
	} else {
		err = json.Unmarshal(file, &config)
		if err != nil {
			return Config{}, err
		}
	}

	return config, nil
}

func (c CompilerInfo) IsInstalled() bool {
	return (CompilerInfo{}) != c
}

func (c Config) GetCompiler(name string) CompilerInfo {
	for _, v := range c.Compilers {
		if v.Name == name {
			return v
		}
	}
	return CompilerInfo{}
}

func (c *Config) AddCompiler(info CompilerInfo) error {
	c.Compilers = append(c.Compilers, info)
	return nil
}

func (c Config) Save(path string) error {
	path = fmt.Sprint(filepath.Dir(os.Args[0]), "\\", path)

	bytes, _ := json.MarshalIndent(c, "", "\t")

	os.WriteFile(path, bytes, 0664)

	return nil
}
