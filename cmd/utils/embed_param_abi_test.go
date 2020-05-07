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
package utils

import (
	"fmt"
	"testing"
)

func TestParseEmbededFunc(t *testing.T) {
	var testEmbededAbi = `{
  "hash": "0xe827bf96529b5780ad0702757b8bad315e2bb8ce",
  "entrypoint": "Main",
  "functions": [
    {
      "name": "Main",
      "parameters": [
        {
          "name": "operation",
          "type": "String"
        },
        {
          "name": "args",
          "type": "Array"
        }
      ],
      "returntype": "Any"
    },
    {
      "name": "Add",
      "parameters": [
        {
          "name": "a",
          "type": "Integer"
        },
        {
          "name": "b",
          "type": "Integer"
        }
      ],
      "returntype": "Integer"
    }
  ],
  "events": []
}`
	contractAbi, err := NewEmbedContractAbi([]byte(testEmbededAbi))
	if err != nil {
		t.Errorf("TestParseEmbededFunc NewEmbedContractAbi error:%s", err)
		return
	}
	funcAbi := contractAbi.GetFunc("Add")
	if funcAbi == nil {
		t.Error("TestParseEmbededFunc cannot find func abi")
		return
	}

	params, err := ParseEmbededFunc([]string{"12", "34"}, funcAbi)
	if err != nil {
		t.Error("TestParseEmbededFunc ParseEmbededFunc error:%s", err)
		return
	}
	fmt.Printf("TestParseEmbededFunc %v\n", params)
}
