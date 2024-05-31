package sampman

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
	Servers   map[string]Server `json:"servers"`
}

type Server struct {
	Type     string `json:"type"` // samp or omp
	Path     string `json:"path"`
	Exec     string `json:"exec"`
	Includes string `json:"includes"`
}

var config Config

func LoadConfig() error {
	if config.Ready {
		return nil
	}

	path := "sampman.json"

	path = fmt.Sprint(filepath.Dir(os.Args[0]), "\\", path)

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Config not found. Creating in", path)

		Save()
	} else {
		err = json.Unmarshal(file, &config)
		if err != nil {
			return err
		}
	}

	if config.Compilers == nil {
		config.Compilers = make(map[string]string)
	}
	if config.Servers == nil {
		config.Servers = make(map[string]Server)
	}

	config.Ready = true

	return nil
}

func IsCompilerInstalled(name string) bool {
	return config.Compilers[name] != ""
}

func AddCompiler(name string, exec string) error {
	config.Compilers[name] = exec
	return nil
}

func GetLatestCompiler() (string, string) {
	latest := config.Compilers["latest"]
	return latest, config.Compilers[latest]
}

func AddServer(name string, server Server) error {
	config.Servers[name] = server
	return nil
}

func IsServerInstalled(name string) bool {
	return config.Servers[name].Type != ""
}

func Save() error {
	if !config.Ready {
		return errors.New("config handle isn't opened")
	}

	path := "sampman.json"
	path = fmt.Sprint(filepath.Dir(os.Args[0]), "\\", path)

	bytes, _ := json.MarshalIndent(config, "", "\t")

	os.WriteFile(path, bytes, 0664)

	return nil
}
