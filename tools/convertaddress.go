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

package main

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/zeepin/ZeepinChain/common"
	base58 "github.com/itchyny/base58-go"
)

const ADDR_LEN = 20

var ADDRESS_EMPTY = common.Address{}

var ADDR_PREFIX = byte(23)

func main() {
	fromAdrr, _ := ConvertAddressFromBase58("AUBek6JUBn5wMbpitRi5agfgYXiXBNFa5K")
	fmt.Println(fromAdrr.ToBase58())
}

// AddressFromBase58 returns Address from encoded base58 string
func ConvertAddressFromBase58(encoded string) (common.Address, error) {
	if encoded == "" {
		return ADDRESS_EMPTY, fmt.Errorf("invalid address")
	}
	decoded, err := base58.BitcoinEncoding.Decode([]byte(encoded))
	if err != nil {
		return ADDRESS_EMPTY, err
	}

	x, ok := new(big.Int).SetString(string(decoded), 10)
	if !ok {
		return ADDRESS_EMPTY, fmt.Errorf("invalid address")
	}

	buf := x.Bytes()
	if len(buf) != 1+ADDR_LEN+4 || buf[0] != byte(ADDR_PREFIX) {
		return ADDRESS_EMPTY, errors.New("wrong encoded address")
	}
	ph, err := common.AddressParseFromBytes(buf[1:21])
	if err != nil {
		return ADDRESS_EMPTY, err
	}

	addr := ph.ToBase58()
	fmt.Println(addr)
	return ph, nil
}
