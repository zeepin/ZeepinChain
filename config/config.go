/*
 * Copyright (C) 2018 The ZeepinChain Authors
 * This file is part of The ZeepinChain library.
 *
 * The ZeepinChain is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ZeepinChain is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ZeepinChain.  If not, see <http://www.gnu.org/licenses/>.
 */

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
