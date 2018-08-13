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

package types

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mileschao/ZeepinChain/common"
	"github.com/mileschao/ZeepinChain/errors"
	common2 "github.com/mileschao/ZeepinChain/p2pserver/common"
)

type DataReq struct {
	DataType common.InventoryType
	Hash     common.Uint256
}

//Serialize message payload
func (this DataReq) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	err := binary.Write(p, binary.LittleEndian, &(this.DataType))
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. DataType:%v", this.DataType))
	}
	err = this.Hash.Serialize(p)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialization error. Hash:%v", this.Hash))
	}

	return p.Bytes(), nil
}

func (this *DataReq) CmdType() string {
	return common2.GET_DATA_TYPE
}

//Deserialize message payload
func (this *DataReq) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, &(this.DataType))
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read DataType error. buf:%v", buf))
	}

	err = this.Hash.Deserialize(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("deserialize Hash error. buf:%v", buf))
	}
	return nil
}
