package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/imZhuFei/zeepin/account"
	"github.com/imZhuFei/zeepin/cmd/utils"
	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/common/constants"
	"github.com/imZhuFei/zeepin/common/password"
	"github.com/imZhuFei/zeepin/config"
	"github.com/imZhuFei/zeepin/consensus/vbft"
	"github.com/imZhuFei/zeepin/consensus/vbft/config"
	"github.com/imZhuFei/zeepin/core/payload"
	"github.com/imZhuFei/zeepin/core/signature"
	"github.com/imZhuFei/zeepin/core/types"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/auth"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/governance"
	nutils "github.com/imZhuFei/zeepin/smartcontract/service/native/utils"
	"github.com/imZhuFei/zeepin/smartcontract/service/native/zpt"
	svrneovm "github.com/imZhuFei/zeepin/smartcontract/service/neovm"
	"github.com/imZhuFei/zeepin/vm/neovm"
	"github.com/ontio/ontology-crypto/keypair"
	sig "github.com/ontio/ontology-crypto/signature"
)

func main() {
	var (
		asset        string
		from         string
		to           string
		branch       string
		walletPath   string
		publicKey    string
		value        string
		gasPrice     uint64
		gasLimit     uint64
		sts          []*zpt.State
		tx           *types.Transaction
		pubKeys      []keypair.PublicKey
		users        []*account.Account
		contractAddr common.Address
		index        int
		accountNum   int
	)
	config.Init()
	flag.StringVar(&branch, "branch", "", "branch")
	flag.StringVar(&asset, "asset", "zpt", "asset address")
	flag.StringVar(&from, "from", "", "transfer sender base58 address")
	flag.StringVar(&to, "to", "", "transfer revicer base58 address")
	flag.StringVar(&publicKey, "pk", "", "register candidate publickey")
	flag.StringVar(&walletPath, "wp", "", "register ontid address")
	flag.StringVar(&value, "value", "0", "transfer amount")
	flag.Uint64Var(&gasPrice, "gasPrice", 0, "gas price")
	flag.Uint64Var(&gasLimit, "gasLimit", 20000, "gas limit")
	flag.IntVar(&index, "index", 1, "account index")
	flag.IntVar(&accountNum, "account_num", 1, "wallet account number")
	flag.Parse()

	/*if len(config.Configuration.Wallets) == 1 && len(config.Configuration.Passwords) == 1 {
		wallet, err := account.Open(config.Configuration.Wallets[0])
		if err != nil {
			fmt.Println("open wallet " + config.Configuration.Wallets[0] + " fail.")
			return
		}
		for i := 1; i <= accountNum; i++ {
			user, err := wallet.GetAccountByIndex(i, []byte(config.Configuration.Passwords[0]))
			//user, err := wallet.GetDefaultAccount([]byte(config.Configuration.Passwords[i]))
			if err != nil {
				fmt.Println("open wallet " + config.Configuration.Wallets[i] + " password error." + err.Error())
				return
			}
			users = append(users, user)
			pubKeys = append(pubKeys, user.PublicKey)
		}
	} else {*/
	for i := 0; i < len(config.Configuration.Wallets); i++ {
		wallet, err := account.Open(config.Configuration.Wallets[i])
		if err != nil {
			fmt.Println("open wallet " + config.Configuration.Wallets[i] + " fail.")
			return
		}
		user, err := wallet.GetDefaultAccount([]byte(config.Configuration.Passwords[i]))
		if err != nil {
			fmt.Println("open wallet " + config.Configuration.Wallets[i] + " password error. " + err.Error())
			return
		}
		//fmt.Println("open wallet " + config.Configuration.Wallets[i])
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}
	//}

	switch branch {
	case "multiAddr":
		//int((5*len(pubKeys)+6)/7)
		mt := int(math.Ceil(float64(len(config.Configuration.Wallets)) * 2.0 / 3.0))
		address, err := types.AddressFromMultiPubKeys(pubKeys, mt)
		if err != nil {
			fmt.Println(fmt.Errorf("return config multi address error:%s", err))
			return
		}
		fmt.Println("config multi address:", address.ToBase58())
		return
	case "asset":
		fromAdrr, err := common.AddressFromBase58(from)
		if err != nil {
			fmt.Println("invalid from address:", from)
			return
		}
		toAdrr, err := common.AddressFromBase58(to)
		if err != nil {
			fmt.Println("invalid to address:", to)
			return
		}

		var amount uint64
		switch strings.ToLower(asset) {
		case "zpt":
			contractAddr = nutils.ZptContractAddress
			amount = utils.ParseZpt(value)
		case "gala":
			contractAddr = nutils.GalaContractAddress
			amount = utils.ParseZpt(value)
		default:
			panic(fmt.Sprintf("Unsupport asset:%s", asset))
		}
		sts = append(sts, &zpt.State{
			From:  fromAdrr,
			To:    toAdrr,
			Value: amount,
		})
		code, err := BuildNativeInvokeCode(contractAddr, 0, "transfer", []interface{}{sts})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, code)
		for i := 0; i < 9; i++ {
			if err := MultiSignToTransaction(tx, uint16(math.Ceil(float64(len(config.Configuration.Wallets))*2.0/3.0)), pubKeys, users[i]); err != nil {
				fmt.Println("sign transaction error!")
				return
			}
		}
	case "unboundgala":
		senderAddr, err := common.AddressFromBase58(from)
		if err != nil {
			return
		}
		fromAddr, err := common.AddressFromBase58(nutils.ZptContractAddress.ToBase58())
		if err != nil {
			return
		}
		toAddr, err := common.AddressFromBase58(to)
		if err != nil {
			return
		}
		balance, err := utils.GetAllowance("gala", nutils.ZptContractAddress.ToBase58(), to)
		if err != nil {
			fmt.Println("get unboundgala error:", to)
			return
		}

		amount, err := strconv.ParseUint(balance, 10, 64)
		if err != nil {
			fmt.Println("parse amount error ", balance)
			return
		}
		if amount <= 0 {
			fmt.Println("unboundgala amount <=0 ")
			return
		}
		transferFrom := &zpt.TransferFrom{
			Sender: senderAddr,
			From:   fromAddr,
			To:     toAddr,
			Value:  amount,
		}
		var version byte
		var contractAddr common.Address
		switch strings.ToLower(asset) {
		case "gala":
			version = byte(0)
			contractAddr = nutils.GalaContractAddress
		default:
			panic(fmt.Sprintf("Unsupport asset:%s", asset))
		}
		invokeCode, err := BuildNativeInvokeCode(contractAddr, version, "transferFrom", []interface{}{transferFrom})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, invokeCode)
		for i := 0; i < len(config.Configuration.Wallets); i++ {
			if err := MultiSignToTransaction(tx, uint16(math.Ceil(float64(len(config.Configuration.Wallets))*2.0/3.0)), pubKeys, users[i]); err != nil {
				fmt.Println("sign transaction error!")
				return
			}
		}
	case "RegGId":
		wallet, err := account.Open(walletPath)
		if err != nil {
			fmt.Println("open register Gid wallet "+walletPath+" fail:", err)
			return
		}
		pwd, err := password.GetPassword()
		if err != nil {
			fmt.Println("get password fail")
			return
		}
		user, err := wallet.GetAccountByIndex(index, pwd)
		if err != nil {
			fmt.Println("open wallet " + walletPath + " password error.")
			return
		}
		params := RegIDWithPublicKeyParam{
			GID:    []byte("GID:ZPT:" + user.Address.ToBase58()),
			Pubkey: keypair.SerializePublicKey(user.PublicKey),
		}
		contractAddress := nutils.GIDContractAddress
		code, err := BuildNativeInvokeCode(contractAddress, 0, "regIDWithPublicKey", []interface{}{params})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, code)
		if err := SignTransaction(user, tx); err != nil {
			fmt.Println("sign transaction error:", err)
			return
		}
	case "assignRole":
		wallet, err := account.Open(walletPath)
		if err != nil {
			fmt.Println("open assign func to role:"+walletPath+" fail:", err)
			return
		}
		pwd, err := password.GetPassword()
		if err != nil {
			fmt.Println("get password fail")
			return
		}
		user, err := wallet.GetAccountByIndex(index, pwd)
		if err != nil {
			fmt.Println("open wallet " + walletPath + " password error.")
			return
		}
		params := &auth.FuncsToRoleParam{
			ContractAddr: nutils.GovernanceContractAddress,
			AdminGID:     []byte("GID:ZPT:" + user.Address.ToBase58()),
			Role:         []byte("TrionesCandidatePeerOwner"),
			FuncNames:    []string{"registerCandidate"},
			KeyNo:        1,
		}
		contractAddress := nutils.AuthContractAddress
		code, err := BuildNativeInvokeCode(contractAddress, 0, "assignFuncsToRole", []interface{}{params})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, code)
		if err := SignTransaction(user, tx); err != nil {
			fmt.Println("sign transaction error:", err)
			return
		}
	case "assignGId":
		wallet, err := account.Open(walletPath)
		if err != nil {
			fmt.Println("open assign func to role:"+walletPath+" fail:", err)
			return
		}
		pwd, err := password.GetPassword()
		if err != nil {
			fmt.Println("get password fail")
			return
		}
		user, err := wallet.GetAccountByIndex(index, pwd)
		if err != nil {
			fmt.Println("open wallet " + walletPath + " password error.")
			return
		}
		params := &auth.GIDsToRoleParam{
			ContractAddr: nutils.GovernanceContractAddress,
			AdminGID:     []byte("GID:ZPT:" + user.Address.ToBase58()),
			Role:         []byte("TrionesCandidatePeerOwner"),
			Persons:      [][]byte{},
			KeyNo:        1,
		}
		for _, gid := range config.Configuration.GIds {
			params.Persons = append(params.Persons, []byte(gid))
		}
		code, err := BuildNativeInvokeCode(nutils.AuthContractAddress, 0, "assignGIDsToRole", []interface{}{params})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, code)
		if err := SignTransaction(user, tx); err != nil {
			fmt.Println("sign transaction error:", err)
			return
		}
	case "RegistCandidate":
		wallet, err := account.Open(walletPath)
		if err != nil {
			fmt.Println("open assign func to role:"+walletPath+" fail:", err)
			return
		}
		pwd, err := password.GetPassword()
		if err != nil {
			fmt.Println("get password fail")
			return
		}
		user, err := wallet.GetAccountByIndex(index, pwd)
		if err != nil {
			fmt.Println("open wallet " + walletPath + " password error.")
			return
		}
		amount := utils.ParseZpt(value)
		if amount <= 0 {
			fmt.Println("value must >=0.")
			return
		}
		params := &governance.RegisterCandidateParam{
			PeerPubkey: publicKey,
			Address:    user.Address,
			InitPos:    amount,
			Caller:     []byte("GID:ZPT:" + user.Address.ToBase58()),
			KeyNo:      1,
		}
		code, err := BuildNativeInvokeCode(nutils.GovernanceContractAddress, 0, "registerCandidate", []interface{}{params})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, code)
		if err := SignTransaction(user, tx); err != nil {
			fmt.Println("sign transaction error:", err)
			return
		}
	case "UnRegisterCandidate":
		wallet, err := account.Open(walletPath)
		if err != nil {
			fmt.Println("open assign func to role:"+walletPath+" fail:", err)
			return
		}
		pwd, err := password.GetPassword()
		if err != nil {
			fmt.Println("get password fail")
			return
		}
		user, err := wallet.GetAccountByIndex(index, pwd)
		if err != nil {
			fmt.Println("open wallet " + walletPath + " password error.")
			return
		}

		params := &governance.UnRegisterCandidateParam{
			PeerPubkey: publicKey,
			Address:    user.Address,
		}
		code, err := BuildNativeInvokeCode(nutils.GovernanceContractAddress, 0, "unRegisterCandidate", []interface{}{params})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, code)
		if err := SignTransaction(user, tx); err != nil {
			fmt.Println("sign transaction error:", err)
			return
		}
	case "VoteForPeer":
		wallet, err := account.Open(walletPath)
		if err != nil {
			fmt.Println("open assign func to role:"+walletPath+" fail:", err)
			return
		}
		pwd, err := password.GetPassword()
		if err != nil {
			fmt.Println("get password fail")
			return
		}
		user, err := wallet.GetAccountByIndex(index, pwd)
		if err != nil {
			fmt.Println("open wallet " + walletPath + " password error.")
			return
		}
		pos := utils.ParseZpt(value)
		if pos <= 0 {
			fmt.Println("value must >=0.")
			return
		}
		params := &governance.VoteForPeerParam{
			Address:        user.Address,
			PeerPubkeyList: []string{publicKey},
			PosList:        []uint64{pos},
		}
		code, err := BuildNativeInvokeCode(nutils.GovernanceContractAddress, 0, "voteForPeer", []interface{}{params})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, code)
		if err := SignTransaction(user, tx); err != nil {
			fmt.Println("sign transaction error:", err)
			return
		}
	case "UnVoteForPeer":
		wallet, err := account.Open(walletPath)
		if err != nil {
			fmt.Println("open assign func to role:"+walletPath+" fail:", err)
			return
		}
		pwd, err := password.GetPassword()
		if err != nil {
			fmt.Println("get password fail")
			return
		}
		user, err := wallet.GetAccountByIndex(index, pwd)
		if err != nil {
			fmt.Println("open wallet " + walletPath + " password error.")
			return
		}
		pos := utils.ParseZpt(value)
		if err != nil {
			fmt.Println("input value error:", err)
			return
		}
		params := &governance.VoteForPeerParam{
			Address:        user.Address,
			PeerPubkeyList: []string{publicKey},
			PosList:        []uint64{pos},
		}
		code, err := BuildNativeInvokeCode(nutils.GovernanceContractAddress, 0, "unVoteForPeer", []interface{}{params})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, code)
		if err := SignTransaction(user, tx); err != nil {
			fmt.Println("sign transaction error:", err)
			return
		}
	case "approveNode":
		params := &governance.ApproveCandidateParam{
			PeerPubkey: publicKey,
		}
		code, err := BuildNativeInvokeCode(nutils.GovernanceContractAddress, 0, "approveCandidate", []interface{}{params})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, code)
		for i := 0; i < len(config.Configuration.Wallets); i++ {
			if err := MultiSignToTransaction(tx, uint16(math.Ceil(float64(len(config.Configuration.Wallets))*2.0/3.0)), pubKeys, users[i]); err != nil {
				fmt.Println("sign transaction error!")
				return
			}
		}
	case "commitPos":
		code, err := BuildNativeInvokeCode(nutils.GovernanceContractAddress, 0, "commitDpos", []interface{}{false})
		if err != nil {
			fmt.Println("build native invoke code error:", err)
			return
		}
		tx = NewInvokeTransaction(gasPrice, gasLimit, code)
		for i := 0; i < len(config.Configuration.Wallets); i++ {
			if err := MultiSignToTransaction(tx, uint16(math.Ceil(float64(len(config.Configuration.Wallets))*2.0/3.0)), pubKeys, users[i]); err != nil {
				fmt.Println("sign transaction error!")
				return
			}
		}
	case "getPeerpoolInfo":
		preresult, err := PrepareInvokeNativeContract(nutils.GovernanceContractAddress, 0, "getPeerPoolInfo", []interface{}{false})
		if err != nil {
			fmt.Println("get PeerpoolInfo fail:", err)
			return
		}

		peerPoolMap := &governance.PeerPoolMap{
			PeerPoolMap: make(map[string]*governance.PeerPoolItem),
		}

		peerPoolMapbytes, _ := common.HexToBytes(preresult.Result.(string))
		fmt.Println(preresult.Result)
		if err := peerPoolMap.Deserialize(bytes.NewBuffer(peerPoolMapbytes)); err != nil {
			fmt.Println("deserialize, deserialize peerPoolMap error!", err)
		}
		for _, v := range peerPoolMap.PeerPoolMap {
			fmt.Printf("peerInfo Index: %d, InitPos:%d \n", v.Index, v.InitPos)
		}
		return
	case "peersInfo":
		data, err := sendRpcRequest("getblockcount", []interface{}{})
		if err != nil {
			fmt.Println("get block count fail:", err)
			return
		}
		height, err := GetUint32(data)
		if err != nil {
			fmt.Println("get height fail:", err)
			return
		}
		data, err = sendRpcRequest("getblock", []interface{}{height - 1})
		if err != nil {
			fmt.Println("get block data fail:", err)
			return
		}
		blk, err := GetBlock(data)
		if err != nil {
			fmt.Println("get block fail:", err)
			return
		}
		block, err := initVbftBlock(blk)
		if err != nil {
			fmt.Println("init Gbft Block fail:", err)
			return
		}
		var cfg vconfig.ChainConfig
		if block.Info.NewChainConfig != nil {
			cfg = *block.Info.NewChainConfig
		} else {
			var cfgBlock *types.Block
			if block.Info.LastConfigBlockNum != math.MaxUint32 {
				data, err = sendRpcRequest("getblock", []interface{}{block.Info.LastConfigBlockNum})
				if err != nil {
					fmt.Println("get block data again fail")
					return
				}
				cfgBlock, err = GetBlock(data)
				if err != nil {
					fmt.Println("get block again fail")
					return
				}
			}
			blk, err := initVbftBlock(cfgBlock)
			if err != nil {
				fmt.Println("init Gbft Block fail")
				return
			}
			if blk.Info.NewChainConfig == nil {
				return
			}
			cfg = *blk.Info.NewChainConfig
		}
		fmt.Printf("block gbft chainConfig, View:%d, N:%d, C:%d, BlockMsgDelay:%v, HashMsgDelay:%v, PeerHandshakeTimeout:%v, MaxBlockChangeView:%d, PosTable:%v\n",
			cfg.View, cfg.N, cfg.C, cfg.BlockMsgDelay, cfg.HashMsgDelay, cfg.PeerHandshakeTimeout, cfg.MaxBlockChangeView, cfg.PosTable)
		for _, p := range cfg.Peers {
			fmt.Printf("peerInfo Index: %d, ID:%v\n", p.Index, p.ID)
		}
		return
	}

	result, err := utils.SendRawTransaction(tx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("send transaction success, transaction hash:", result)
}

type RegIDWithPublicKeyParam struct {
	GID    []byte
	Pubkey []byte
}

func GetUint32(data []byte) (uint32, error) {
	count := uint32(0)
	err := json.Unmarshal(data, &count)
	if err != nil {
		return 0, fmt.Errorf("json.Unmarshal:%s error:%s", data, err)
	}
	return count, nil
}

func GetBlock(data []byte) (*types.Block, error) {
	hexStr := ""
	err := json.Unmarshal(data, &hexStr)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error:%s", err)
	}
	blockData, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString error:%s", err)
	}
	block := &types.Block{}
	buf := bytes.NewBuffer(blockData)
	err = block.Deserialize(buf)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func NewInvokeTransaction(gasPrice, gasLimit uint64, code []byte) *types.Transaction {
	invokePayload := &payload.InvokeCode{
		Code: code,
	}
	tx := &types.Transaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Nonce:    uint32(time.Now().Unix()),
		Payload:  invokePayload,
		Sigs:     make([]*types.Sig, 0, 0),
	}
	return tx
}

//JsonRpcRequest object in rpc
type JsonRpcRequest struct {
	Version string        `json:"jsonrpc"`
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

//JsonRpcResponse object response for JsonRpcRequest
type JsonRpcResponse struct {
	Error  int64           `json:"error"`
	Desc   string          `json:"desc"`
	Result json.RawMessage `json:"result"`
}

func PrepareInvokeNativeContract(
	contractAddress common.Address,
	version byte,
	method string,
	params []interface{}) (*cstates.PreExecResult, error) {
	tx, err := httpcom.NewNativeInvokeTransaction(0, 0, contractAddress, version, method, params)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	err = tx.Serialize(&buffer)
	if err != nil {
		return nil, fmt.Errorf("Serialize error:%s", err)
	}
	txData := hex.EncodeToString(buffer.Bytes())
	data, err := sendRpcRequest("sendrawtransaction", []interface{}{txData, 1})

	if err != nil {
		return nil, err
	}
	preResult := &cstates.PreExecResult{}
	err = json.Unmarshal(data, &preResult)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal PreExecResult:%s error:%s", data, err)
	}
	return preResult, nil
}
func sendRpcRequest(method string, params []interface{}) ([]byte, error) {
	rpcReq := &JsonRpcRequest{
		Version: "1.0",
		Id:      "cli",
		Method:  method,
		Params:  params,
	}
	data, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("JsonRpcRequest json.Marsha error:%s", err)
	}
	fmt.Println("config.Configuration.RpcUr:", config.Configuration.RpcUrl)
	resp, err := http.Post(config.Configuration.RpcUrl, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("http post request:%s error:%s", data, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rpc response body error:%s", err)
	}
	rpcRsp := &JsonRpcResponse{}
	err = json.Unmarshal(body, rpcRsp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal JsonRpcResponse:%s error:%s", body, err)
	}
	if rpcRsp.Error != 0 {
		return nil, fmt.Errorf("error code:%d desc:%s", rpcRsp.Error, rpcRsp.Desc)
	}
	return rpcRsp.Result, nil
}

func initVbftBlock(block *types.Block) (*vbft.Block, error) {
	if block == nil {
		return nil, fmt.Errorf("nil block in initVbftBlock")
	}

	blkInfo := &vconfig.VbftBlockInfo{}
	if err := json.Unmarshal(block.Header.ConsensusPayload, blkInfo); err != nil {
		return nil, fmt.Errorf("unmarshal blockInfo: %s", err)
	}

	return &vbft.Block{
		Block: block,
		Info:  blkInfo,
	}, nil
}

func BuildNativeInvokeCode(contractAddress common.Address, version byte, method string, params []interface{}) ([]byte, error) {
	builder := neovm.NewParamsBuilder(new(bytes.Buffer))
	err := BuildNeoVMParam(builder, params)
	if err != nil {
		return nil, err
	}
	builder.EmitPushByteArray([]byte(method))
	builder.EmitPushByteArray(contractAddress[:])
	builder.EmitPushInteger(new(big.Int).SetInt64(int64(version)))
	builder.Emit(neovm.SYSCALL)
	builder.EmitPushByteArray([]byte(svrneovm.NATIVE_INVOKE_NAME))
	return builder.ToArray(), nil
}

func BuildNeoVMParam(builder *neovm.ParamsBuilder, smartContractParams []interface{}) error {
	for i := len(smartContractParams) - 1; i >= 0; i-- {
		switch v := smartContractParams[i].(type) {
		case bool:
			builder.EmitPushBool(v)
		case byte:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case int:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case uint:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case int32:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case uint32:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case int64:
			builder.EmitPushInteger(big.NewInt(int64(v)))
		case common.Fixed64:
			builder.EmitPushInteger(big.NewInt(int64(v.GetData())))
		case uint64:
			val := big.NewInt(0)
			builder.EmitPushInteger(val.SetUint64(uint64(v)))
		case string:
			builder.EmitPushByteArray([]byte(v))
		case *big.Int:
			builder.EmitPushInteger(v)
		case []byte:
			builder.EmitPushByteArray(v)
		case common.Address:
			builder.EmitPushByteArray(v[:])
		case common.Uint256:
			builder.EmitPushByteArray(v.ToArray())
		case []interface{}:
			err := BuildNeoVMParam(builder, v)
			if err != nil {
				return err
			}
			builder.EmitPushInteger(big.NewInt(int64(len(v))))
			builder.Emit(neovm.PACK)
		default:
			object := reflect.ValueOf(v)
			kind := object.Kind().String()
			if kind == "ptr" {
				object = object.Elem()
				kind = object.Kind().String()
			}
			switch kind {
			case "slice":
				ps := make([]interface{}, 0)
				for i := 0; i < object.Len(); i++ {
					ps = append(ps, object.Index(i).Interface())
				}
				err := BuildNeoVMParam(builder, []interface{}{ps})
				if err != nil {
					return err
				}
			case "struct":
				builder.EmitPushInteger(big.NewInt(0))
				builder.Emit(neovm.NEWSTRUCT)
				builder.Emit(neovm.TOALTSTACK)
				for i := 0; i < object.NumField(); i++ {
					field := object.Field(i)
					err := BuildNeoVMParam(builder, []interface{}{field.Interface()})
					if err != nil {
						return err
					}
					builder.Emit(neovm.DUPFROMALTSTACK)
					builder.Emit(neovm.SWAP)
					builder.Emit(neovm.APPEND)
				}
				builder.Emit(neovm.FROMALTSTACK)
			default:
				return fmt.Errorf("unsupported param:%s", v)
			}
		}
	}
	return nil
}

func SignTransaction(signer *account.Account, tx *types.Transaction) error {
	tx.Payer = signer.Address
	txHash := tx.Hash()
	sigData, err := Sign(txHash.ToArray(), signer)
	if err != nil {
		return fmt.Errorf("sign error:%s", err)
	}
	sig := &types.Sig{
		PubKeys: []keypair.PublicKey{signer.PublicKey},
		M:       1,
		SigData: [][]byte{sigData},
	}
	tx.Sigs = []*types.Sig{sig}
	return nil
}

func Sign(data []byte, signer *account.Account) ([]byte, error) {
	s, err := sig.Sign(signer.SigScheme, signer.PrivateKey, data, nil)
	if err != nil {
		return nil, err
	}
	sigData, err := sig.Serialize(s)
	if err != nil {
		return nil, fmt.Errorf("sig.Serialize error:%s", err)
	}
	return sigData, nil
}

func MultiSignToTransaction(tx *types.Transaction, m uint16, pubKeys []keypair.PublicKey, signer *account.Account) error {
	pkSize := len(pubKeys)
	if m == 0 || int(m) > pkSize || pkSize > constants.MULTI_SIG_MAX_PUBKEY_SIZE {
		return fmt.Errorf("both m and number of pub key must larger than 0, and small than %d, and m must smaller than pub key number", constants.MULTI_SIG_MAX_PUBKEY_SIZE)
	}
	validPubKey := false
	for _, pk := range pubKeys {
		if keypair.ComparePublicKey(pk, signer.PublicKey) {
			validPubKey = true
			break
		}
	}
	if !validPubKey {
		return fmt.Errorf("Invalid signer")
	}
	var emptyAddress = common.Address{}
	if tx.Payer == emptyAddress {
		payer, err := types.AddressFromMultiPubKeys(pubKeys, int(m))
		if err != nil {
			return fmt.Errorf("AddressFromMultiPubKeys error:%s", err)
		}
		tx.Payer = payer
	}
	txHash := tx.Hash()
	if len(tx.Sigs) == 0 {
		tx.Sigs = make([]*types.Sig, 0)
	}
	sigData, err := SignToData(txHash.ToArray(), signer)
	if err != nil {
		return fmt.Errorf("SignToData error:%s", err)
	}
	hasMutilSig := false
	for i, sigs := range tx.Sigs {
		if pubKeysEqual(sigs.PubKeys, pubKeys) {
			hasMutilSig = true
			if hasAlreadySig(txHash.ToArray(), signer.PublicKey, sigs.SigData) {
				break
			}
			sigs.SigData = append(sigs.SigData, sigData)
			tx.Sigs[i] = sigs
			break
		}
	}
	if !hasMutilSig {
		tx.Sigs = append(tx.Sigs, &types.Sig{
			PubKeys: pubKeys,
			M:       m,
			SigData: [][]byte{sigData},
		})
	}
	return nil
}

func pubKeysEqual(pks1, pks2 []keypair.PublicKey) bool {
	if len(pks1) != len(pks2) {
		return false
	}
	size := len(pks1)
	if size == 0 {
		return true
	}
	pkstr1 := make([]string, 0, size)
	for _, pk := range pks1 {
		pkstr1 = append(pkstr1, hex.EncodeToString(keypair.SerializePublicKey(pk)))
	}
	pkstr2 := make([]string, 0, size)
	for _, pk := range pks2 {
		pkstr2 = append(pkstr2, hex.EncodeToString(keypair.SerializePublicKey(pk)))
	}
	sort.Strings(pkstr1)
	sort.Strings(pkstr2)
	for i := 0; i < size; i++ {
		if pkstr1[i] != pkstr2[i] {
			return false
		}
	}
	return true
}

func hasAlreadySig(data []byte, pk keypair.PublicKey, sigDatas [][]byte) bool {
	for _, sigData := range sigDatas {
		err := signature.Verify(pk, data, sigData)
		if err == nil {
			return true
		}
	}
	return false
}

func SignToData(data []byte, signer *account.Account) ([]byte, error) {
	s, err := sig.Sign(signer.SigScheme, signer.PrivateKey, data, nil)
	if err != nil {
		return nil, err
	}
	sigData, err := sig.Serialize(s)
	if err != nil {
		return nil, fmt.Errorf("sig.Serialize error:%s", err)
	}
	return sigData, nil
}
