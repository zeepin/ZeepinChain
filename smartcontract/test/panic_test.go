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
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/zeepin/ZeepinChain/common/log"
	"github.com/zeepin/ZeepinChain/common/serialization"
	"github.com/zeepin/ZeepinChain/core/types"
	"github.com/zeepin/ZeepinChain/embed/simulator"
	. "github.com/zeepin/ZeepinChain/smartcontract"
	"github.com/zeepin/ZeepinChain/smartcontract/service/native/embed"
	"github.com/stretchr/testify/assert"
)

func TestRandomCodeCrash(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("Log")
	}()

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}

	var code []byte
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("code %x \n", code)
		}
	}()

	for i := 1; i < 10; i++ {
		fmt.Printf("test round:%d \n", i)
		code := make([]byte, i)
		for j := 0; j < 10; j++ {
			rand.Read(code)

			//cache := storage.NewCloneCache(testBatch)
			sc := SmartContract{
				Config:     config,
				Gas:        10000,
				CloneCache: nil,
			}
			engine, _ := sc.NewExecuteEngine(code)
			engine.Invoke()
		}
	}
}

func TestOpCodeDUP(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("Log")
	}()

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}

	var code = []byte{byte(simulator.DUP)}

	sc := SmartContract{
		Config:     config,
		Gas:        10000,
		CloneCache: nil,
	}
	engine, _ := sc.NewExecuteEngine(code)
	_, err := engine.Invoke()

	assert.NotNil(t, err)
}

func TestOpReadMemAttack(t *testing.T) {
	log.InitLog(4)
	defer func() {
		os.RemoveAll("Log")
	}()

	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}

	bf := new(bytes.Buffer)
	builder := simulator.NewParamsBuilder(bf)
	builder.Emit(simulator.SYSCALL)
	bs := bytes.NewBuffer(builder.ToArray())
	builder.EmitPushByteArray([]byte(embed.NATIVE_INVOKE_NAME))
	l := 0X7fffffc7 - 1
	serialization.WriteVarUint(bs, uint64(l))
	b := make([]byte, 4)
	bs.Write(b)

	sc := SmartContract{
		Config:     config,
		Gas:        100000,
		CloneCache: nil,
	}
	engine, _ := sc.NewExecuteEngine(bs.Bytes())
	_, err := engine.Invoke()

	assert.NotNil(t, err)

}
