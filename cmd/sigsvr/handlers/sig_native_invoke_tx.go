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

package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	"github.com/imZhuFei/zeepin/cmd/abi"
	clisvrcom "github.com/imZhuFei/zeepin/cmd/sigsvr/common"
	cliutil "github.com/imZhuFei/zeepin/cmd/utils"
	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/common/log"
)

type SigNativeInvokeTxReq struct {
	GasPrice uint64        `json:"gas_price"`
	GasLimit uint64        `json:"gas_limit"`
	Address  string        `json:"address"`
	Method   string        `json:"method"`
	Params   []interface{} `json:"params"`
	Version  byte          `json:"version"`
}

type SigNativeInvokeTxRsp struct {
	SignedTx string `json:"signed_tx"`
}

func SigNativeInvokeTx(req *clisvrcom.CliRpcRequest, resp *clisvrcom.CliRpcResponse) {
	rawReq := &SigNativeInvokeTxReq{}
	err := json.Unmarshal(req.Params, rawReq)
	if err != nil {
		log.Infof("Cli Qid:%s SigNativeInvokeTx json.Unmarshal SigNativeInvokeTxReq:%s error:%s", req.Qid, req.Params, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	contractAddr, err := common.AddressFromHexString(rawReq.Address)
	if err != nil {
		log.Infof("Cli Qid:%s SigNativeInvokeTx AddressParseFromBytes:%s error:%s", req.Qid, rawReq.Address, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	nativeAbi := abi.DefAbiMgr.GetNativeAbi(rawReq.Address)
	if nativeAbi == nil {
		resp.ErrorCode = clisvrcom.CLIERR_ABI_NOT_FOUND
		return
	}
	funcAbi := nativeAbi.GetFunc(rawReq.Method)
	if funcAbi == nil {
		resp.ErrorCode = clisvrcom.CLIERR_ABI_NOT_FOUND
		return
	}
	tx, err := cliutil.NewNativeInvokeTransaction(rawReq.GasPrice, rawReq.GasLimit, contractAddr, rawReq.Version, rawReq.Params, funcAbi)
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		resp.ErrorInfo = err.Error()
		return
	}
	if err != nil {
		log.Infof("Cli Qid:%s SigNativeInvokeTx InvokeNativeContractTx error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	signer := clisvrcom.DefAccount
	err = cliutil.SignTransaction(signer, tx)
	if err != nil {
		log.Infof("Cli Qid:%s SigNativeInvokeTx SignTransaction error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	immutable, err := tx.IntoImmutable()
	if err != nil {
		log.Infof("convert to immutable transaction error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	buf := bytes.NewBuffer(nil)
	err = immutable.Serialize(buf)
	if err != nil {
		log.Infof("Cli Qid:%s SigNativeInvokeTx tx Serialize error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	resp.Result = &SigNativeInvokeTxRsp{
		SignedTx: hex.EncodeToString(buf.Bytes()),
	}
}
