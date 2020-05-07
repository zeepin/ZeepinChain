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

package gala

import (
	"bytes"
	"math/big"

	"fmt"

	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/common/constants"
	"github.com/imZhuFei/zeepin/common/serialization"
	scommon "github.com/imZhuFei/zeepin/core/store/common"
	"github.com/imZhuFei/zeepin/errors"
	"github.com/imZhuFei/zeepin/smartcontract/service/native"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/utils"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/zpt"
	"github.com/imZhuFei/zeepin/embed/simulator/types"
)

func InitGala() {
	native.Contracts[utils.GalaContractAddress] = RegisterGalaContract
}

func RegisterGalaContract(native *native.NativeService) {
	native.Register(zpt.INIT_NAME, GalaInit)
	native.Register(zpt.TRANSFER_NAME, GalaTransfer)
	native.Register(zpt.APPROVE_NAME, GalaApprove)
	native.Register(zpt.TRANSFERFROM_NAME, GalaTransferFrom)
	native.Register(zpt.NAME_NAME, GalaName)
	native.Register(zpt.SYMBOL_NAME, GalaSymbol)
	native.Register(zpt.DECIMALS_NAME, GalaDecimals)
	native.Register(zpt.TOTALSUPPLY_NAME, GalaTotalSupply)
	native.Register(zpt.BALANCEOF_NAME, GalaBalanceOf)
	native.Register(zpt.ALLOWANCE_NAME, GalaAllowance)
}

func GalaInit(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, zpt.GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if amount > 0 {
		return utils.BYTE_FALSE, errors.NewErr("Init gala has been completed!")
	}

	distribute := make(map[common.Address]uint64)
	buf, err := serialization.ReadVarBytes(bytes.NewBuffer(native.Input))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadVarBytes, contract params deserialize error!")
	}
	input := bytes.NewBuffer(buf)
	num, err := utils.ReadVarUint(input)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("read number error:%v", err)
	}
	sum := uint64(0)
	overflow := false
	for i := uint64(0); i < num; i++ {
		addr, err := utils.ReadAddress(input)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("read address error:%v", err)
		}
		value, err := utils.ReadVarUint(input)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("read value error:%v", err)
		}
		sum, overflow = common.SafeAdd(sum, value)
		if overflow {
			return utils.BYTE_FALSE, errors.NewErr("wrong config. overflow detected")
		}
		distribute[addr] += value
	}
	if sum != constants.GALA_TOTAL_SUPPLY-constants.GALA_UNBOUND_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("wrong config. total supply %d != %d", sum, constants.GALA_TOTAL_SUPPLY-constants.GALA_UNBOUND_SUPPLY)
	}

	for addr, val := range distribute {
		balanceKey := zpt.GenBalanceKey(contract, addr)
		item := utils.GenUInt64StorageItem(val)
		native.CloneCache.Add(scommon.ST_STORAGE, balanceKey, item)
		zpt.AddNotifications(native, contract, &zpt.State{To: addr, Value: val})
	}
	item := utils.GenUInt64StorageItem(constants.GALA_TOTAL_SUPPLY)
	unboundItem := utils.GenUInt64StorageItem(constants.GALA_UNBOUND_SUPPLY)
	native.CloneCache.Add(scommon.ST_STORAGE, zpt.GenTotalSupplyKey(contract), item)
	native.CloneCache.Add(scommon.ST_STORAGE, append(contract[:], utils.ZptContractAddress[:]...), unboundItem)
	zpt.AddNotifications(native, contract, &zpt.State{To: utils.ZptContractAddress, Value: constants.GALA_UNBOUND_SUPPLY})
	return utils.BYTE_TRUE, nil

}

func GalaTransfer(native *native.NativeService) ([]byte, error) {
	var transfers zpt.Transfers
	source := common.NewZeroCopySource(native.Input)
	if err := transfers.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GalaTransfer] Transfers deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for _, v := range transfers.States {
		if v.Value == 0 {
			continue
		}
		if v.Value > constants.GALA_TOTAL_SUPPLY {
			return utils.BYTE_FALSE, fmt.Errorf("transfer gala amount:%d over totalSupply:%d", v.Value, constants.GALA_TOTAL_SUPPLY)
		}
		if _, _, err := zpt.Transfer(native, contract, &v); err != nil {
			return utils.BYTE_FALSE, err
		}
		zpt.AddNotifications(native, contract, &v)
	}
	return utils.BYTE_TRUE, nil
}

func GalaApprove(native *native.NativeService) ([]byte, error) {
	var state zpt.State
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GalaApprove] state deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.GALA_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve gala amount:%d over totalSupply:%d", state.Value, constants.GALA_TOTAL_SUPPLY)
	}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("authentication failed!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	native.CloneCache.Add(scommon.ST_STORAGE, zpt.GenApproveKey(contract, state.From, state.To), utils.GenUInt64StorageItem(state.Value))
	return utils.BYTE_TRUE, nil
}

func GalaTransferFrom(native *native.NativeService) ([]byte, error) {
	var state zpt.TransferFrom
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[ZptTransferFrom] State deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.GALA_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve gala amount:%d over totalSupply:%d", state.Value, constants.GALA_TOTAL_SUPPLY)
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	if _, _, err := zpt.TransferedFrom(native, contract, &state); err != nil {
		return utils.BYTE_FALSE, err
	}
	zpt.AddNotifications(native, contract, &zpt.State{From: state.From, To: state.To, Value: state.Value})
	return utils.BYTE_TRUE, nil
}

func GalaName(native *native.NativeService) ([]byte, error) {
	return []byte(constants.GALA_NAME), nil
}

func GalaDecimals(native *native.NativeService) ([]byte, error) {
	return big.NewInt(int64(constants.GALA_DECIMALS)).Bytes(), nil
}

func GalaSymbol(native *native.NativeService) ([]byte, error) {
	return []byte(constants.GALA_SYMBOL), nil
}

func GalaTotalSupply(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, zpt.GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[ZptTotalSupply] get totalSupply error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func GalaBalanceOf(native *native.NativeService) ([]byte, error) {
	return zpt.GetBalanceValue(native, zpt.TRANSFER_FLAG)
}

func GalaAllowance(native *native.NativeService) ([]byte, error) {
	return zpt.GetBalanceValue(native, zpt.APPROVE_FLAG)
}
