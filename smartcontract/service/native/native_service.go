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

package native

import (
	"bytes"
	"fmt"

	"github.com/mileschao/ZeepinChain/common"
	"github.com/mileschao/ZeepinChain/core/types"
	"github.com/mileschao/ZeepinChain/errors"
	"github.com/mileschao/ZeepinChain/smartcontract/context"
	"github.com/mileschao/ZeepinChain/smartcontract/event"
	"github.com/mileschao/ZeepinChain/smartcontract/states"
	sstates "github.com/mileschao/ZeepinChain/smartcontract/states"
	"github.com/mileschao/ZeepinChain/smartcontract/storage"
)

type (
	Handler         func(native *NativeService) ([]byte, error)
	RegisterService func(native *NativeService)
)

var (
	Contracts = make(map[common.Address]RegisterService)
)

// Native service struct
// Invoke a native smart contract, new a native service
type NativeService struct {
	CloneCache    *storage.CloneCache
	ServiceMap    map[string]Handler
	Notifications []*event.NotifyEventInfo
	Code          []byte
	Input         []byte
	Tx            *types.Transaction
	Height        uint32
	Time          uint32
	ContextRef    context.ContextRef
}

func (this *NativeService) Register(methodName string, handler Handler) {
	this.ServiceMap[methodName] = handler
}

func (this *NativeService) Invoke() (interface{}, error) {
	bf := bytes.NewBuffer(this.Code)
	contract := new(sstates.Contract)
	if err := contract.Deserialize(bf); err != nil {
		return false, err
	}
	services, ok := Contracts[contract.Address]
	if !ok {
		return false, fmt.Errorf("Native contract address %x haven't been registered.", contract.Address)
	}
	services(this)
	service, ok := this.ServiceMap[contract.Method]
	if !ok {
		return false, fmt.Errorf("Native contract %x doesn't support this function %s.",
			contract.Address, contract.Method)
	}
	args := this.Input
	this.Input = contract.Args
	this.ContextRef.PushContext(&context.Context{ContractAddress: contract.Address})
	notifications := this.Notifications
	this.Notifications = []*event.NotifyEventInfo{}
	result, err := service(this)
	if err != nil {
		return result, errors.NewDetailErr(err, errors.ErrNoCode, "[Invoke] Native serivce function execute error!")
	}
	this.ContextRef.PopContext()
	this.ContextRef.PushNotifications(this.Notifications)
	this.Notifications = notifications
	this.Input = args
	return result, nil
}

func (this *NativeService) NativeCall(address common.Address, method string, args []byte) (interface{}, error) {
	bf := new(bytes.Buffer)
	c := states.Contract{
		Address: address,
		Method:  method,
		Args:    args,
	}
	if err := c.Serialize(bf); err != nil {
		return nil, err
	}
	this.Code = bf.Bytes()
	return this.Invoke()
}
