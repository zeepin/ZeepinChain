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
package states

import (
	"bytes"
	"testing"

	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/core/types"
)

func TestContract_Serialize_Deserialize(t *testing.T) {
	addr := types.AddressFromVmCode([]byte{1})

	c := &Contract{
		Version: 0,
		Address: addr,
		Method:  "init",
		Args:    []byte{2},
	}
	bf := new(bytes.Buffer)
	if err := c.Serialize(bf); err != nil {
		t.Fatalf("Contract serialize error: %v", err)
	}

	v := new(Contract)
	if err := v.Deserialize(bf); err != nil {
		t.Fatalf("Contract deserialize error: %v", err)
	}
}

func TestDeserialize(t *testing.T) {
	addr, err := common.AddressFromHexString("204495bc64f8e78bdf590a6b93a8996e66345f3f")
	if err != nil {
		t.Errorf("%s", err)
	}
	c := &Contract{
		Version: '1',
		Address: addr,
		Method:  "addStorage",
		Args:    []byte(`{"Params":[{"type":"string","value":"haha"},{"type":"string","value":"eee"}]}`),
	}
	bf := new(bytes.Buffer)
	if err := c.Serialize(bf); err != nil {
		t.Fatalf("Contract serialize error: %v", err)
	}
	t.Errorf("len%d, %+v", len(bf.Bytes()), bf.Bytes())
}
