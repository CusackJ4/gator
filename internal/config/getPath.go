package config

import (
	"os"
	"path/filepath"
)

func GetPath() (string, error) {
	home, err := os.UserHomeDir() // Retrieves home directory
	if err != nil {
		return "", err
	}
	path := filepath.Join(home + "/.gatorconfig.json")

	return path, nil
}
