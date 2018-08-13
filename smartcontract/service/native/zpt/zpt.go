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

package zpt

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/mileschao/ZeepinChain/common"
	"github.com/mileschao/ZeepinChain/common/constants"
	"github.com/mileschao/ZeepinChain/common/log"
	"github.com/mileschao/ZeepinChain/common/serialization"
	scommon "github.com/mileschao/ZeepinChain/core/store/common"
	"github.com/mileschao/ZeepinChain/errors"
	"github.com/mileschao/ZeepinChain/smartcontract/service/native"
	"github.com/mileschao/ZeepinChain/smartcontract/service/native/utils"
	"github.com/mileschao/ZeepinChain/vm/neovm/types"
)

const (
	TRANSFER_FLAG byte = 1
	APPROVE_FLAG  byte = 2
)

func InitZpt() {
	native.Contracts[utils.ZptContractAddress] = RegisterZptContract
}

func RegisterZptContract(native *native.NativeService) {
	native.Register(INIT_NAME, ZptInit)
	native.Register(TRANSFER_NAME, ZptTransfer)
	native.Register(APPROVE_NAME, ZptApprove)
	native.Register(TRANSFERFROM_NAME, ZptTransferFrom)
	native.Register(NAME_NAME, ZptName)
	native.Register(SYMBOL_NAME, ZptSymbol)
	native.Register(DECIMALS_NAME, ZptDecimals)
	native.Register(TOTALSUPPLY_NAME, ZptTotalSupply)
	native.Register(BALANCEOF_NAME, ZptBalanceOf)
	native.Register(ALLOWANCE_NAME, ZptAllowance)
}

func ZptInit(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if amount > 0 {
		return utils.BYTE_FALSE, errors.NewErr("Init zpt has been completed!")
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
	if sum != constants.ZPT_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("wrong config. total supply %d != %d", sum, constants.ZPT_TOTAL_SUPPLY)
	}

	for addr, val := range distribute {
		balanceKey := GenBalanceKey(contract, addr)
		item := utils.GenUInt64StorageItem(val)
		native.CloneCache.Add(scommon.ST_STORAGE, balanceKey, item)
		AddNotifications(native, contract, &State{To: addr, Value: val})
	}
	native.CloneCache.Add(scommon.ST_STORAGE, GenTotalSupplyKey(contract), utils.GenUInt64StorageItem(constants.ZPT_TOTAL_SUPPLY))

	return utils.BYTE_TRUE, nil
}

func ZptTransfer(native *native.NativeService) ([]byte, error) {
	transfers := new(Transfers)
	if err := transfers.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[Transfer] Transfers deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for _, v := range transfers.States {
		if v.Value == 0 {
			continue
		}
		if v.Value > constants.ZPT_TOTAL_SUPPLY {
			return utils.BYTE_FALSE, fmt.Errorf("transfer zpt amount:%d over totalSupply:%d", v.Value, constants.ZPT_TOTAL_SUPPLY)
		}
		fromBalance, toBalance, err := Transfer(native, contract, v)
		if err != nil {
			return utils.BYTE_FALSE, err
		}

		if err := grantGala(native, contract, v.From, fromBalance); err != nil {
			return utils.BYTE_FALSE, err
		}

		if err := grantGala(native, contract, v.To, toBalance); err != nil {
			return utils.BYTE_FALSE, err
		}

		AddNotifications(native, contract, v)
	}
	return utils.BYTE_TRUE, nil
}

func ZptTransferFrom(native *native.NativeService) ([]byte, error) {
	state := new(TransferFrom)
	if err := state.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[ZptTransferFrom] State deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.ZPT_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("transferFrom zpt amount:%d over totalSupply:%d", state.Value, constants.ZPT_TOTAL_SUPPLY)
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	fromBalance, toBalance, err := TransferedFrom(native, contract, state)
	if err != nil {
		return utils.BYTE_FALSE, err
	}
	if err := grantGala(native, contract, state.From, fromBalance); err != nil {
		return utils.BYTE_FALSE, err
	}
	if err := grantGala(native, contract, state.To, toBalance); err != nil {
		return utils.BYTE_FALSE, err
	}
	AddNotifications(native, contract, &State{From: state.From, To: state.To, Value: state.Value})
	return utils.BYTE_TRUE, nil
}

func ZptApprove(native *native.NativeService) ([]byte, error) {
	state := new(State)
	if err := state.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GalaApprove] state deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.ZPT_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve zpt amount:%d over totalSupply:%d", state.Value, constants.ZPT_TOTAL_SUPPLY)
	}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("authentication failed!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	native.CloneCache.Add(scommon.ST_STORAGE, GenApproveKey(contract, state.From, state.To), utils.GenUInt64StorageItem(state.Value))
	return utils.BYTE_TRUE, nil
}

func ZptName(native *native.NativeService) ([]byte, error) {
	return []byte(constants.ZPT_NAME), nil
}

func ZptDecimals(native *native.NativeService) ([]byte, error) {
	return types.BigIntToBytes(big.NewInt(int64(constants.ZPT_DECIMALS))), nil
}

func ZptSymbol(native *native.NativeService) ([]byte, error) {
	return []byte(constants.ZPT_SYMBOL), nil
}

func ZptTotalSupply(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[ZptTotalSupply] get totalSupply error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func ZptBalanceOf(native *native.NativeService) ([]byte, error) {
	return GetBalanceValue(native, TRANSFER_FLAG)
}

func ZptAllowance(native *native.NativeService) ([]byte, error) {
	return GetBalanceValue(native, APPROVE_FLAG)
}

func GetBalanceValue(native *native.NativeService, flag byte) ([]byte, error) {
	var key []byte
	buf := bytes.NewBuffer(native.Input)
	from, err := utils.ReadAddress(buf)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] get from address error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	if flag == APPROVE_FLAG {
		to, err := utils.ReadAddress(buf)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] get from address error!")
		}
		key = GenApproveKey(contract, from, to)
	} else if flag == TRANSFER_FLAG {
		key = GenBalanceKey(contract, from)
	}
	amount, err := utils.GetStorageUInt64(native, key)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] address parse error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func grantGala(native *native.NativeService, contract, address common.Address, balance uint64) error {
	startOffset, err := getUnboundOffset(native, contract, address)
	if err != nil {
		return err
	}
	if native.Time <= constants.GENESIS_BLOCK_TIMESTAMP {
		return nil
	}
	endOffset := native.Time - constants.GENESIS_BLOCK_TIMESTAMP
	log.Debugf("grantGala: startOffset: %d, endOffset:%d", startOffset, endOffset)
	if endOffset < startOffset {
		errstr := fmt.Sprintf("grant Gala error: wrong timestamp endOffset: %d < startOffset: %d", endOffset, startOffset)
		log.Error(errstr)
		return errors.NewErr(errstr)
	} else if endOffset == startOffset {
		return nil
	}

	if balance != 0 {
		value := utils.CalcUnbindGala(balance, startOffset, endOffset)

		args, err := getApproveArgs(native, contract, utils.GalaContractAddress, address, value)
		if err != nil {
			return err
		}

		if _, err := native.NativeCall(utils.GalaContractAddress, "approve", args); err != nil {
			return err
		}
	}

	native.CloneCache.Add(scommon.ST_STORAGE, genAddressUnboundOffsetKey(contract, address), utils.GenUInt32StorageItem(endOffset))
	return nil
}

func getApproveArgs(native *native.NativeService, contract, galaContract, address common.Address, value uint64) ([]byte, error) {
	bf := new(bytes.Buffer)
	approve := &State{
		From:  contract,
		To:    address,
		Value: value,
	}

	stateValue, err := utils.GetStorageUInt64(native, GenApproveKey(galaContract, approve.From, approve.To))
	if err != nil {
		return nil, err
	}

	approve.Value += stateValue

	if err := approve.Serialize(bf); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}
