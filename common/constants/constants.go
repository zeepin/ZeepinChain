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

package constants

import (
	"time"
)

// genesis constants
var (
	//TODO: modify this when on mainnet
	GENESIS_BLOCK_TIMESTAMP = uint32(time.Date(2018, time.August, 31, 0, 0, 0, 0, time.UTC).Unix())
)

// zpt constants
const (
	ZPT_NAME         = "ZPT Token"
	ZPT_SYMBOL       = "ZPT"
	ZPT_DECIMALS     = 4
	ZPT_TOTAL_SUPPLY = uint64(10000000000000)
)

// gala constants
const (
	GALA_NAME           = "GALA Token"
	GALA_SYMBOL         = "GALA"
	GALA_DECIMALS       = 4
	GALA_TOTAL_SUPPLY   = uint64(1000000000000000)
	GALA_UNBOUND_SUPPLY = uint64(180000000000000)
)

// zpt/gala unbound model constants
const UNBOUND_TIME_INTERVAL = uint32(31536000)

var UNBOUND_GENERATION_AMOUNT = [18]uint64{89, 89, 55, 55, 55, 34, 34, 34, 21, 21, 21, 13, 13, 13, 8, 8, 5, 5}

// the end of unbound timestamp offset from genesis block's timestamp
var UNBOUND_DEADLINE = (func() uint32 {
	count := uint64(0)
	for _, m := range UNBOUND_GENERATION_AMOUNT {
		count += m
	}
	count *= uint64(UNBOUND_TIME_INTERVAL)

	numInterval := len(UNBOUND_GENERATION_AMOUNT)

	if UNBOUND_GENERATION_AMOUNT[numInterval-1] != 5 ||
		!((count-(uint64(UNBOUND_TIME_INTERVAL)*5))*10000 < GALA_UNBOUND_SUPPLY && GALA_UNBOUND_SUPPLY <= count*10000) {
		panic("incompatible constants setting")
	}

	return UNBOUND_TIME_INTERVAL*uint32(numInterval) - uint32((count-(GALA_UNBOUND_SUPPLY/10000))/5)
})()

// multi-sig constants
const MULTI_SIG_MAX_PUBKEY_SIZE = 16

// transaction constants
const TX_MAX_SIG_SIZE = 16

// network magic number
const (
	NETWORK_MAGIC_MAINNET = 0x8c77ab60
	NETWORK_MAGIC_POLARIS = 0x2d8829df
)
