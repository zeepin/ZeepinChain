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
package abi

import "strings"

const (
	EMBED_PARAM_TYPE_BOOL       = "boolean"
	EMBED_PARAM_TYPE_STRING     = "string"
	EMBED_PARAM_TYPE_INTEGER    = "integer"
	EMBED_PARAM_TYPE_ARRAY      = "array"
	EMBED_PARAM_TYPE_BYTE_ARRAY = "bytearray"
	EMBED_PARAM_TYPE_VOID       = "void"
	EMBED_PARAM_TYPE_ANY        = "any"
)

type EmbedContractAbi struct {
	Address    string                        `json:"hash"`
	EntryPoint string                        `json:"entrypoint"`
	Functions  []*EmbededContractFunctionAbi `json:"functions"`
	Events     []*EmbededContractEventAbi      `json:"events"`
}

func (this *EmbedContractAbi) GetFunc(method string) *EmbededContractFunctionAbi {
	method = strings.ToLower(method)
	for _, funcAbi := range this.Functions {
		if strings.ToLower(funcAbi.Name) == method {
			return funcAbi
		}
	}
	return nil
}

func (this *EmbedContractAbi) GetEvent(evt string) *EmbededContractEventAbi {
	evt = strings.ToLower(evt)
	for _, evtAbi := range this.Events {
		if strings.ToLower(evtAbi.Name) == evt {
			return evtAbi
		}
	}
	return nil
}

type EmbededContractFunctionAbi struct {
	Name       string                      `json:"name"`
	Parameters []*EmbededContractParamsAbi `json:"parameters"`
	ReturnType string                      `json:"returntype"`
}

type EmbededContractParamsAbi struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type EmbededContractEventAbi struct {
	Name       string                      `json:"name"`
	Parameters []*EmbededContractParamsAbi `json:"parameters"`
	ReturnType string                      `json:"returntype"`
}
