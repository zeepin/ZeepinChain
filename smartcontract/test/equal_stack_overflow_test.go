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

package test

import (
	"os"
	"testing"

	"github.com/zeepin/ZeepinChain/common/log"
	"github.com/zeepin/ZeepinChain/core/types"
	"github.com/zeepin/ZeepinChain/embed/simulator"
	. "github.com/zeepin/ZeepinChain/smartcontract"
	"github.com/stretchr/testify/assert"
)

func TestEqualStackOverflow(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("./Log")
	}()

	code := []byte{
		byte(simulator.PUSH1),    // {1}
		byte(simulator.NEWARRAY), // {[]}
		byte(simulator.DUP),      // {[],[]}
		byte(simulator.DUP),      // {[],[],[]}
		byte(simulator.PUSH0),    // {[],[],[],0}
		byte(simulator.ROT),      // {[],[],0,[]}
		byte(simulator.SETITEM),  // {[[]]}
		byte(simulator.DUP),      // {[[]],[[]]}
		byte(simulator.EQUAL),
	}

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}
	sc := SmartContract{
		Config:     config,
		Gas:        10000,
		CloneCache: nil,
	}
	engine, _ := sc.NewExecuteEngine(code)
	_, err := engine.Invoke()

	assert.Nil(t, err)
}
