package main

import (
	"encoding/json"
	"os"
)

func exportToJson() error {
	if len(golds) == 0 {
		return nil
	}

	file, err := os.Create("result.json")
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(golds)
	if err != nil {
		return err
	}

	return nil
}
