package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ActionMap map[string]Action

func (a *ActionMap) UnmarshalJSON(data []byte) error {
	actions := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &actions); err != nil {
		return err
	}
	avActs := GetAvailableActions()
	result := make(ActionMap)
	for k, v := range actions {
		if creator, ok := avActs[k]; ok {
			actionMap := make(map[string]json.RawMessage)
			if err := json.Unmarshal(v, &actionMap); err != nil {
				return err
			}
			for name, value := range actionMap {
				action := creator()
				if err := json.Unmarshal(value, &action); err != nil {
					return err
				}
				result[name] = action
			}
		}
	}
	*a = result
	return nil
}

type Config struct {
	Actions ActionMap `json:"action"`
}

func GetConfig(file string) (*Config, error) {
	// TODO: parse all files within directory and compose them (check ending .json)
	return GetConfigFromJsonFile(file)
}

func GetConfigFromJsonFile(file string) (*Config, error) {
	byteContent, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err // TODO: wrap this error
	}
	c := &Config{}
	if err := json.Unmarshal(byteContent, &c); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return c, nil
}
