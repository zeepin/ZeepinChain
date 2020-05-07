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
	"encoding/hex"
	"encoding/json"
	"fmt"

	clisvrcom "github.com/imZhuFei/zeepin/cmd/sigsvr/common"
	cliutil "github.com/imZhuFei/zeepin/cmd/utils"
	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/common/log"
	httpcom "github.com/imZhuFei/zeepin/http/base/common"
)

type SigEmbededInvokeTxReq struct {
	GasPrice uint64        `json:"gas_price"`
	GasLimit uint64        `json:"gas_limit"`
	Address  string        `json:"address"`
	Params   []interface{} `json:"params"`
}

type SigEmbededInvokeTxRsp struct {
	SignedTx string `json:"signed_tx"`
}

func SigEmbededInvokeTx(req *clisvrcom.CliRpcRequest, resp *clisvrcom.CliRpcResponse) {
	rawReq := &SigEmbededInvokeTxReq{}
	err := json.Unmarshal(req.Params, rawReq)
	if err != nil {
		log.Infof("SigEmbededInvokeTx json.Unmarshal SigEmbededInvokeTxReq:%s error:%s", req.Params, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	params, err := cliutil.ParseEmbeddedInvokeParams(rawReq.Params)
	if err != nil {
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		resp.ErrorInfo = fmt.Sprintf("ParseEmbeddedInvokeParams error:%s", err)
		return
	}
	contAddr, err := common.AddressFromHexString(rawReq.Address)
	if err != nil {
		log.Infof("Cli Qid:%s SigEmbededInvokeTx AddressParseFromBytes:%s error:%s", req.Qid, rawReq.Address, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	mutable, err := httpcom.NewEmbeddedInvokeTransaction(rawReq.GasPrice, rawReq.GasLimit, contAddr, params)
	if err != nil {
		log.Infof("Cli Qid:%s SigEmbededInvokeTx InvokeEmbeddedContractTx error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
		return
	}
	if rawReq.Payer != "" {
		payerAddress, err := common.AddressFromBase58(rawReq.Payer)
		if err != nil {
			log.Infof("Cli Qid:%s SigEmbededInvokeTx AddressFromBase58 error:%s", req.Qid, err)
			resp.ErrorCode = clisvrcom.CLIERR_INVALID_PARAMS
			return
		}
		mutable.Payer = payerAddress
	}
	signer := clisvrcom.DefAccount
	err = cliutil.SignTransaction(signer, mutable)
	if err != nil {
		log.Infof("Cli Qid:%s SigEmbededInvokeTx SignTransaction error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	tx, err := mutable.IntoImmutable()
	if err != nil {
		log.Infof("Cli Qid:%s SigEmbededInvokeTx mutable Serialize error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	sink := common.ZeroCopySink{}
	err = tx.Serialization(&sink)
	if err != nil {
		log.Infof("Cli Qid:%s SigEmbededInvokeTx tx Serialize error:%s", req.Qid, err)
		resp.ErrorCode = clisvrcom.CLIERR_INTERNAL_ERR
		return
	}
	resp.Result = &SigEmbededInvokeTxRsp{
		SignedTx: hex.EncodeToString(sink.Bytes()),
	}
}
