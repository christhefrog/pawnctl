package sampman

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Compilers []string `json:"compilers"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	path = fmt.Sprint(filepath.Dir(os.Args[0]), "\\", path)

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Config not found. Creating in", path)
		bytes, _ := json.MarshalIndent(config, "", "\t")
		os.WriteFile(path, bytes, 0664)
	} else {
		err = json.Unmarshal(file, &config)
		if err != nil {
			return Config{}, err
		}
	}

	return config, nil
}
