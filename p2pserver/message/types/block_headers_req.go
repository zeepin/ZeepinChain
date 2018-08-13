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

	comm "github.com/mileschao/ZeepinChain/common"
	"github.com/mileschao/ZeepinChain/errors"
	"github.com/mileschao/ZeepinChain/p2pserver/common"
)

type HeadersReq struct {
	Len       uint8
	HashStart comm.Uint256
	HashEnd   comm.Uint256
}

//Serialize message payload
func (this *HeadersReq) Serialization() ([]byte, error) {
	p := new(bytes.Buffer)
	err := binary.Write(p, binary.LittleEndian, this)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("binary.Write payload error. payload:%v", this))
	}

	return p.Bytes(), nil
}

func (this *HeadersReq) CmdType() string {
	return common.GET_HEADERS_TYPE
}

//Deserialize message payload
func (this *HeadersReq) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, this)

	return err
}
