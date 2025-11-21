package storage

import (
	"encoding/json"
	"os"
)

func Load[T any](filename string) ([]T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var items []T
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func Save[T any](filename string, items []T) error {
	data, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)

}
