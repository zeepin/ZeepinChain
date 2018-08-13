package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Wallets   []string `json:"Wallets"`
	Passwords []string `json:"Passwords"`
	RpcUrl    string   `json:"RpcUrl"`
	GIds      []string `json:"GIds"`
}

const (
	DefaultConfigFilename = "./config.json"
)

type ConfigFile struct {
	ConfigFile Config `json:"Configuration"`
}

var Configuration *Config

func Init() {
	file, e := ioutil.ReadFile(DefaultConfigFilename)
	if e != nil {
		panic(e)
	}

	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))

	config := ConfigFile{}
	e = json.Unmarshal(file, &config)
	if e != nil {
		panic(e)
	}
	Configuration = &(config.ConfigFile)
}
