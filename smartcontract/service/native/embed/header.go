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
	"github.com/imZhuFei/zeepin/errors"
)

// HeaderGetHash put header's hash to vm stack
func HeaderGetHash(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	d, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	var data *types.Header
	if b, ok := d.(*types.Block); ok {
		data = b.Header
	} else if h, ok := d.(*types.Header); ok {
		data = h
	} else {
		return errors.NewErr("[HeaderGetHash] Wrong type!")
	}
	h := data.Hash()
	vm.PushData(engine, h.ToArray())
	return nil
}

// HeaderGetVersion put header's version to vm stack
func HeaderGetVersion(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	d, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	var data *types.Header
	if b, ok := d.(*types.Block); ok {
		data = b.Header
	} else if h, ok := d.(*types.Header); ok {
		data = h
	} else {
		return errors.NewErr("[HeaderGetVersion] Wrong type!")
	}
	vm.PushData(engine, data.Version)
	return nil
}

// HeaderGetPrevHash put header's prevblockhash to vm stack
func HeaderGetPrevHash(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	d, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	var data *types.Header
	if b, ok := d.(*types.Block); ok {
		data = b.Header
	} else if h, ok := d.(*types.Header); ok {
		data = h
	} else {
		return errors.NewErr("[HeaderGetPrevHash] Wrong type!")
	}
	vm.PushData(engine, data.PrevBlockHash.ToArray())
	return nil
}

// HeaderGetMerkleRoot put header's merkleroot to vm stack
func HeaderGetMerkleRoot(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	d, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	var data *types.Header
	if b, ok := d.(*types.Block); ok {
		data = b.Header
	} else if h, ok := d.(*types.Header); ok {
		data = h
	} else {
		return errors.NewErr("[HeaderGetMerkleRoot] Wrong type!")
	}
	vm.PushData(engine, data.TransactionsRoot.ToArray())
	return nil
}

// HeaderGetIndex put header's height to vm stack
func HeaderGetIndex(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	d, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	var data *types.Header
	if b, ok := d.(*types.Block); ok {
		data = b.Header
	} else if h, ok := d.(*types.Header); ok {
		data = h
	} else {
		return errors.NewErr("[HeaderGetIndex] Wrong type!")
	}
	vm.PushData(engine, data.Height)
	return nil
}

// HeaderGetTimestamp put header's timestamp to vm stack
func HeaderGetTimestamp(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	d, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	var data *types.Header
	if b, ok := d.(*types.Block); ok {
		data = b.Header
	} else if h, ok := d.(*types.Header); ok {
		data = h
	} else {
		return errors.NewErr("[HeaderGetTimestamp] Wrong type!")
	}
	vm.PushData(engine, data.Timestamp)
	return nil
}

// HeaderGetConsensusData put header's consensus data to vm stack
func HeaderGetConsensusData(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	d, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	var data *types.Header
	if b, ok := d.(*types.Block); ok {
		data = b.Header
	} else if h, ok := d.(*types.Header); ok {
		data = h
	} else {
		return errors.NewErr("[HeaderGetConsensusData] Wrong type!")
	}
	vm.PushData(engine, data.ConsensusData)
	return nil
}

// HeaderGetNextConsensus put header's consensus to vm stack
func HeaderGetNextConsensus(service *EmbeddedService, engine *vm.ExecutionEngine) error {
	d, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	var data *types.Header
	if b, ok := d.(*types.Block); ok {
		data = b.Header
	} else if h, ok := d.(*types.Header); ok {
		data = h
	} else {
		return errors.NewErr("[HeaderGetNextConsensus] Wrong type!")
	}
	vm.PushData(engine, data.NextBookkeeper[:])
	return nil
}
