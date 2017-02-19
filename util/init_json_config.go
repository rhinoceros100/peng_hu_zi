package util

import (
	"io/ioutil"
	"encoding/json"
)

func InitJsonConfig(file string, config interface{}) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return err
	}

	return nil
}
