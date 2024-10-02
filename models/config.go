package models

import (
	ini "gopkg.in/ini.v1"
)

var config *ini.File


func Config(section string, key string) (string, error) {
	var err error
	if config == nil {
		config, err = ini.Load("config.ini")
		if err != nil {
			panic(err)
		}
	}

	sec, err := config.GetSection(section)
	if err != nil {
		return "", err
	}

	result, err := sec.GetKey(key)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
