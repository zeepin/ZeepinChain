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

package simulator

import (
	"bytes"
	"math/big"
	"testing"

	vtypes "github.com/imZhuFei/zeepin/embed/simulator/types"
)

func TestOpCat(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("aaa"))))
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("bbb"))))
	e.EvaluationStack = stack

	opCat(&e)
	v, err := PeekNByteArray(0, &e)
	if err != nil {
		t.Fatal("embed OpCat test failed.")
	}
	if Count(&e) != 1 || !bytes.Equal(v, []byte("aaabbb")) {
		t.Fatalf("embed OpCat test failed, expect aaabbb, got %s.", string(v))
	}
}

func TestOpSubStr(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("12345"))))
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(1))))
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(4))))
	e.EvaluationStack = stack

	opSubStr(&e)
	v, err := PeekNByteArray(0, &e)
	if err != nil {
		t.Fatal("embed OpSubStr test failed.")
	}
	if !bytes.Equal(v, []byte("2345")) {
		t.Fatalf("embed OpSubStr test failed, expect 234, got %s.", string(v))
	}
}

func TestOpLeft(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("12345"))))
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(4))))
	e.EvaluationStack = stack

	opLeft(&e)
	v, err := PeekNByteArray(0, &e)
	if err != nil {
		t.Fatal("embed OpLeft test failed.")
	}
	if !bytes.Equal(v, []byte("1234")) {
		t.Fatalf("embed OpLeft test failed, expect 1234, got %s.", string(v))
	}
}

func TestOpRight(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("12345"))))
	stack.Push(NewStackItem(vtypes.NewInteger(big.NewInt(3))))
	e.EvaluationStack = stack

	opRight(&e)
	v, err := PeekNByteArray(0, &e)
	if err != nil {
		t.Fatal("embed OpRight test failed.")
	}
	if !bytes.Equal(v, []byte("345")) {
		t.Fatalf("embed OpRight test failed, expect 345, got %s.", string(v))
	}
}

func TestOpSize(t *testing.T) {
	var e ExecutionEngine
	stack := NewRandAccessStack()
	stack.Push(NewStackItem(vtypes.NewByteArray([]byte("12345"))))
	e.EvaluationStack = stack

	opSize(&e)
	v, err := PeekInt(&e)
	if err != nil {
		t.Fatal("embed OpSize test failed.")
	}
	if v != 5 {
		t.Fatalf("embed OpSize test failed, expect 5, got %d.", v)
	}
}
