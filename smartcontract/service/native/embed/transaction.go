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

package embed

import (
	"github.com/imZhuFei/zeepin/core/types"
	vm "github.com/imZhuFei/zeepin/embed/simulator"
	vmtypes "github.com/imZhuFei/zeepin/embed/simulator/types"
)

// GetExecutingAddress push transaction's hash to vm stack
func TransactionGetHash(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	txn, _ := vm.PopInteropInterface(engine)
	tx := txn.(*types.Transaction)
	txHash := tx.Hash()
	vm.PushData(engine, txHash.ToArray())
	return nil
}

// TransactionGetType push transaction's type to vm stack
func TransactionGetType(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	txn, _ := vm.PopInteropInterface(engine)
	tx := txn.(*types.Transaction)
	vm.PushData(engine, int(tx.TxType))
	return nil
}

// TransactionGetAttributes push transaction's attributes to vm stack
func TransactionGetAttributes(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	vm.PopInteropInterface(engine)
	attributList := make([]vmtypes.StackItems, 0)
	vm.PushData(engine, attributList)
	return nil
}
