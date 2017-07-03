package config

import (
	"os"
)

// Key -
type Key string

var config = make(map[Key]string)

// Get -
func Get(key Key) string {
	rs, ok := config[key]
	if !ok {
		rs = os.Getenv(string(key))
		if rs != "" {
			config[key] = rs
		}
	}
	return rs
}

func Set(key Key, value string) error {
	err := os.Setenv(string(key), value)
	if err != nil {
		return err
	}

	config[key] = value
	return nil
}
