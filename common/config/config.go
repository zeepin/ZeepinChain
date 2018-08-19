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

package config

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/common/constants"
	"github.com/imZhuFei/zeepin/common/serialization"
	"github.com/imZhuFei/zeepin/errors"
	"github.com/ontio/ontology-crypto/keypair"
)

var Version = "" //Set value when build project

const (
	DEFAULT_CONFIG_FILE_NAME = "./config.json"
	DEFAULT_WALLET_FILE_NAME = "./wallet.dat"
	MIN_GEN_BLOCK_TIME       = 2
	DEFAULT_GEN_BLOCK_TIME   = 6
	DBFT_MIN_NODE_NUM        = 4 //min node number of dbft consensus
	SOLO_MIN_NODE_NUM        = 1 //min node number of solo consensus
	VBFT_MIN_NODE_NUM        = 4 //min node number of vbft consensus

	CONSENSUS_TYPE_DBFT = "dbft"
	CONSENSUS_TYPE_SOLO = "solo"
	CONSENSUS_TYPE_VBFT = "gbft"

	DEFAULT_LOG_LEVEL                       = 1
	DEFAULT_MAX_LOG_SIZE                    = 100 //MByte
	DEFAULT_NODE_PORT                       = uint(20338)
	DEFAULT_CONSENSUS_PORT                  = uint(20339)
	DEFAULT_RPC_PORT                        = uint(20336)
	DEFAULT_RPC_LOCAL_PORT                  = uint(20337)
	DEFAULT_REST_PORT                       = uint(20334)
	DEFAULT_WS_PORT                         = uint(20335)
	DEFAULT_MAX_CONN_IN_BOUND               = uint(1024)
	DEFAULT_MAX_CONN_OUT_BOUND              = uint(1024)
	DEFAULT_MAX_CONN_IN_BOUND_FOR_SINGLE_IP = uint(16)
	DEFAULT_HTTP_INFO_PORT                  = uint(0)
	DEFAULT_MAX_TX_IN_BLOCK                 = 60000
	DEFAULT_MAX_SYNC_HEADER                 = 500
	DEFAULT_ENABLE_CONSENSUS                = true
	DEFAULT_ENABLE_EVENT_LOG                = true
	DEFAULT_CLI_RPC_PORT                    = uint(20000)
	DEFAULT_GAS_LIMIT                       = 20000
	DEFAULT_GAS_PRICE                       = 1

	DEFAULT_DATA_DIR      = "./Chain"
	DEFAULT_RESERVED_FILE = "./peers.rsv"
)

const (
	NETWORK_ID_MAIN_NET      = 1
	NETWORK_ID_POLARIS_NET   = 2
	NETWORK_ID_SOLO_NET      = 3
	NETWORK_NAME_MAIN_NET    = "zeepin"
	NETWORK_NAME_POLARIS_NET = "polaris"
	NETWORK_NAME_SOLO_NET    = "testmode"
)

var NETWORK_MAGIC = map[uint32]uint32{
	NETWORK_ID_MAIN_NET:    constants.NETWORK_MAGIC_MAINNET, //Network main
	NETWORK_ID_POLARIS_NET: constants.NETWORK_MAGIC_POLARIS, //Network polaris
	NETWORK_ID_SOLO_NET:    0,                               //Network solo
}

var NETWORK_NAME = map[uint32]string{
	NETWORK_ID_MAIN_NET:    NETWORK_NAME_MAIN_NET,
	NETWORK_ID_POLARIS_NET: NETWORK_NAME_POLARIS_NET,
	NETWORK_ID_SOLO_NET:    NETWORK_NAME_SOLO_NET,
}

func GetNetworkMagic(id uint32) uint32 {
	nid, ok := NETWORK_MAGIC[id]
	if ok {
		return nid
	}
	return id
}

func GetNetworkName(id uint32) string {
	name, ok := NETWORK_NAME[id]
	if ok {
		return name
	}
	return fmt.Sprintf("%d", id)
}

var PolarisConfig = &GenesisConfig{
	SeedList: []string{
		"test1.zeepin.net:20338",
		"test2.zeepin.net:20338",
		"test3.zeepin.net:20338",
		"test4.zeepin.net:20338"},
	ConsensusType: CONSENSUS_TYPE_VBFT,
	GBFT: &VBFTConfig{
		N:                    9,
		C:                    3,
		K:                    9,
		L:                    144,
		BlockMsgDelay:        10000,
		HashMsgDelay:         10000,
		PeerHandshakeTimeout: 10,
		MaxBlockChangeView:   3000,
		AdminGID:             "GID:ZPT:ZPN5huXL7NkBwuBZWfksumeMxDhHkXShB4",
		MinInitStake:         10000000000,
		VrfValue:             "1c9810aa9822e511d5804a9c4db9dd08497c31087b0daafa34d768a3253441fa20515e2f30f81741102af0ca3cefc4818fef16adb825fbaa8cad78647f3afb590e",
		VrfProof:             "c57741f934042cb8d8b087b44b161db56fc3ffd4ffb675d36cd09f83935be853d8729f3f5298d12d6fd28d45dde515a4b9d7f67682d182ba5118abf451ff1988",
		Peers: []*VBFTPeerStakeInfo{
			{
				Index:      1,
				PeerPubkey: "0282ee0668c461823bb569aa58091ff8e96ff940246a8aafdce2e28f417903261e",
				Address:    "ZWbEMfyTkJf5tH4rE8daT5F9GvkSYwguHd",
				InitPos:    10000000000,
			},
			{
				Index:      2,
				PeerPubkey: "033c3b04d435b936e117ff8c73d251a9b40daef78245af30eba058b7cfc260e77b",
				Address:    "ZR8BGPm7XRAh8GZseGVZar1WyQgARxHB28",
				InitPos:    20000000000,
			},
			{
				Index:      3,
				PeerPubkey: "0290a3b090fc02c3a9ba38b7bc2b57c4605807961cb19c48bbd2413323a7a25017",
				Address:    "ZDSseFDdBA25JjfJvQPfYbhGTuGhsNJmcb",
				InitPos:    30000000000,
			},
			{
				Index:      4,
				PeerPubkey: "03a19b0f8d94a1f12611c03dbb037ecb5be46ee5f86d3d554c098967c0fba5b28e",
				Address:    "ZG4aVjPrEX4RrXfsyPLFjk8TcpEvannekX",
				InitPos:    40000000000,
			},
			{
				Index:      5,
				PeerPubkey: "0386e6bf4e00e96032c0b993a4638aaf3ee596196c0d7b1ffe4db9c29bcef23b96",
				Address:    "ZSEkQoqsP3zma6KqPnPGkJpA71VSPgdw54",
				InitPos:    30000000000,
			},
			{
				Index:      6,
				PeerPubkey: "02f748c70fe96c5e7feb1bf477943368443d21c44141184515375f53d21ccea091",
				Address:    "ZGemTxNiQ5arMHtDtReNUnwuUvrv3XA5Go",
				InitPos:    20000000000,
			},
			{
				Index:      7,
				PeerPubkey: "03aedf7e3f4049abeae005f2438e54c9d432aa280d3c8877152f4733bd04e1ed19",
				Address:    "ZKyCKN4adbXjKMi8Gah71j9C6aPDEu16VF",
				InitPos:    10000000000,
			},
			{
				Index:      8,
				PeerPubkey: "03726e7679d59daf013c4b509a477098a4893fa75e4e3fe8c89d18a34a48770952",
				Address:    "ZZMSKtiuH7uzyF1FwjNDu1FqtARYfoc6f3",
				InitPos:    20000000000,
			},
			{
				Index:      9,
				PeerPubkey: "0233b8deedd64705a4282b3198aa9411665815e4ae991e6b5e95426e06c0717e9f",
				Address:    "ZSbCwdP1UVatSzjV343F8Xih2qSyge9pCW",
				InitPos:    20000000000,
			},
		},
	},
	DBFT: &DBFTConfig{},
	SOLO: &SOLOConfig{},
}

var MainNetConfig = &GenesisConfig{
	SeedList: []string{
		"main1.zeepin.net:20338",
		"main2.zeepin.net:20338",
		"main3.zeepin.net:20338",
		"main4.zeepin.net:20338",
		"main5.zeepin.net:20338"},
	ConsensusType: CONSENSUS_TYPE_VBFT,
	GBFT: &VBFTConfig{
		N:                    9,
		C:                    3,
		K:                    9,
		L:                    144,
		BlockMsgDelay:        10000,
		HashMsgDelay:         10000,
		PeerHandshakeTimeout: 10,
		MaxBlockChangeView:   120000,
		AdminGID:             "GID:ZPT:ZVbasGQ1jRNPmERADj2KrA95kp7aYFnUAQ",
		MinInitStake:         10000000000,
		VrfValue:             "1c9810aa9822e511d5804a9c4db9dd08497c31087b0daafa34d768a3253441fa20515e2f30f81741102af0ca3cefc4818fef16adb825fbaa8cad78647f3afb590e",
		VrfProof:             "c57741f934042cb8d8b087b44b161db56fc3ffd4ffb675d36cd09f83935be853d8729f3f5298d12d6fd28d45dde515a4b9d7f67682d182ba5118abf451ff1988",
		Peers: []*VBFTPeerStakeInfo{
			{
				Index:      1,
				PeerPubkey: "03b3da1cb21ade9278d602b9cf5f13ffa81ec4fc7c1474a3ae60b4c882b1bad54c",
				Address:    "ZJGp9ZNe6qzKJhPN7TMgH4b2hRWNQ6SHRA",
			},
			{
				Index:      2,
				PeerPubkey: "034cf74659ff023d9d5a71b60c73b423403e2b001f29237ba3e7f00a260de91fdb",
				Address:    "ZP2MyRmZMxWPbUs6JLBn4zjywSutpdAraw",
			},
			{
				Index:      3,
				PeerPubkey: "03e392fbb0ac8faf262584f594657059070078a21e4f21f6cc8602a65254dd71e7",
				Address:    "ZanCZ8srSFFs54s56WAFyPTKUg9L3WcnQ8",
			},
			{
				Index:      4,
				PeerPubkey: "022313331552c569c6c7922fe0d34574bb4269bd0899f4e36ba0272157ce6869d7",
				Address:    "ZSgP8YQ2nKiAqfy6bdo1aozEm1HrkjN2m7",
			},
			{
				Index:      5,
				PeerPubkey: "02c7e2a9da20f28a57c68d91785c3970eb4ed776c2d50ada5c0c574e9be412a881",
				Address:    "ZP1hFo6GHJTRMYNWADdCkv2xrFt6sKZy7p",
			},
			{
				Index:      6,
				PeerPubkey: "027f3a24262dc3e726b0fc6940dc6b6a396129717a869a0f8a3190a10131693e9f",
				Address:    "ZMFTVYpfq2wpHJ4PxBhPSKkvx8yuMSWC4E",
			},
			{
				Index:      7,
				PeerPubkey: "02358d9ceec30ad4d66757770ea2b0955e61d41f56939510cf75de45183c95a5af",
				Address:    "ZR8qTxRvxDumKAqeRcV8J6G3bdpscUrXNa",
			},
			{
				Index:      8,
				PeerPubkey: "02632c7f74d81be26e25b5ffccd24aefa72c7244170a5ede7935648d1a5e2d3941",
				Address:    "ZK6sUYqYJUihymSFmknwDZwAjqXezXGaAh",
			},
			{
				Index:      9,
				PeerPubkey: "0340529e1019658c4e76fadfc4d1660e512e0ee82126db9cec9bacd521999231a3",
				Address:    "ZQxXkGvc9r1VPtsjfq58dAJH356SPZLu2q",
			},
		},
	},
	DBFT: &DBFTConfig{},
	SOLO: &SOLOConfig{},
}

var DefConfig = NewZeepinChainConfig()

type GenesisConfig struct {
	SeedList      []string
	ConsensusType string
	GBFT          *VBFTConfig
	DBFT          *DBFTConfig
	SOLO          *SOLOConfig
}

func NewGenesisConfig() *GenesisConfig {
	return &GenesisConfig{
		SeedList:      make([]string, 0),
		ConsensusType: CONSENSUS_TYPE_DBFT,
		GBFT:          &VBFTConfig{},
		DBFT:          &DBFTConfig{},
		SOLO:          &SOLOConfig{},
	}
}

//
// VBFT genesis config, from local config file
//
type VBFTConfig struct {
	N                    uint32               `json:"n"` // network size
	C                    uint32               `json:"c"` // consensus quorum
	K                    uint32               `json:"k"`
	L                    uint32               `json:"l"`
	BlockMsgDelay        uint32               `json:"block_msg_delay"`
	HashMsgDelay         uint32               `json:"hash_msg_delay"`
	PeerHandshakeTimeout uint32               `json:"peer_handshake_timeout"`
	MaxBlockChangeView   uint32               `json:"max_block_change_view"`
	MinInitStake         uint64               `json:"min_init_stake"`
	AdminGID             string               `json:"admin_gid"`
	VrfValue             string               `json:"vrf_value"`
	VrfProof             string               `json:"vrf_proof"`
	Peers                []*VBFTPeerStakeInfo `json:"peers"`
}

func (this *VBFTConfig) Serialize(w io.Writer) error {
	if err := serialization.WriteUint32(w, this.N); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize n error!")
	}
	if err := serialization.WriteUint32(w, this.C); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize c error!")
	}
	if err := serialization.WriteUint32(w, this.K); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize k error!")
	}
	if err := serialization.WriteUint32(w, this.L); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize l error!")
	}
	if err := serialization.WriteUint32(w, this.BlockMsgDelay); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize block_msg_delay error!")
	}
	if err := serialization.WriteUint32(w, this.HashMsgDelay); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize hash_msg_delay error!")
	}
	if err := serialization.WriteUint32(w, this.PeerHandshakeTimeout); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize peer_handshake_timeout error!")
	}
	if err := serialization.WriteUint32(w, this.MaxBlockChangeView); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize max_block_change_view error!")
	}
	if err := serialization.WriteUint64(w, this.MinInitStake); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize min_init_stake error!")
	}
	if err := serialization.WriteString(w, this.AdminGID); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteString, serialize admin_gid error!")
	}
	if err := serialization.WriteString(w, this.VrfValue); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteString, serialize vrf_value error!")
	}
	if err := serialization.WriteString(w, this.VrfProof); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteString, serialize vrf_proof error!")
	}
	if err := serialization.WriteVarUint(w, uint64(len(this.Peers))); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteVarUint, serialize peer length error!")
	}
	for _, peer := range this.Peers {
		if err := peer.Serialize(w); err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "serialize peer error!")
		}
	}
	return nil
}

func (this *VBFTConfig) Deserialize(r io.Reader) error {
	n, err := serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize n error!")
	}
	c, err := serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize c error!")
	}
	k, err := serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize k error!")
	}
	l, err := serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize l error!")
	}
	blockMsgDelay, err := serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize blockMsgDelay error!")
	}
	hashMsgDelay, err := serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize hashMsgDelay error!")
	}
	peerHandshakeTimeout, err := serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize peerHandshakeTimeout error!")
	}
	maxBlockChangeView, err := serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize maxBlockChangeView error!")
	}
	minInitStake, err := serialization.ReadUint64(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize minInitStake error!")
	}
	adminGID, err := serialization.ReadString(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadString, deserialize adminGID error!")
	}
	vrfValue, err := serialization.ReadString(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadString, deserialize vrfValue error!")
	}
	vrfProof, err := serialization.ReadString(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadString, deserialize vrfProof error!")
	}
	length, err := serialization.ReadVarUint(r, 0)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadVarUint, deserialize peer length error!")
	}
	peers := make([]*VBFTPeerStakeInfo, 0)
	for i := 0; uint64(i) < length; i++ {
		peer := new(VBFTPeerStakeInfo)
		err = peer.Deserialize(r)
		if err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "deserialize peer error!")
		}
		peers = append(peers, peer)
	}
	this.N = n
	this.C = c
	this.K = k
	this.L = l
	this.BlockMsgDelay = blockMsgDelay
	this.HashMsgDelay = hashMsgDelay
	this.PeerHandshakeTimeout = peerHandshakeTimeout
	this.MaxBlockChangeView = maxBlockChangeView
	this.MinInitStake = minInitStake
	this.AdminGID = adminGID
	this.VrfValue = vrfValue
	this.VrfProof = vrfProof
	this.Peers = peers
	return nil
}

type VBFTPeerStakeInfo struct {
	Index      uint32 `json:"index"`
	PeerPubkey string `json:"peerPubkey"`
	Address    string `json:"address"`
	InitPos    uint64 `json:"initPos"`
}

func (this *VBFTPeerStakeInfo) Serialize(w io.Writer) error {
	if err := serialization.WriteUint32(w, this.Index); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize index error!")
	}
	if err := serialization.WriteString(w, this.PeerPubkey); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize peerPubkey error!")
	}
	address, err := common.AddressFromBase58(this.Address)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "common.AddressFromBase58, address format error!")
	}
	if err := address.Serialize(w); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize address error!")
	}
	if err := serialization.WriteUint64(w, this.InitPos); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.WriteUint32, serialize initPos error!")
	}
	return nil
}

func (this *VBFTPeerStakeInfo) Deserialize(r io.Reader) error {
	index, err := serialization.ReadUint32(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize index error!")
	}
	peerPubkey, err := serialization.ReadString(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize peerPubkey error!")
	}
	address := new(common.Address)
	err = address.Deserialize(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "address.Deserialize, deserialize address error!")
	}
	initPos, err := serialization.ReadUint64(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint32, deserialize initPos error!")
	}
	this.Index = index
	this.PeerPubkey = peerPubkey
	this.Address = address.ToBase58()
	this.InitPos = initPos
	return nil
}

type DBFTConfig struct {
	GenBlockTime uint
	Bookkeepers  []string
}

type SOLOConfig struct {
	GenBlockTime uint
	Bookkeepers  []string
}

type CommonConfig struct {
	LogLevel       uint
	NodeType       string
	EnableEventLog bool
	SystemFee      map[string]int64
	GasLimit       uint64
	GasPrice       uint64
	DataDir        string
}

type ConsensusConfig struct {
	EnableConsensus bool
	MaxTxInBlock    uint
}

type P2PRsvConfig struct {
	ReservedPeers []string `json:"reserved"`
	MaskPeers     []string `json:"mask"`
}

type P2PNodeConfig struct {
	ReservedPeersOnly         bool
	ReservedCfg               *P2PRsvConfig
	NetworkMagic              uint32
	NetworkId                 uint32
	NetworkName               string
	NodePort                  uint
	NodeConsensusPort         uint
	DualPortSupport           bool
	IsTLS                     bool
	CertPath                  string
	KeyPath                   string
	CAPath                    string
	HttpInfoPort              uint
	MaxHdrSyncReqs            uint
	MaxConnInBound            uint
	MaxConnOutBound           uint
	MaxConnInBoundForSingleIP uint
}

type RpcConfig struct {
	EnableHttpJsonRpc bool
	HttpJsonPort      uint
	HttpLocalPort     uint
}

type RestfulConfig struct {
	EnableHttpRestful bool
	HttpRestPort      uint
	HttpCertPath      string
	HttpKeyPath       string
}

type WebSocketConfig struct {
	EnableHttpWs bool
	HttpWsPort   uint
	HttpCertPath string
	HttpKeyPath  string
}

type ZeepinChainConfig struct {
	Genesis   *GenesisConfig
	Common    *CommonConfig
	Consensus *ConsensusConfig
	P2PNode   *P2PNodeConfig
	Rpc       *RpcConfig
	Restful   *RestfulConfig
	Ws        *WebSocketConfig
}

func NewZeepinChainConfig() *ZeepinChainConfig {
	return &ZeepinChainConfig{
		Genesis: MainNetConfig,
		Common: &CommonConfig{
			LogLevel:       DEFAULT_LOG_LEVEL,
			EnableEventLog: DEFAULT_ENABLE_EVENT_LOG,
			SystemFee:      make(map[string]int64),
			GasLimit:       DEFAULT_GAS_LIMIT,
			DataDir:        DEFAULT_DATA_DIR,
		},
		Consensus: &ConsensusConfig{
			EnableConsensus: true,
			MaxTxInBlock:    DEFAULT_MAX_TX_IN_BLOCK,
		},
		P2PNode: &P2PNodeConfig{
			ReservedCfg:               &P2PRsvConfig{},
			ReservedPeersOnly:         false,
			NetworkId:                 NETWORK_ID_MAIN_NET,
			NetworkName:               GetNetworkName(NETWORK_ID_MAIN_NET),
			NetworkMagic:              GetNetworkMagic(NETWORK_ID_MAIN_NET),
			NodePort:                  DEFAULT_NODE_PORT,
			NodeConsensusPort:         DEFAULT_CONSENSUS_PORT,
			DualPortSupport:           true,
			IsTLS:                     false,
			CertPath:                  "",
			KeyPath:                   "",
			CAPath:                    "",
			HttpInfoPort:              DEFAULT_HTTP_INFO_PORT,
			MaxHdrSyncReqs:            DEFAULT_MAX_SYNC_HEADER,
			MaxConnInBound:            DEFAULT_MAX_CONN_IN_BOUND,
			MaxConnOutBound:           DEFAULT_MAX_CONN_OUT_BOUND,
			MaxConnInBoundForSingleIP: DEFAULT_MAX_CONN_IN_BOUND_FOR_SINGLE_IP,
		},
		Rpc: &RpcConfig{
			EnableHttpJsonRpc: true,
			HttpJsonPort:      DEFAULT_RPC_PORT,
			HttpLocalPort:     DEFAULT_RPC_LOCAL_PORT,
		},
		Restful: &RestfulConfig{
			EnableHttpRestful: true,
			HttpRestPort:      DEFAULT_REST_PORT,
		},
		Ws: &WebSocketConfig{
			EnableHttpWs: true,
			HttpWsPort:   DEFAULT_WS_PORT,
		},
	}
}

func (this *ZeepinChainConfig) GetBookkeepers() ([]keypair.PublicKey, error) {
	var bookKeepers []string
	switch this.Genesis.ConsensusType {
	case CONSENSUS_TYPE_VBFT:
		for _, peer := range this.Genesis.GBFT.Peers {
			bookKeepers = append(bookKeepers, peer.PeerPubkey)
		}
	case CONSENSUS_TYPE_DBFT:
		bookKeepers = this.Genesis.DBFT.Bookkeepers
	case CONSENSUS_TYPE_SOLO:
		bookKeepers = this.Genesis.SOLO.Bookkeepers
	default:
		return nil, fmt.Errorf("Does not support %s consensus", this.Genesis.ConsensusType)
	}

	pubKeys := make([]keypair.PublicKey, 0, len(bookKeepers))
	for _, key := range bookKeepers {
		pubKey, err := hex.DecodeString(key)
		k, err := keypair.DeserializePublicKey(pubKey)
		if err != nil {
			return nil, fmt.Errorf("Incorrectly book keepers key:%s", key)
		}
		pubKeys = append(pubKeys, k)
	}
	keypair.SortPublicKeys(pubKeys)
	return pubKeys, nil
}

func (this *ZeepinChainConfig) GetDefaultNetworkId() (uint32, error) {
	defaultNetworkId, err := this.getDefNetworkIDFromGenesisConfig(this.Genesis)
	if err != nil {
		return 0, err
	}
	mainNetId, err := this.getDefNetworkIDFromGenesisConfig(MainNetConfig)
	if err != nil {
		return 0, err
	}
	polaridId, err := this.getDefNetworkIDFromGenesisConfig(PolarisConfig)
	if err != nil {
		return 0, err
	}
	switch defaultNetworkId {
	case mainNetId:
		return NETWORK_ID_MAIN_NET, nil
	case polaridId:
		return NETWORK_ID_POLARIS_NET, nil
	}
	return defaultNetworkId, nil
}

func (this *ZeepinChainConfig) getDefNetworkIDFromGenesisConfig(genCfg *GenesisConfig) (uint32, error) {
	var configData []byte
	var err error
	switch this.Genesis.ConsensusType {
	case CONSENSUS_TYPE_VBFT:
		configData, err = json.Marshal(genCfg.GBFT)
	case CONSENSUS_TYPE_DBFT:
		configData, err = json.Marshal(genCfg.DBFT)
	case CONSENSUS_TYPE_SOLO:
		return NETWORK_ID_SOLO_NET, nil
	default:
		return 0, fmt.Errorf("unknown consensus type:%s", this.Genesis.ConsensusType)
	}
	if err != nil {
		return 0, fmt.Errorf("json.Marshal error:%s", err)
	}
	data := sha256.Sum256(configData)
	return binary.LittleEndian.Uint32(data[0:4]), nil
}
