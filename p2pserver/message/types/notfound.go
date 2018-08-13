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
	"fmt"

	"github.com/mileschao/ZeepinChain/common"
	"github.com/mileschao/ZeepinChain/errors"
	common2 "github.com/mileschao/ZeepinChain/p2pserver/common"
)

type NotFound struct {
	Hash common.Uint256
}

//Serialize message payload
func (this NotFound) Serialization() ([]byte, error) {
	return this.Hash[:], nil
}

func (this NotFound) CmdType() string {
	return common2.NOT_FOUND_TYPE
}

//Deserialize message payload
func (this *NotFound) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)

	err := this.Hash.Deserialize(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("deserialize Hash error. buf:%v", buf))
	}
	return nil
}
