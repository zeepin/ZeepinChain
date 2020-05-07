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
package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/imZhuFei/zeepin/cmd/abi"
)

func NewEmbedContractAbi(abiData []byte) (*abi.EmbedContractAbi, error) {
	abi := &abi.EmbedContractAbi{}
	err := json.Unmarshal(abiData, abi)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal EmbedContractAbi error:%s", err)
	}
	return abi, nil
}

func ParseEmbededFunc(rawParams []string, funcAbi *abi.EmbededContractFunctionAbi) ([]interface{}, error) {
	res := make([]interface{}, 0)
	funcName := convertEmbededFuncName(funcAbi.Name)
	res = append(res, funcName)
	params, err := ParseEmbededParam(rawParams, funcAbi.Parameters)
	if err != nil {
		return nil, err
	}
	res = append(res, params)
	return res, nil
}

//Embeded func name in Camel-Case. For example: transfer, transferFrom
func convertEmbededFuncName(name string) string {
	if name == "" {
		return name
	}
	data := []byte(name)
	data[0] = strings.ToLower(string(data[0]))[0]
	return string(data)
}

func ParseEmbededParam(params []string, paramsAbi []*abi.EmbededContractParamsAbi) ([]interface{}, error) {
	if len(params) != len(paramsAbi) {
		return nil, fmt.Errorf("abi param not match")
	}
	val := make([]interface{}, 0)
	for i, rawParam := range params {
		paramAbi := paramsAbi[i]
		rawParam = strings.TrimSpace(rawParam)
		var res interface{}
		var err error
		switch strings.ToLower(paramAbi.Type) {
		case abi.EMBED_PARAM_TYPE_INTEGER:
			res, err = ParseEmbededParamInteger(rawParam)
		case abi.EMBED_PARAM_TYPE_BOOL:
			res, err = ParseEmbededParamBoolean(rawParam)
		case abi.EMBED_PARAM_TYPE_STRING:
			res, err = ParseEmbededParamString(rawParam)
		case abi.EMBED_PARAM_TYPE_BYTE_ARRAY:
			res, err = ParseEmbededParamByteArray(rawParam)
		default:
			return nil, fmt.Errorf("unknown param type:%s", paramAbi.Type)
		}
		if err != nil {
			return nil, fmt.Errorf("parse param:%s value:%s type:%s error:%s", paramAbi.Name, rawParam, paramAbi.Type, err)
		}
		val = append(val, res)
	}
	return val, nil
}

func ParseEmbededParamString(param string) (interface{}, error) {
	return param, nil
}

func ParseEmbededParamInteger(param string) (interface{}, error) {
	if param == "" {
		return nil, fmt.Errorf("invalid integer")
	}
	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse integer param:%s error:%s", param, err)
	}
	return value, nil
}

func ParseEmbededParamBoolean(param string) (interface{}, error) {
	var res bool
	switch strings.ToLower(param) {
	case "true":
		res = true
	case "false":
		res = false
	default:
		return nil, fmt.Errorf("parse boolean param:%s failed", param)
	}
	return res, nil
}

func ParseEmbededParamByteArray(param string) (interface{}, error) {
	res, err := hex.DecodeString(param)
	if err != nil {
		return nil, fmt.Errorf("parse byte array param:%s error:%s", param, err)
	}
	return res, nil
}
