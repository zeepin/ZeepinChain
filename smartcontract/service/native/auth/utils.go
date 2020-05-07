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

package auth

import (
	"bytes"
	"fmt"
	"io"

	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/common/serialization"
	"github.com/imZhuFei/zeepin/smartcontract/event"
	"github.com/imZhuFei/zeepin/smartcontract/service/native"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/utils"
)

var (
	PreAdmin          = []byte{0x01}
	PreRoleFunc       = []byte{0x02}
	PreRoleToken      = []byte{0x03}
	PreDelegateStatus = []byte{0x04}
)

//type(this.contractAddr.Admin) = []byte
func concatContractAdminKey(native *native.NativeService, contractAddr common.Address) []byte {
	this := native.ContextRef.CurrentContext().ContractAddress
	adminKey := append(this[:], contractAddr[:]...)
	adminKey = append(adminKey, PreAdmin...)

	return adminKey
}

func getContractAdmin(native *native.NativeService, contractAddr common.Address) ([]byte, error) {
	key := concatContractAdminKey(native, contractAddr)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, err
	}
	if item == nil { //is not set
		return nil, nil
	}
	return item.Value, nil
}

func putContractAdmin(native *native.NativeService, contractAddr common.Address, adminGID []byte) error {
	key := concatContractAdminKey(native, contractAddr)
	utils.PutBytes(native, key, adminGID)
	return nil
}

//type(this.contractAddr.RoleFunc.role) = roleFuncs
func concatRoleFuncKey(native *native.NativeService, contractAddr common.Address, role []byte) []byte {
	this := native.ContextRef.CurrentContext().ContractAddress
	roleFuncKey := append(this[:], contractAddr[:]...)
	roleFuncKey = append(roleFuncKey, PreRoleFunc...)
	roleFuncKey = append(roleFuncKey, role...)

	return roleFuncKey
}

func getRoleFunc(native *native.NativeService, contractAddr common.Address, role []byte) (*roleFuncs, error) {
	key := concatRoleFuncKey(native, contractAddr, role)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, err
	}
	if item == nil { //is not set
		return nil, nil
	}
	rd := bytes.NewReader(item.Value)
	rF := new(roleFuncs)
	err = rF.Deserialize(rd)
	if err != nil {
		return nil, fmt.Errorf("deserialize roleFuncs object failed. data: %x", item.Value)
	}
	return rF, nil
}

func putRoleFunc(native *native.NativeService, contractAddr common.Address, role []byte, funcs *roleFuncs) error {
	key := concatRoleFuncKey(native, contractAddr, role)
	bf := new(bytes.Buffer)
	err := funcs.Serialize(bf)
	if err != nil {
		return fmt.Errorf("serialize roleFuncs failed, caused by %v", err)
	}
	utils.PutBytes(native, key, bf.Bytes())
	return nil
}

//type(this.contractAddr.RoleP.GID) = roleTokens
func concatGIDTokenKey(native *native.NativeService, contractAddr common.Address, GID []byte) []byte {
	this := native.ContextRef.CurrentContext().ContractAddress
	tokenKey := append(this[:], contractAddr[:]...)
	tokenKey = append(tokenKey, PreRoleToken...)
	tokenKey = append(tokenKey, GID...)

	return tokenKey
}

func getGIDToken(native *native.NativeService, contractAddr common.Address, GID []byte) (*roleTokens, error) {
	key := concatGIDTokenKey(native, contractAddr, GID)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, err
	}
	if item == nil { //is not set
		return nil, nil
	}
	rd := bytes.NewReader(item.Value)
	rT := new(roleTokens)
	err = rT.Deserialize(rd)
	if err != nil {
		return nil, fmt.Errorf("deserialize roleTokens object failed. data: %x", item.Value)
	}
	return rT, nil
}

func putGIDToken(native *native.NativeService, contractAddr common.Address, GID []byte, tokens *roleTokens) error {
	key := concatGIDTokenKey(native, contractAddr, GID)
	bf := new(bytes.Buffer)
	err := tokens.Serialize(bf)
	if err != nil {
		return fmt.Errorf("serialize roleFuncs failed, caused by %v", err)
	}
	utils.PutBytes(native, key, bf.Bytes())
	return nil
}

//type(this.contractAddr.DelegateStatus.GID)
func concatDelegateStatusKey(native *native.NativeService, contractAddr common.Address, GID []byte) []byte {
	this := native.ContextRef.CurrentContext().ContractAddress
	key := append(this[:], contractAddr[:]...)
	key = append(key, PreDelegateStatus...)
	key = append(key, GID...)

	return key
}

func getDelegateStatus(native *native.NativeService, contractAddr common.Address, GID []byte) (*Status, error) {
	key := concatDelegateStatusKey(native, contractAddr, GID)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	status := new(Status)
	rd := bytes.NewReader(item.Value)
	err = status.Deserialize(rd)
	if err != nil {
		return nil, fmt.Errorf("deserialize Status object failed. data: %x", item.Value)
	}
	return status, nil
}

func putDelegateStatus(native *native.NativeService, contractAddr common.Address, GID []byte, status *Status) error {
	key := concatDelegateStatusKey(native, contractAddr, GID)
	bf := new(bytes.Buffer)
	err := status.Serialize(bf)
	if err != nil {
		return fmt.Errorf("serialize Status failed, caused by %v", err)
	}
	utils.PutBytes(native, key, bf.Bytes())
	return nil
}

//remote duplicates in the slice of string
func stringSliceUniq(s []string) []string {
	smap := make(map[string]int)
	for i, str := range s {
		if str == "" {
			continue
		}
		smap[str] = i
	}
	ret := make([]string, len(smap))
	i := 0
	for str, _ := range smap {
		ret[i] = str
		i++
	}
	return ret
}

func pushEvent(native *native.NativeService, s interface{}) {
	event := new(event.NotifyEventInfo)
	event.ContractAddress = native.ContextRef.CurrentContext().ContractAddress
	event.States = s
	native.Notifications = append(native.Notifications, event)
}

func serializeAddress(w io.Writer, addr common.Address) error {
	err := serialization.WriteVarBytes(w, addr[:])
	if err != nil {
		return err
	}
	return nil
}
