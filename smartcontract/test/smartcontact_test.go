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

package test

import (
	"os"
	"testing"

	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/core/types"
	. "github.com/imZhuFei/zeepin/smartcontract"
	"github.com/stretchr/testify/assert"
)

func TestSmartContratc(t *testing.T) {
	evilBytecode, _ := common.HexToBytes("53c56b6c766b00527ac46c766b51527ac46c766b00c303507574876438006c766b51c3c0529c63080000616c75666c766b51c300c36c766b51c351c36c766b52527ac46c766b52c3617ce001023a00616c75666c766b00c303476574876424006c766b51c3c0519c63080000616c75666c766b51c300c361e001015f00616c756600616c756652c56b6c766b00527ac46c766b51527ac461681953797374656d2e53746f726167652e476574436f6e746578746c766b00c36c766b51c3615272681253797374656d2e53746f726167652e50757451616c756651c56b6c766b00527ac461681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681253797374656d2e53746f726167652e476574616c7566")
	dbFile := "test"
	defer func() {
		os.RemoveAll(dbFile)
	}()
	//testLevelDB, err := leveldbstore.NewLevelDBStore(dbFile)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//store := statestore.NewMemDatabase()
	//testBatch := statestore.NewStateStoreBatch(store, testLevelDB)
	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}
	//cache := storage.NewCloneCache(testBatch)
	sc := SmartContract{
		Config:     config,
		Gas:        10000,
		CloneCache: nil,
	}
	engine, err := sc.NewExecuteEngine(evilBytecode)
	if err != nil {
		t.Fatal(err)
	}
	_, err = engine.Invoke()
	assert.Equal(t, "[NeoVmService] vm execute error!: the biginteger over max size 32bit", err.Error())
}
