package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadMovies(filePath string, movies interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %w", err)
	}

	if err := json.Unmarshal(data, &movies); err != nil {
		return fmt.Errorf("error deserializing JSON: %w", err)
	}

	return nil
}