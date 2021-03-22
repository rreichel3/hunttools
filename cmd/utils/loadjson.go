package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func LoadJsonList(sourceFilePath string) ([]map[string]interface{}, error) {
	jsonFile, err := os.Open(sourceFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil

}

func LoadJsonMap(sourceFilePath string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(sourceFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil

}
