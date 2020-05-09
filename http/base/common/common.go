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

// Package common privides functions for http handler call
package common

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/common/log"
	"github.com/imZhuFei/zeepin/common/serialization"
	"github.com/imZhuFei/zeepin/core/payload"
	"github.com/imZhuFei/zeepin/core/types"
	"github.com/imZhuFei/zeepin/embed/simulator"
	ontErrors "github.com/imZhuFei/zeepin/errors"
	bactor "github.com/imZhuFei/zeepin/http/base/actor"
	"github.com/imZhuFei/zeepin/smartcontract/event"
	embed "github.com/imZhuFei/zeepin/smartcontract/service/native/embed"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/utils"
	"github.com/imZhuFei/zeepin/smartcontract/service/wasmvm"
	cstates "github.com/imZhuFei/zeepin/smartcontract/states"
	"github.com/imZhuFei/zeepin/vm/wasmvm/exec"
	"github.com/zeepin/zeepinchain-crypto/keypair"
)

const MAX_SEARCH_HEIGHT uint32 = 100

type BalanceOfRsp struct {
	Zpt  string `json:"zpt"`
	Gala string `json:"gala"`
}

type MerkleProof struct {
	Type             string
	TransactionsRoot string
	BlockHeight      uint32
	CurBlockRoot     string
	CurBlockHeight   uint32
	TargetHashes     []string
}

type LogEventArgs struct {
	TxHash          string
	ContractAddress string
	Message         string
}

type ExecuteNotify struct {
	TxHash      string
	State       byte
	GasConsumed uint64
	Notify      []NotifyEventInfo
}

type NotifyEventInfo struct {
	ContractAddress string
	States          interface{}
}

type TxAttributeInfo struct {
	Usage types.TransactionAttributeUsage
	Data  string
}

type AmountMap struct {
	Key   common.Uint256
	Value common.Fixed64
}

type Fee struct {
	Amount common.Fixed64
	Payer  string
}

type Sig struct {
	PubKeys []string
	M       uint16
	SigData []string
}
type Transactions struct {
	Version    byte
	Nonce      uint32
	GasPrice   uint64
	GasLimit   uint64
	Payer      string
	TxType     types.TransactionType
	Payload    PayloadInfo
	Attributes []TxAttributeInfo
	Sigs       []Sig
	Hash       string
	Height     uint32
}

type BlockHead struct {
	Version          uint32
	PrevBlockHash    string
	TransactionsRoot string
	BlockRoot        string
	Timestamp        uint32
	Height           uint32
	ConsensusData    uint64
	ConsensusPayload string
	NextBookkeeper   string

	Bookkeepers []string
	SigData     []string

	Hash string
}

type BlockInfo struct {
	Hash         string
	Size         int
	Header       *BlockHead
	Transactions []*Transactions
}

type NodeInfo struct {
	NodeState   uint   // node status
	NodePort    uint16 // The nodes's port
	ID          uint64 // The nodes's id
	NodeTime    int64
	NodeVersion uint32   // The network protocol the node used
	NodeType    uint64   // The services the node supplied
	Relay       bool     // The relay capability of the node (merge into capbility flag)
	Height      uint32   // The node latest block height
	TxnCnt      []uint32 // The transactions in pool
	//RxTxnCnt uint64 // The transaction received by this node
}

type ConsensusInfo struct {
	// TODO
}

type TXNAttrInfo struct {
	Height  uint32
	Type    int
	ErrCode int
}

type TXNEntryInfo struct {
	State []TXNAttrInfo // the result from each validator
}

func GetLogEvent(obj *event.LogEventArgs) (map[string]bool, LogEventArgs) {
	hash := obj.TxHash
	addr := obj.ContractAddress.ToHexString()
	contractAddrs := map[string]bool{addr: true}
	return contractAddrs, LogEventArgs{hash.ToHexString(), addr, obj.Message}
}

func GetExecuteNotify(obj *event.ExecuteNotify) (map[string]bool, ExecuteNotify) {
	evts := []NotifyEventInfo{}
	var contractAddrs = make(map[string]bool)
	for _, v := range obj.Notify {
		evts = append(evts, NotifyEventInfo{v.ContractAddress.ToHexString(), v.States})
		contractAddrs[v.ContractAddress.ToHexString()] = true
	}
	txhash := obj.TxHash.ToHexString()
	return contractAddrs, ExecuteNotify{txhash, obj.State, obj.GasConsumed, evts}
}

func TransArryByteToHexString(ptx *types.Transaction) *Transactions {
	trans := new(Transactions)
	trans.TxType = ptx.TxType
	trans.Nonce = ptx.Nonce
	trans.GasLimit = ptx.GasLimit
	trans.GasPrice = ptx.GasPrice
	trans.Payer = ptx.Payer.ToBase58()
	trans.Payload = TransPayloadToHex(ptx.Payload)

	trans.Attributes = make([]TxAttributeInfo, 0)
	trans.Sigs = []Sig{}
	for _, sig := range ptx.Sigs {
		e := Sig{M: sig.M}
		for i := 0; i < len(sig.PubKeys); i++ {
			key := keypair.SerializePublicKey(sig.PubKeys[i])
			e.PubKeys = append(e.PubKeys, common.ToHexString(key))
		}
		for i := 0; i < len(sig.SigData); i++ {
			e.SigData = append(e.SigData, common.ToHexString(sig.SigData[i]))
		}
		trans.Sigs = append(trans.Sigs, e)
	}

	mhash := ptx.Hash()
	trans.Hash = mhash.ToHexString()
	return trans
}

func SendTxToPool(txn *types.Transaction) (ontErrors.ErrCode, string) {
	if errCode, desc := bactor.AppendTxToPool(txn); errCode != ontErrors.ErrNoError {
		log.Warn("TxnPool verify error:", errCode.Error())
		return errCode, desc
	}
	return ontErrors.ErrNoError, ""
}

func GetBlockInfo(block *types.Block) BlockInfo {
	hash := block.Hash()
	var bookkeepers = []string{}
	var sigData = []string{}
	for i := 0; i < len(block.Header.SigData); i++ {
		s := common.ToHexString(block.Header.SigData[i])
		sigData = append(sigData, s)
	}
	for i := 0; i < len(block.Header.Bookkeepers); i++ {
		e := block.Header.Bookkeepers[i]
		key := keypair.SerializePublicKey(e)
		bookkeepers = append(bookkeepers, common.ToHexString(key))
	}

	blockHead := &BlockHead{
		Version:          block.Header.Version,
		PrevBlockHash:    block.Header.PrevBlockHash.ToHexString(),
		TransactionsRoot: block.Header.TransactionsRoot.ToHexString(),
		BlockRoot:        block.Header.BlockRoot.ToHexString(),
		Timestamp:        block.Header.Timestamp,
		Height:           block.Header.Height,
		ConsensusData:    block.Header.ConsensusData,
		ConsensusPayload: common.ToHexString(block.Header.ConsensusPayload),
		NextBookkeeper:   block.Header.NextBookkeeper.ToBase58(),
		Bookkeepers:      bookkeepers,
		SigData:          sigData,
		Hash:             hash.ToHexString(),
	}

	trans := make([]*Transactions, len(block.Transactions))
	for i := 0; i < len(block.Transactions); i++ {
		trans[i] = TransArryByteToHexString(block.Transactions[i])
	}

	b := BlockInfo{
		Hash:         hash.ToHexString(),
		Size:         len(block.ToArray()),
		Header:       blockHead,
		Transactions: trans,
	}
	return b
}

func GetBalance(address common.Address) (*BalanceOfRsp, error) {
	zpt, err := GetContractBalance(0, utils.ZptContractAddress, address)
	if err != nil {
		return nil, fmt.Errorf("get zpt balance error:%s", err)
	}
	gala, err := GetContractBalance(0, utils.GalaContractAddress, address)
	if err != nil {
		return nil, fmt.Errorf("get gala balance error:%s", err)
	}
	return &BalanceOfRsp{
		Zpt:  fmt.Sprintf("%d", zpt),
		Gala: fmt.Sprintf("%d", gala),
	}, nil
}

func GetAllowance(asset string, from, to common.Address) (string, error) {
	var contractAddr common.Address
	switch strings.ToLower(asset) {
	case "zpt":
		contractAddr = utils.ZptContractAddress
	case "gala":
		contractAddr = utils.GalaContractAddress
	default:
		return "", fmt.Errorf("unsupport asset")
	}
	allowance, err := GetContractAllowance(0, contractAddr, from, to)
	if err != nil {
		return "", fmt.Errorf("get allowance error:%s", err)
	}
	return fmt.Sprintf("%v", allowance), nil
}

func GetContractBalance(cVersion byte, contractAddr, accAddr common.Address) (uint64, error) {
	mutable, err := NewNativeInvokeTransaction(0, 0, contractAddr, cVersion, "balanceOf", []interface{}{accAddr[:]})
	if err != nil {
		return 0, fmt.Errorf("NewNativeInvokeTransaction error:%s", err)
	}
	tx, err := mutable.IntoImmutable()
	if err != nil {
		return 0, err
	}
	result, err := bactor.PreExecuteContract(tx)
	if err != nil {
		return 0, fmt.Errorf("PrepareInvokeContract error:%s", err)
	}
	if result.State == 0 {
		return 0, fmt.Errorf("prepare invoke failed")
	}
	data, err := hex.DecodeString(result.Result.(string))
	if err != nil {
		return 0, fmt.Errorf("hex.DecodeString error:%s", err)
	}

	balance := common.BigIntFromEmbeddedBytes(data)
	return balance.Uint64(), nil
}

func GetContractAllowance(cVersion byte, contractAddr, fromAddr, toAddr common.Address) (uint64, error) {
	type allowanceStruct struct {
		From common.Address
		To   common.Address
	}
	mutable, err := NewNativeInvokeTransaction(0, 0, contractAddr, cVersion, "allowance",
		[]interface{}{&allowanceStruct{
			From: fromAddr,
			To:   toAddr,
		}})
	if err != nil {
		return 0, fmt.Errorf("NewNativeInvokeTransaction error:%s", err)
	}
	tx, err := mutable.IntoImmutable()
	if err != nil {
		return 0, err
	}
	result, err := bactor.PreExecuteContract(tx)
	if err != nil {
		return 0, fmt.Errorf("PrepareInvokeContract error:%s", err)
	}
	if result.State == 0 {
		return 0, fmt.Errorf("prepare invoke failed")
	}
	data, err := hex.DecodeString(result.Result.(string))
	if err != nil {
		return 0, fmt.Errorf("hex.DecodeString error:%s", err)
	}
	allowance := common.BigIntFromEmbeddedBytes(data)
	return allowance.Uint64(), nil
}

func GetGasPrice() (map[string]interface{}, error) {
	start := bactor.GetCurrentBlockHeight()
	var gasPrice uint64 = 0
	var height uint32 = 0
	var end uint32 = 0
	if start > MAX_SEARCH_HEIGHT {
		end = start - MAX_SEARCH_HEIGHT
	}
	for i := start; i >= end; i-- {
		head, err := bactor.GetHeaderByHeight(i)
		if err == nil && head.TransactionsRoot != common.UINT256_EMPTY {
			height = i
			blk, err := bactor.GetBlockByHeight(i)
			if err != nil {
				return nil, err
			}
			for _, v := range blk.Transactions {
				gasPrice += v.GasPrice
			}
			gasPrice = gasPrice / uint64(len(blk.Transactions))
			break
		}
	}
	result := map[string]interface{}{"gasprice": gasPrice, "height": height}
	return result, nil
}

func GetBlockTransactions(block *types.Block) interface{} {
	trans := make([]string, len(block.Transactions))
	for i := 0; i < len(block.Transactions); i++ {
		t := block.Transactions[i].Hash()
		trans[i] = t.ToHexString()
	}
	hash := block.Hash()
	type BlockTransactions struct {
		Hash         string
		Height       uint32
		Transactions []string
	}
	b := BlockTransactions{
		Hash:         hash.ToHexString(),
		Height:       block.Header.Height,
		Transactions: trans,
	}
	return b
}

//NewNativeInvokeTransaction return native contract invoke transaction
func NewNativeInvokeTransaction(gasPirce, gasLimit uint64, contractAddress common.Address, version byte, method string, params []interface{}) (*types.MutableTransaction, error) {
	invokeCode, err := BuildNativeInvokeCode(contractAddress, version, method, params)
	if err != nil {
		return nil, err
	}
	return NewSmartContractTransaction(gasPirce, gasLimit, invokeCode, 0)
}

func NewEmbeddedInvokeTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, params []interface{}) (*types.MutableTransaction, error) {
	invokeCode, err := BuildEmbeddedInvokeCode(contractAddress, params)
	if err != nil {
		return nil, err
	}
	return NewSmartContractTransaction(gasPrice, gasLimit, invokeCode, 0)
}

func NewWASMVMInvokeTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, methodName string, paramType wasmvm.ParamType, version byte, params []interface{}) (*types.MutableTransaction, error) {
	invokeCode, err := BuildWasmVMInvokeCode(contractAddress, methodName, paramType, version, params)
	if err != nil {
		return nil, err
	}
	immut, err := NewSmartContractTransaction(gasPrice, gasLimit, invokeCode, 1)
	return immut, err
}

func NewSmartContractTransaction(gasPrice, gasLimit uint64, invokeCode []byte, attr byte) (*types.MutableTransaction, error) {
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice:   gasPrice,
		GasLimit:   gasLimit,
		TxType:     types.Invoke,
		Nonce:      uint32(time.Now().Unix()),
		Payload:    invokePayload,
		Attributes: attr,
		Sigs:       nil,
	}
	return tx, nil
}

func BuildNativeInvokeCode(contractAddress common.Address, version byte, method string, params []interface{}) ([]byte, error) {
	builder := simulator.NewParamsBuilder(new(bytes.Buffer))
	err := BuildEmbeddedParam(builder, params)
	if err != nil {
		return nil, err
	}
	builder.EmitPushByteArray([]byte(method))
	builder.EmitPushByteArray(contractAddress[:])
	builder.EmitPushInteger(new(big.Int).SetInt64(int64(version)))
	builder.Emit(simulator.SYSCALL)
	builder.EmitPushByteArray([]byte(embed.NATIVE_INVOKE_NAME))
	return builder.ToArray(), nil
}

//BuildEmbeddedInvokeCode build Embed Invoke code for params
func BuildEmbeddedInvokeCode(smartContractAddress common.Address, params []interface{}) ([]byte, error) {
	builder := simulator.NewParamsBuilder(new(bytes.Buffer))
	err := BuildEmbeddedParam(builder, params)
	if err != nil {
		return nil, err
	}
	args := append(builder.ToArray(), 0x67)
	args = append(args, smartContractAddress[:]...)
	return args, nil
}

//buildEmbeddedParamInter build embed invoke param code
func BuildEmbeddedParam(builder *simulator.ParamsBuilder, smartContractParams []interface{}) error {
	//VM load params in reverse order
	for i := len(smartContractParams) - 1; i >= 0; i-- {
		switch v := smartContractParams[i].(type) {
		case bool:
			builder.EmitPushBool(v)
		case byte:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case int:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case uint:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case int32:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case uint32:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case int64:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case common.Fixed64:
			builder.EmitPushInteger(big.NewInt(int64(v.GetData())))
		case uint64:
			val := big.NewInt(0)
			builder.EmitPushInteger(val.SetUint64(uint64(v)))
		case string:
			builder.EmitPushByteArray([]byte(v))
		case *big.Int:
			builder.EmitPushInteger(v)
		case []byte:
			builder.EmitPushByteArray(v)
		case common.Address:
			builder.EmitPushByteArray(v[:])
		case common.Uint256:
			builder.EmitPushByteArray(v.ToArray())
		case []interface{}:
			err := BuildEmbeddedParam(builder, v)
			if err != nil {
				return err
			}
			builder.EmitPushInteger(big.NewInt(int64(len(v))))
			builder.Emit(simulator.PACK)
		default:
			object := reflect.ValueOf(v)
			kind := object.Kind().String()
			if kind == "ptr" {
				object = object.Elem()
				kind = object.Kind().String()
			}
			switch kind {
			case "slice":
				ps := make([]interface{}, 0)
				for i := 0; i < object.Len(); i++ {
					ps = append(ps, object.Index(i).Interface())
				}
				err := BuildEmbeddedParam(builder, []interface{}{ps})
				if err != nil {
					return err
				}
			case "struct":
				builder.EmitPushInteger(big.NewInt(0))
				builder.Emit(simulator.NEWSTRUCT)
				builder.Emit(simulator.TOALTSTACK)
				for i := 0; i < object.NumField(); i++ {
					field := object.Field(i)
					err := BuildEmbeddedParam(builder, []interface{}{field.Interface()})
					if err != nil {
						return err
					}
					builder.Emit(simulator.DUPFROMALTSTACK)
					builder.Emit(simulator.SWAP)
					builder.Emit(simulator.APPEND)
				}
				builder.Emit(simulator.FROMALTSTACK)
			default:
				return fmt.Errorf("unsupported param:%s", v)
			}
		}
	}
	return nil
}

//for wasm vm
//build param bytes for wasm contract
func buildWasmContractParam(params []interface{}, paramType wasmvm.ParamType) ([]byte, error) {
	switch paramType {
	case wasmvm.Json:
		args := make([]exec.Param, len(params))

		for i, param := range params {
			switch param.(type) {
			case string:
				arg := exec.Param{Ptype: "string", Pval: param.(string)}
				args[i] = arg
			case int:
				arg := exec.Param{Ptype: "int", Pval: strconv.Itoa(param.(int))}
				args[i] = arg
			case int64:
				arg := exec.Param{Ptype: "int64", Pval: strconv.FormatInt(param.(int64), 10)}
				args[i] = arg
			case []int:
				bf := bytes.NewBuffer(nil)
				array := param.([]int)
				for i, tmp := range array {
					bf.WriteString(strconv.Itoa(tmp))
					if i != len(array)-1 {
						bf.WriteString(",")
					}
				}
				arg := exec.Param{Ptype: "int_array", Pval: bf.String()}
				args[i] = arg
			case []int64:
				bf := bytes.NewBuffer(nil)
				array := param.([]int64)
				for i, tmp := range array {
					bf.WriteString(strconv.FormatInt(tmp, 10))
					if i != len(array)-1 {
						bf.WriteString(",")
					}
				}
				arg := exec.Param{Ptype: "int_array", Pval: bf.String()}
				args[i] = arg
			default:
				return nil, fmt.Errorf("not a supported type :%v\n", param)
			}
		}

		bs, err := json.Marshal(exec.Args{args})
		if err != nil {
			return nil, err
		}
		return bs, nil
	case wasmvm.Raw:
		bf := bytes.NewBuffer(nil)
		for _, param := range params {
			switch param.(type) {
			case string:
				tmp := bytes.NewBuffer(nil)
				serialization.WriteString(tmp, param.(string))
				bf.Write(tmp.Bytes())

			case int:
				tmpBytes := make([]byte, 4)
				binary.LittleEndian.PutUint32(tmpBytes, uint32(param.(int)))
				bf.Write(tmpBytes)

			case int64:
				tmpBytes := make([]byte, 8)
				binary.LittleEndian.PutUint64(tmpBytes, uint64(param.(int64)))
				bf.Write(tmpBytes)

			default:
				return nil, fmt.Errorf("not a supported type :%v\n", param)
			}
		}
		return bf.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

//BuildWasmVMInvokeCode return wasn vm invoke code
func BuildWasmVMInvokeCode(smartcodeAddress common.Address, methodName string, paramType wasmvm.ParamType, version byte, params []interface{}) ([]byte, error) {
	contract := &cstates.Contract{}
	contract.Address = smartcodeAddress
	contract.Method = methodName
	contract.Version = version

	argbytes, err := buildWasmContractParam(params, paramType)

	if err != nil {
		return nil, fmt.Errorf("build wasm contract param failed:%s", err)
	}
	contract.Args = argbytes
	bf := bytes.NewBuffer(nil)
	contract.Serialize(bf)
	return bf.Bytes(), nil
}
