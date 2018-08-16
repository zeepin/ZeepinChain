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

package genesis

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/common/config"
	"github.com/imZhuFei/zeepin/common/constants"
	"github.com/imZhuFei/zeepin/consensus/vbft/config"
	"github.com/imZhuFei/zeepin/core/types"
	"github.com/imZhuFei/zeepin/core/utils"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/global_params"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/governance"
	nutils "github.com/imZhuFei/zeepin/smartcontract/service/native/utils"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/zpt"
	"github.com/imZhuFei/zeepin/smartcontract/service/neovm"
	"github.com/ontio/ontology-crypto/keypair"
)

const (
	BlockVersion uint32 = 0
	GenesisNonce uint64 = 2083236893
)

var (
	ZPTToken    = newGoverningToken()
	GALAToken   = newUtilityToken()
	ZPTTokenID  = ZPTToken.Hash()
	GALATokenID = GALAToken.Hash()
)

var GenBlockTime = (config.DEFAULT_GEN_BLOCK_TIME * time.Second)

var INIT_PARAM = map[string]string{
	"gasPrice": "0",
}

var GenesisBookkeepers []keypair.PublicKey

// BuildGenesisBlock returns the genesis block with default consensus bookkeeper list
func BuildGenesisBlock(defaultBookkeeper []keypair.PublicKey, genesisConfig *config.GenesisConfig) (*types.Block, error) {
	//getBookkeeper
	GenesisBookkeepers = defaultBookkeeper
	nextBookkeeper, err := types.AddressFromBookkeepers(defaultBookkeeper)
	if err != nil {
		return nil, fmt.Errorf("[Block],BuildGenesisBlock err with GetBookkeeperAddress: %s", err)
	}
	conf := bytes.NewBuffer(nil)
	if genesisConfig.GBFT != nil {
		genesisConfig.GBFT.Serialize(conf)
	}
	govConfig := newGoverConfigInit(conf.Bytes())
	consensusPayload, err := vconfig.GenesisConsensusPayload(govConfig.Hash(), 0)
	if err != nil {
		return nil, fmt.Errorf("consensus genesus init failed: %s", err)
	}
	//blockdata
	genesisHeader := &types.Header{
		Version:          BlockVersion,
		PrevBlockHash:    common.Uint256{},
		TransactionsRoot: common.Uint256{},
		Timestamp:        constants.GENESIS_BLOCK_TIMESTAMP,
		Height:           uint32(0),
		ConsensusData:    GenesisNonce,
		NextBookkeeper:   nextBookkeeper,
		ConsensusPayload: consensusPayload,

		Bookkeepers: nil,
		SigData:     nil,
	}

	//block
	zpt := newGoverningToken()
	gala := newUtilityToken()
	param := newParamContract()
	oid := deployGIDContract()
	auth := deployAuthContract()
	config := newConfig()

	genesisBlock := &types.Block{
		Header: genesisHeader,
		Transactions: []*types.Transaction{
			zpt,
			gala,
			param,
			oid,
			auth,
			config,
			newGoverningInit(),
			newUtilityInit(),
			newParamInit(),
			govConfig,
		},
	}
	genesisBlock.RebuildMerkleRoot()
	return genesisBlock, nil
}

func newGoverningToken() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.ZptContractAddress[:], "ZPT", "1.0",
		"zeepin Team", "contact@zeepin.io", "zeepin Network ZPT Token", true)
	return tx
}

func newUtilityToken() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.GalaContractAddress[:], "GALA", "1.0",
		"zeepin Team", "contact@zeepin.io", "zeepin Network GALA Token", true)
	return tx
}

func newParamContract() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.ParamContractAddress[:],
		"ParamConfig", "1.0", "zeepin Team", "contact@zeepin.io",
		"Chain Global Environment Variables Manager ", true)
	return tx
}

func newConfig() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.GovernanceContractAddress[:], "CONFIG", "1.0",
		"zeepin Team", "contact@zeepin.io", "zeepin Network Consensus Config", true)
	return tx
}

func deployAuthContract() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.AuthContractAddress[:], "AuthContract", "1.0",
		"zeepin Team", "contact@zeepin.io", "zeepin Network Authorization Contract", true)
	return tx
}

func deployGIDContract() *types.Transaction {
	tx := utils.NewDeployTransaction(nutils.GIDContractAddress[:], "GID", "1.0",
		"zeepin Team", "contact@zeepin.io", "zeepin Network GID", true)
	return tx
}

func newGoverningInit() *types.Transaction {
	bookkeepers, _ := config.DefConfig.GetBookkeepers()

	var addr common.Address
	if len(bookkeepers) == 1 {
		addr = types.AddressFromPubKey(bookkeepers[0])
	} else {
		//m := (5*len(bookkeepers) + 6) / 7
		m := int(math.Ceil(float64(len(bookkeepers)) * 2.0 / 3.0))
		temp, err := types.AddressFromMultiPubKeys(bookkeepers, m)
		if err != nil {
			panic(fmt.Sprint("wrong bookkeeper config, caused by", err))
		}
		addr = temp
	}

	distribute := []struct {
		addr  common.Address
		value uint64
	}{{addr, constants.ZPT_TOTAL_SUPPLY}}

	args := bytes.NewBuffer(nil)
	nutils.WriteVarUint(args, uint64(len(distribute)))
	for _, part := range distribute {
		nutils.WriteAddress(args, part.addr)
		nutils.WriteVarUint(args, part.value)
	}

	return utils.BuildNativeTransaction(nutils.ZptContractAddress, zpt.INIT_NAME, args.Bytes())
}

func newUtilityInit() *types.Transaction {
	bookkeepers, _ := config.DefConfig.GetBookkeepers()

	var addr common.Address
	if len(bookkeepers) == 1 {
		addr = types.AddressFromPubKey(bookkeepers[0])
	} else {
		//m := (5*len(bookkeepers) + 6) / 7
		m := int(math.Ceil(float64(len(bookkeepers)) * 2.0 / 3.0))
		temp, err := types.AddressFromMultiPubKeys(bookkeepers, m)
		if err != nil {
			panic(fmt.Sprint("wrong bookkeeper config, caused by", err))
		}
		addr = temp
	}

	distribute := []struct {
		addr  common.Address
		value uint64
	}{{addr, constants.GALA_TOTAL_SUPPLY - constants.GALA_UNBOUND_SUPPLY}}

	args := bytes.NewBuffer(nil)
	nutils.WriteVarUint(args, uint64(len(distribute)))
	for _, part := range distribute {
		nutils.WriteAddress(args, part.addr)
		nutils.WriteVarUint(args, part.value)
	}

	return utils.BuildNativeTransaction(nutils.GalaContractAddress, zpt.INIT_NAME, args.Bytes())
}

func newParamInit() *types.Transaction {
	params := new(global_params.Params)
	var s []string
	for k, _ := range INIT_PARAM {
		s = append(s, k)
	}

	neovm.GAS_TABLE.Range(func(key, value interface{}) bool {
		INIT_PARAM[key.(string)] = strconv.FormatUint(value.(uint64), 10)
		s = append(s, key.(string))
		return true
	})

	sort.Strings(s)
	for _, v := range s {
		params.SetParam(global_params.Param{Key: v, Value: INIT_PARAM[v]})
	}
	bf := new(bytes.Buffer)
	params.Serialize(bf)

	bookkeepers, _ := config.DefConfig.GetBookkeepers()
	var addr common.Address
	if len(bookkeepers) == 1 {
		addr = types.AddressFromPubKey(bookkeepers[0])
	} else {
		//m := (5*len(bookkeepers) + 6) / 7
		m := int(math.Ceil(float64(len(bookkeepers)) * 2.0 / 3.0))
		temp, err := types.AddressFromMultiPubKeys(bookkeepers, m)
		if err != nil {
			panic(fmt.Sprint("wrong bookkeeper config, caused by", err))
		}
		addr = temp
	}
	nutils.WriteAddress(bf, addr)

	return utils.BuildNativeTransaction(nutils.ParamContractAddress, global_params.INIT_NAME, bf.Bytes())
}

func newGoverConfigInit(config []byte) *types.Transaction {
	return utils.BuildNativeTransaction(nutils.GovernanceContractAddress, governance.INIT_CONFIG, config)
}
