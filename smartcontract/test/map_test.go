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
	"fmt"
	"testing"

	"github.com/imZhuFei/zeepin/core/types"
	"github.com/imZhuFei/zeepin/embed/simulator"
	"github.com/imZhuFei/zeepin/smartcontract"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	byteCode := []byte{
		byte(simulator.NEWMAP),
		byte(simulator.DUP),   // dup map
		byte(simulator.PUSH0), // key (index)
		byte(simulator.PUSH0), // key (index)
		byte(simulator.SETITEM),

		byte(simulator.DUP),   // dup map
		byte(simulator.PUSH0), // key (index)
		byte(simulator.PUSH1), // value (newItem)
		byte(simulator.SETITEM),
	}

	// pick a value out
	byteCode = append(byteCode,
		[]byte{ // extract element
			byte(simulator.DUP),   // dup map (items)
			byte(simulator.PUSH0), // key (index)

			byte(simulator.PICKITEM),
			byte(simulator.JMPIF), // dup map (items)
			0x04, 0x00,            // skip a drop?
			byte(simulator.DROP),
		}...)

	// count faults vs successful executions
	N := 1024
	faults := 0

	//dbFile := "/tmp/test"
	//os.RemoveAll(dbFile)
	//testLevelDB, err := leveldbstore.NewLevelDBStore(dbFile)
	//if err != nil {
	//	panic(err)
	//}

	for n := 0; n < N; n++ {
		// Setup Execution Environment
		//store := statestore.NewMemDatabase()
		//testBatch := statestore.NewStateStoreBatch(store, testLevelDB)
		config := &smartcontract.Config{
			Time:   10,
			Height: 10,
			Tx:     &types.Transaction{},
		}
		//cache := storage.NewCloneCache(testBatch)
		sc := smartcontract.SmartContract{
			Config:     config,
			Gas:        100,
			CloneCache: nil,
		}
		engine, err := sc.NewExecuteEngine(byteCode)

		_, err = engine.Invoke()
		if err != nil {
			fmt.Println("err:", err)
			faults += 1
		}
	}
	assert.Equal(t, faults, 0)

}
