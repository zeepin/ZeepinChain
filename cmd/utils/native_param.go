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

 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/imZhuFei/zeepin/cmd/abi"
	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/core/types"
	"github.com/imZhuFei/zeepin/embed/simulator"
	httpcom "github.com/imZhuFei/zeepin/http/base/common"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/embed"
)

func NewNativeInvokeTransaction(gasPrice, gasLimit uint64, contractAddr common.Address, version byte, params []interface{}, funcAbi *abi.NativeContractFunctionAbi) (*types.MutableTransaction, error) {
	builder := simulator.NewParamsBuilder(new(bytes.Buffer))
	err := ParseNativeFuncParam(builder, funcAbi.Name, params, funcAbi.Parameters)
	if err != nil {
		return nil, err
	}
	builder.EmitPushByteArray([]byte(funcAbi.Name))
	builder.EmitPushByteArray(contractAddr[:])
	builder.EmitPushInteger(new(big.Int).SetInt64(int64(version)))
	builder.Emit(simulator.SYSCALL)
	builder.EmitPushByteArray([]byte(embed.NATIVE_INVOKE_NAME))
	invokeCode := builder.ToArray()
	return httpcom.NewSmartContractTransaction(gasPrice, gasLimit, invokeCode, 0)
}

func ParseNativeFuncParam(builder *simulator.ParamsBuilder, funName string, params []interface{}, paramsAbi []*abi.NativeContractParamAbi) error {
	size := len(paramsAbi)
	if size == 0 {
		//Params cannot empty, if params is empty, fulfil with func name
		params = []interface{}{funName}
		paramsAbi = []*abi.NativeContractParamAbi{{
			Name: "funcName",
			Type: abi.NATIVE_PARAM_TYPE_STRING,
		}}
	} else if size > 1 {
		//If more than one param in func, must using struct
		paramRoot := &abi.NativeContractParamAbi{
			Name:    "root",
			Type:    abi.NATIVE_PARAM_TYPE_STRUCT,
			SubType: paramsAbi,
		}
		params = []interface{}{params}
		paramsAbi = []*abi.NativeContractParamAbi{paramRoot}
	}
	return ParseNativeParams(builder, params, paramsAbi)
}

func ParseNativeParams(builder *simulator.ParamsBuilder, params []interface{}, paramsAbi []*abi.NativeContractParamAbi) error {
	if len(params) != len(paramsAbi) {
		return fmt.Errorf("abi unmatch")
	}
	var err error
	for i, param := range params {
		paramAbi := paramsAbi[i]
		switch strings.ToLower(paramAbi.Type) {
		case abi.NATIVE_PARAM_TYPE_STRUCT:
			err = ParseNativeParamStruct(builder, param, paramAbi)
			if err != nil {
				return fmt.Errorf("param:%s parse:%v error:%s", paramAbi.Name, param, err)
			}
		case abi.NATIVE_PARAM_TYPE_ARRAY:
			err = ParseNativeParamArray(builder, param, paramAbi)
		default:
			rawParam, ok := param.(string)
			if !ok {
				return fmt.Errorf("param:%v assert to string failed", param)
			}
			rawParam = strings.TrimSpace(rawParam)
			switch strings.ToLower(paramAbi.Type) {
			case abi.NATIVE_PARAM_TYPE_ADDRESS:
				err = ParseNativeParamAddress(builder, rawParam)
			case abi.NATIVE_PARAM_TYPE_BOOL:
				err = ParseNativeParamBool(builder, rawParam)
			case abi.NATIVE_PARAM_TYPE_BYTE:
				err = ParseNativeParamByte(builder, rawParam)
			case abi.NATIVE_PARAM_TYPE_BYTEARRAY:
				err = ParseNativeParamByteArray(builder, rawParam)
			case abi.NATIVE_PARAM_TYPE_INTEGER:
				err = ParseNativeParamInteger(builder, rawParam)
			case abi.NATIVE_PARAM_TYPE_STRING:
				err = ParseNativeParamString(builder, rawParam)
			case abi.NATIVE_PARAM_TYPE_UINT256:
				err = ParseNativeParamUint256(builder, rawParam)
			default:
				return fmt.Errorf("unknown param type:%s", paramAbi.Type)
			}
		}
		if err != nil {
			return fmt.Errorf("param:%s parse:%v error:%s", paramAbi.Name, param, err)
		}
	}

	return nil
}

func ParseNativeParamStruct(builder *simulator.ParamsBuilder, param interface{}, structAbi *abi.NativeContractParamAbi) error {
	params, ok := param.([]interface{})
	if !ok {
		return fmt.Errorf("assert to []interface{} failed")
	}
	if len(params) != len(structAbi.SubType) {
		return fmt.Errorf("struct abi not match")
	}
	builder.EmitPushInteger(big.NewInt(0))
	builder.Emit(simulator.NEWSTRUCT)
	builder.Emit(simulator.TOALTSTACK)
	for i, param := range params {
		paramAbi := structAbi.SubType[i]
		err := ParseNativeParams(builder, []interface{}{param}, []*abi.NativeContractParamAbi{paramAbi})
		if err != nil {
			return fmt.Errorf("params struct:%s item:%s error:%s", structAbi.Name, paramAbi.Name, err)
		}
		builder.Emit(simulator.DUPFROMALTSTACK)
		builder.Emit(simulator.SWAP)
		builder.Emit(simulator.APPEND)
	}
	builder.Emit(simulator.FROMALTSTACK)
	return nil
}

func ParseNativeParamArray(builder *simulator.ParamsBuilder, param interface{}, arrayAbi *abi.NativeContractParamAbi) error {
	params, ok := param.([]interface{})
	if !ok {
		return fmt.Errorf("assert to []interface{} failed")
	}
	abis := make([]*abi.NativeContractParamAbi, 0, len(params))
	for i := 0; i < len(params); i++ {
		abis = append(abis, &abi.NativeContractParamAbi{
			Name:    fmt.Sprintf("%s_%d", arrayAbi.Name, i),
			Type:    arrayAbi.SubType[0].Type,
			SubType: arrayAbi.SubType[0].SubType,
		})
	}
	err := ParseNativeParams(builder, params, abis)
	if err != nil {
		return fmt.Errorf("parse array error:%s", err)
	}
	builder.EmitPushInteger(big.NewInt(int64(len(params))))
	builder.Emit(simulator.PACK)
	return nil
}

func ParseNativeParamByte(builder *simulator.ParamsBuilder, param string) error {
	if param == "" {
		return fmt.Errorf("invalid byte")
	}
	i, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		return fmt.Errorf("parse int error:%s", err)
	}
	if i > math.MaxUint8 || i < 0 {
		return fmt.Errorf("invalid byte value")
	}
	builder.EmitPushInteger(new(big.Int).SetInt64(i))
	return nil
}

func ParseNativeParamByteArray(builder *simulator.ParamsBuilder, param string) error {
	data, err := hex.DecodeString(param)
	if err != nil {
		return fmt.Errorf("hex decode string error:%s", err)
	}
	builder.EmitPushByteArray(data)
	return nil
}

func ParseNativeParamUint256(builder *simulator.ParamsBuilder, param string) error {
	if param == "" {
		return fmt.Errorf("invalid uint256")
	}
	uint256, err := common.Uint256FromHexString(param)
	if err != nil {
		return fmt.Errorf("invalid uint256")
	}
	builder.EmitPushByteArray(uint256.ToArray())
	return nil
}

func ParseNativeParamString(builder *simulator.ParamsBuilder, param string) error {
	builder.EmitPushByteArray([]byte(param))
	return nil
}

func ParseNativeParamInteger(builder *simulator.ParamsBuilder, param string) error {
	if param == "" {
		return fmt.Errorf("invalid integer")
	}
	i, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return fmt.Errorf("parse int error:%s", err)
	}
	builder.EmitPushInteger(new(big.Int).SetInt64(i))
	return nil
}

func ParseNativeParamBool(builder *simulator.ParamsBuilder, param string) error {
	var b bool
	switch strings.ToLower(param) {
	case "true":
		b = true
	case "false":
		b = false
	default:
		return fmt.Errorf("invalid bool value")
	}
	builder.EmitPushBool(b)
	return nil
}

func ParseNativeParamAddress(builder *simulator.ParamsBuilder, param string) error {
	if param == "" {
		return fmt.Errorf("invalid address")
	}
	var addr common.Address
	var err error
	//Maybe param is a contract address
	addr, err = common.AddressFromHexString(param)
	if err != nil {
		//Maybe param is a account address
		addr, err = common.AddressFromBase58(param)
		if err != nil {
			return fmt.Errorf("invalid address")
		}
	}

	builder.EmitPushByteArray(addr[:])
	return nil
}
