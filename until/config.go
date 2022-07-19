package until

import (
	"encoding/json"
	"os"
)

func LoadConfigFromFile(path string, obj any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, obj)
	return err
}

func SaveConfigToFile(path string, obj any) error {

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
