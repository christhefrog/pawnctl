package sampman

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Compilers map[string]string `json:"compilers"`
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

	if config.Compilers == nil {
		config.Compilers = make(map[string]string)
	}

	return config, nil
}

func (c *Config) SetCompiler(name string, exec string) error {
	c.Compilers[name] = exec
	return nil
}

func (c Config) Save(path string) error {
	path = fmt.Sprint(filepath.Dir(os.Args[0]), "\\", path)

	bytes, _ := json.MarshalIndent(c, "", "\t")

	os.WriteFile(path, bytes, 0664)

	return nil
}
