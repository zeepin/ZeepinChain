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
 * This file is part of The ZeepinChain library.
 */
package gid

import (
	"encoding/hex"
	"errors"

	com "github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/core/states"
	"github.com/imZhuFei/zeepin/core/store/common"
	"github.com/imZhuFei/zeepin/core/types"
	"github.com/imZhuFei/zeepin/smartcontract/service/native"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/utils"
	"github.com/zeepin/zeepinchain-crypto/keypair"
)

const flag_exist = 0x01

func checkIDExistence(srvc *native.NativeService, encID []byte) bool {
	val, err := srvc.CloneCache.Get(common.ST_STORAGE, encID)
	if err == nil {
		t, ok := val.(*states.StorageItem)
		if ok {
			if len(t.Value) > 0 && t.Value[0] == flag_exist {
				return true
			}
		}
	}
	return false
}

const (
	FIELD_PK byte = 1 + iota
	FIELD_ATTR
	FIELD_RECOVERY
)

func encodeID(id []byte) ([]byte, error) {
	length := len(id)
	if length == 0 || length > 255 {
		return nil, errors.New("encode GID error: invalid ID length")
	}
	enc := []byte{byte(length)}
	enc = append(enc, id...)
	return enc, nil
}

func decodeID(data []byte) ([]byte, error) {
	if len(data) == 0 || len(data) != int(data[0])+1 {
		return nil, errors.New("decode GID error: invalid data length")
	}
	return data[1:], nil
}

func setRecovery(srvc *native.NativeService, encID []byte, recovery com.Address) error {
	key := append(encID, FIELD_RECOVERY)
	val := &states.StorageItem{Value: recovery[:]}
	srvc.CloneCache.Add(common.ST_STORAGE, key, val)
	return nil
}

func getRecovery(srvc *native.NativeService, encID []byte) ([]byte, error) {
	key := append(encID, FIELD_RECOVERY)
	item, err := utils.GetStorageItem(srvc, key)
	if err != nil {
		return nil, errors.New("get recovery error: " + err.Error())
	} else if item == nil {
		return nil, nil
	}
	return item.Value, nil
}

func checkWitness(srvc *native.NativeService, key []byte) error {
	// try as if key is a public key
	pk, err := keypair.DeserializePublicKey(key)
	if err == nil {
		addr := types.AddressFromPubKey(pk)
		if srvc.ContextRef.CheckWitness(addr) {
			return nil
		}
	}

	// try as if key is an address
	addr, err := com.AddressParseFromBytes(key)
	if srvc.ContextRef.CheckWitness(addr) {
		return nil
	}

	return errors.New("check witness failed, " + hex.EncodeToString(key))
}
