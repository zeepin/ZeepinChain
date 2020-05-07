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

package embed

import (
	"bytes"

	"fmt"

	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/core/states"
	scommon "github.com/imZhuFei/zeepin/core/store/common"
	vm "github.com/imZhuFei/zeepin/embed/simulator"
	"github.com/imZhuFei/zeepin/errors"
)

// StoragePut put smart contract storage item to cache
func StoragePut(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	if vm.EvaluationStackCount(engine) < 3 {
		return errors.NewErr("[Context] Too few input parameters ")
	}
	context, err := getContext(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StoragePut] get pop context error!")
	}
	if context.IsReadOnly {
		return fmt.Errorf("%s", "[StoragePut] storage read only!")
	}
	if err := checkStorageContext(service, context); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StoragePut] check context error!")
	}

	key, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	if len(key) > 1024 {
		return errors.NewErr("[StoragePut] Storage key to long")
	}

	value, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	service.CloneCache.Add(scommon.ST_STORAGE, getStorageKey(context.Address, key), &states.StorageItem{Value: value})
	return nil
}

// StorageDelete delete smart contract storage item from cache
func StorageDelete(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	if vm.EvaluationStackCount(engine) < 2 {
		return errors.NewErr("[Context] Too few input parameters ")
	}
	context, err := getContext(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageDelete] get pop context error!")
	}
	if context.IsReadOnly {
		return fmt.Errorf("%s", "[StorageDelete] storage read only!")
	}
	if err := checkStorageContext(service, context); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageDelete] check context error!")
	}
	ba, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	service.CloneCache.Delete(scommon.ST_STORAGE, getStorageKey(context.Address, ba))

	return nil
}

// StorageGet push smart contract storage item from cache to vm stack
func StorageGet(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	if vm.EvaluationStackCount(engine) < 2 {
		return errors.NewErr("[Context] Too few input parameters ")
	}
	context, err := getContext(engine)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[StorageGet] get pop context error!")
	}
	ba, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	item, err := service.CloneCache.Get(scommon.ST_STORAGE, getStorageKey(context.Address, ba))
	if err != nil {
		return err
	}

	if item == nil {
		vm.PushData(engine, []byte{})
	} else {
		vm.PushData(engine, item.(*states.StorageItem).Value)
	}
	return nil
}

// StorageGetContext push smart contract storage context to vm stack
func StorageGetContext(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	vm.PushData(engine, NewStorageContext(service.ContextRef.CurrentContext().ContractAddress))
	return nil
}

func StorageGetReadOnlyContext(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	context := NewStorageContext(service.ContextRef.CurrentContext().ContractAddress)
	context.IsReadOnly = true
	vm.PushData(engine, context)
	return nil
}

func checkStorageContext(service *EmbeddedService, context *StorageContext) error {
	item, err := service.CloneCache.Get(scommon.ST_CONTRACT, context.Address[:])
	if err != nil || item == nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[CheckStorageContext] get context fail!")
	}
	return nil
}

func getContext(engine *vm.ExecutionEngine) (*StorageContext, error) {
	opInterface, err := vm.PopInteropInterface(engine)
	if err != nil {
		return nil, err
	}
	if opInterface == nil {
		return nil, errors.NewErr("[Context] Get storageContext nil")
	}
	context, ok := opInterface.(*StorageContext)
	if !ok {
		return nil, errors.NewErr("[Context] Get storageContext invalid")
	}
	return context, nil
}

func getStorageKey(address common.Address, key []byte) []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(address[:])
	buf.Write(key)
	return buf.Bytes()
}
