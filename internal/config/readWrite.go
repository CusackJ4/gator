package config

import (
	"encoding/json"
	"os"
)

func Read() (Config, error) {
	// 1. Load empty Config struct
	cfg := Config{}

	// 2. Get filepath
	path, err := GetPath()

	// 3. Retrieve file using path
	jsonFile, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer jsonFile.Close()

	// 4. Create file decoder and decode file, storing in empty Config struct
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&cfg)

	return cfg, err
}

func writeJSONfile(c Config) error {

	path, err := GetPath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(c)

	return err
}
