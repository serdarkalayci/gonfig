package gonfig

import (
	"encoding/json"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var myReadFile = os.ReadFile

func readJSON(filePath string) (map[string]interface{}, error) {
	readBytes, err := myReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var output map[string]interface{}
	err = json.Unmarshal(readBytes, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func readYaml(filePath string) (map[string]interface{}, error) {
	readBytes, err := myReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var output map[string]interface{}
	err = yaml.Unmarshal(readBytes, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
