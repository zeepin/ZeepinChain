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

package handlers

import (
	"encoding/json"
	"testing"

	clisvrcom "github.com/imZhuFei/zeepin/cmd/sigsvr/common"
	"github.com/imZhuFei/zeepin/cmd/utils"
	"github.com/imZhuFei/zeepin/common"
)

func TestSigEmbededInvokeTx(t *testing.T) {
	addr1 := common.Address([20]byte{1})
	address1 := addr1.ToHexString()
	invokeReq := &SigEmbededInvokeTxReq{
		GasPrice: 0,
		GasLimit: 0,
		Address:  address1,
		Params: []interface{}{
			&utils.EmbeddedInvokeParam{
				Type:  "string",
				Value: "foo",
			},
			&utils.EmbeddedInvokeParam{
				Type: "array",
				Value: []interface{}{
					&utils.EmbeddedInvokeParam{
						Type:  "int",
						Value: "0",
					},
					&utils.EmbeddedInvokeParam{
						Type:  "bool",
						Value: "true",
					},
				},
			},
		},
	}
	data, err := json.Marshal(invokeReq)
	if err != nil {
		t.Errorf("json.Marshal SigEmbededInvokeTxReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "sigembededinvoketx",
		Params: data,
	}
	rsp := &clisvrcom.CliRpcResponse{}
	SigEmbededInvokeTx(req, rsp)
	if rsp.ErrorCode != 0 {
		t.Errorf("SigEmbededInvokeTx failed. ErrorCode:%d ErrorInfo:%s", rsp.ErrorCode, rsp.ErrorInfo)
		return
	}
}
