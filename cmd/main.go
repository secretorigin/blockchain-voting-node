package main

import (
	"encoding/json"
	"io"
	"os"
)

func JsonParser(path string, config interface{}) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(byteValue, config)
}

func main() {

}
