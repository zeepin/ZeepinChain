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
	"crypto/sha256"
	"errors"
	"math"

	"github.com/mileschao/ZeepinChain/common"
	"github.com/mileschao/ZeepinChain/common/constants"
	"github.com/mileschao/ZeepinChain/core/program"
	"github.com/ontio/ontology-crypto/keypair"
	"golang.org/x/crypto/ripemd160"
)

func AddressFromPubKey(pubkey keypair.PublicKey) common.Address {
	prog := program.ProgramFromPubKey(pubkey)

	return AddressFromVmCode(prog)
}

func AddressFromMultiPubKeys(pubkeys []keypair.PublicKey, m int) (common.Address, error) {
	var addr common.Address
	n := len(pubkeys)
	if !(1 <= m && m <= n && n > 1 && n <= constants.MULTI_SIG_MAX_PUBKEY_SIZE) {
		return addr, errors.New("wrong multi-sig param")
	}

	prog, err := program.ProgramFromMultiPubKey(pubkeys, m)
	if err != nil {
		return addr, err
	}

	return AddressFromVmCode(prog), nil
}

func AddressFromVmCode(code []byte) common.Address {
	var addr common.Address
	temp := sha256.Sum256(code)
	md := ripemd160.New()
	md.Write(temp[:])
	md.Sum(addr[:0])

	return addr
}

func AddressFromBookkeepers(bookkeepers []keypair.PublicKey) (common.Address, error) {
	if len(bookkeepers) == 1 {
		return AddressFromPubKey(bookkeepers[0]), nil
	}
	m := int(math.Ceil(float64(len(bookkeepers)) * 2.0 / 3.0))
	return AddressFromMultiPubKeys(bookkeepers, m)
}
