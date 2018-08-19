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
		AdminGID:             "GID:ZPT:ZZdjHbif4HHSkyAnuT6RfyFazcmWwWc3sD",
		MinInitStake:         10000000000,
		VrfValue:             "1c9810aa9822e511d5804a9c4db9dd08497c31087b0daafa34d768a3253441fa20515e2f30f81741102af0ca3cefc4818fef16adb825fbaa8cad78647f3afb590e",
		VrfProof:             "c57741f934042cb8d8b087b44b161db56fc3ffd4ffb675d36cd09f83935be853d8729f3f5298d12d6fd28d45dde515a4b9d7f67682d182ba5118abf451ff1988",
		Peers: []*VBFTPeerStakeInfo{
			{
				Index:      1,
				PeerPubkey: "02d1b5f1aed605c6a2c3599d51293493acaeab430435b55f86f85f5c3a913c5fd2",
				Address:    "ZWVxDL3Adpu2hwHjWsKrJ8ECXcupzxHCUe",
				InitPos:    10000000000,
			},
			{
				Index:      2,
				PeerPubkey: "02e2bf0ff8c6dfd703693d93b0ea0a83cd65dcf70037028aab0f603b12d5d5b65c",
				Address:    "ZNS7q11dWyzm3Mv4aFEHxgkt4FPXgNzba7",
				InitPos:    20000000000,
			},
			{
				Index:      3,
				PeerPubkey: "026b228b66d56a1e756170a6060e5b9979ae536a23aad7b3302db686bf652de2c9",
				Address:    "ZT3t31uSamjPhGgTFxsnFFVW8vQfbaoC7q",
				InitPos:    30000000000,
			},
			{
				Index:      4,
				PeerPubkey: "020f6ea09dbefe59f09b04a4e5c7ab5e6e295ba5c160ff4fbd552abea888e021ca",
				Address:    "ZKP2FrcJKiYspWHhq6Somo3LtWxYj9mMoW",
				InitPos:    40000000000,
			},
			{
				Index:      5,
				PeerPubkey: "034410daea0782b9a65e2bab6180489c954b58492814fbad19ee6f52575ed0498e",
				Address:    "ZYAsjhq7vhPfTSYyNSJu9FdkPTC9XMrsKE",
				InitPos:    30000000000,
			},
			{
				Index:      6,
				PeerPubkey: "02955d97ae4596a2d150a35ca24937537d2bf29af99a54220be007fc66e5d1e463",
				Address:    "ZPvGsKUAatEU2qQPJPxtQu7D798Y34bHQK",
				InitPos:    20000000000,
			},
			{
				Index:      7,
				PeerPubkey: "022dc33945b90f244a47e4c001f9b07cfcba503c65283addee44c7bf61c62252c6",
				Address:    "ZKemxih7nMau99Es6crzBv3K9kiAgh9tG7",
				InitPos:    10000000000,
			},
			{
				Index:      8,
				PeerPubkey: "036a575a507aa202fd5c4b16222f84de0035eec4c91b269f10f524fc64682b89d9",
				Address:    "ZYyJdBugsx48w9mB5PisbEir7aM7QtUZxF",
				InitPos:    20000000000,
			},
			{
				Index:      9,
				PeerPubkey: "02306c7a136080b181b9a60ff849be320adb87709b4cf2687d1c970ec6d9ea404e",
				Address:    "ZKQ314rrcfNBeDfhC8uJuLts4uZbyQXFrN",
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
		AdminGID:             "GID:ZPT:ZayqjARGEFJ7gS7Ydud9CymSJPZRHW3haK",
		MinInitStake:         10000000000,
		VrfValue:             "1c9810aa9822e511d5804a9c4db9dd08497c31087b0daafa34d768a3253441fa20515e2f30f81741102af0ca3cefc4818fef16adb825fbaa8cad78647f3afb590e",
		VrfProof:             "c57741f934042cb8d8b087b44b161db56fc3ffd4ffb675d36cd09f83935be853d8729f3f5298d12d6fd28d45dde515a4b9d7f67682d182ba5118abf451ff1988",
		Peers: []*VBFTPeerStakeInfo{
			{
				Index:      1,
				PeerPubkey: "0382db58efff04b97222521317ded058097979caf9f42f2bcb63b8c52fd3a1618f",
				Address:    "ZUsDJzfJ8k9gmh1dPjWHB8fr76FrMctu7Z",
			},
			{
				Index:      2,
				PeerPubkey: "025d17e332db266117747595716f3f89933d4ed14292c4abe084cddb0493604b70",
				Address:    "ZE5ZynNQBEvJeKXa4ZgvYRrUv7LjLRHaLT",
			},
			{
				Index:      3,
				PeerPubkey: "0357cf491a64fb3a26b48affe7d7f4783911ca8b7059e650d7c2728690f503c9b3",
				Address:    "ZULDP1KcUkWz9zuYsVmH1JMrcQsxg1Qikm",
			},
			{
				Index:      4,
				PeerPubkey: "02d8c7fd2a99a688a30c296cb1dfef062a890093e48d99c4c014eb0fca06c869e0",
				Address:    "ZXbs6BsgBoG3mbcRwAXgZukEMBZuZx55sc",
			},
			{
				Index:      5,
				PeerPubkey: "03cba415e260fcd555538a248c524cb2475de610bc7a5e896a35e8f9d21428f21a",
				Address:    "ZFmyLhXz83xhsAKAjJpq6PWyVJNuVNdZFw",
			},
			{
				Index:      6,
				PeerPubkey: "039a874a39b3b4f985a4b3ccc118d479b62e79a4cde1ba223bf8b3f9ed3136720c",
				Address:    "ZQHNjEv2MNSxJLbkBno21LXpJMs8epCkGp",
			},
			{
				Index:      7,
				PeerPubkey: "030da20ac5ca163f696c4efe2b8da9131e9aad66d58a783f5996a8693b99bbf366",
				Address:    "ZFLbv9caNgbGqZ7fXKJyVat8VcGnZuuBAz",
			},
			{
				Index:      8,
				PeerPubkey: "031521d198501230a3bc0905ea75e1ba867f6c12a27604758f339b8a172cb19880",
				Address:    "ZUPn74WjnLqaa9eSXLtbSCnWV5pLQRwBL7",
			},
			{
				Index:      9,
				PeerPubkey: "036b3d9e875712bc39752f3188244c17871cc7ff5b8df1c1c9b6a57914256c6c0c",
				Address:    "ZWjDzaEYffPjnSNRc7iv7UGmxgJSFTtrDT",
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
