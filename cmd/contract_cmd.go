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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	cmdcom "github.com/imZhuFei/zeepin/cmd/common"
	"github.com/imZhuFei/zeepin/cmd/utils"
	"github.com/imZhuFei/zeepin/common"
	"github.com/imZhuFei/zeepin/common/config"
	"github.com/imZhuFei/zeepin/core/types"
	httpcom "github.com/imZhuFei/zeepin/http/base/common"
	"github.com/imZhuFei/zeepin/smartcontract/service/wasmvm"
	cstates "github.com/imZhuFei/zeepin/smartcontract/states"
	"github.com/urfave/cli"
)

var (
	ContractCommand = cli.Command{
		Name:        "contract",
		Action:      cli.ShowSubcommandHelp,
		Usage:       "Deploy or invoke smart contract",
		ArgsUsage:   " ",
		Description: `Smart contract operations support the deployment of WASMVM smart contract, and the pre-execution and execution of WASMVM smart contract.`,
		Subcommands: []cli.Command{
			{
				Action:    deployContract,
				Name:      "deploy",
				Usage:     "Deploy a smart contract to ontolgoy",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					utils.RPCPortFlag,
					utils.TransactionGasPriceFlag,
					utils.TransactionGasLimitFlag,
					utils.ContractStorageFlag,
					utils.ContractCodeFileFlag,
					utils.ContractNameFlag,
					utils.ContractVersionFlag,
					utils.ContractAuthorFlag,
					utils.ContractEmailFlag,
					utils.ContractDescFlag,
					utils.ContractPrepareDeployFlag,
					utils.WalletFileFlag,
					utils.AccountAddressFlag,
					utils.ContractAttrFlag,
				},
			},
			{
				Action:    invokeContract,
				Name:      "invoke",
				Usage:     "Invoke smart contract",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					utils.RPCPortFlag,
					utils.TransactionGasPriceFlag,
					utils.TransactionGasLimitFlag,
					utils.ContractAddrFlag,
					utils.ContractParamsFlag,
					utils.ContractAttrFlag,
					utils.ContractMethodFlag,
					utils.ContractParamTypeFlag,
					utils.ContractVersionFlag,
					utils.ContractPrepareInvokeFlag,
					utils.ContractReturnTypeFlag,
					utils.WalletFileFlag,
					utils.AccountAddressFlag,
				},
			},
			{
				Action:    invokeCodeContract,
				Name:      "invokeCode",
				Usage:     "Invoke smart contract by code",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					utils.RPCPortFlag,
					utils.ContractCodeFileFlag,
					utils.TransactionGasPriceFlag,
					utils.TransactionGasLimitFlag,
					utils.WalletFileFlag,
					utils.ContractPrepareInvokeFlag,
					utils.AccountAddressFlag,
				},
			},
		},
	}
)

func deployContract(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if !ctx.IsSet(utils.GetFlagName(utils.ContractCodeFileFlag)) ||
		!ctx.IsSet(utils.GetFlagName(utils.ContractNameFlag)) {
		fmt.Errorf("Missing code or name argument\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	store := ctx.Bool(utils.GetFlagName(utils.ContractStorageFlag))
	codeFile := ctx.String(utils.GetFlagName(utils.ContractCodeFileFlag))
	if "" == codeFile {
		return fmt.Errorf("Please specific code file")
	}
	codeStr, err := ioutil.ReadFile(codeFile)
	if err != nil {
		return fmt.Errorf("Read code:%s error:%s", codeFile, err)
	}

	name := ctx.String(utils.GetFlagName(utils.ContractNameFlag))
	version := ctx.String(utils.GetFlagName(utils.ContractVersionFlag))
	author := ctx.String(utils.GetFlagName(utils.ContractAuthorFlag))
	email := ctx.String(utils.GetFlagName(utils.ContractEmailFlag))
	desc := ctx.String(utils.GetFlagName(utils.ContractDescFlag))
	code := strings.TrimSpace(string(codeStr))
	gasPrice := ctx.Uint64(utils.GetFlagName(utils.TransactionGasPriceFlag))
	gasLimit := ctx.Uint64(utils.GetFlagName(utils.TransactionGasLimitFlag))
	networkId, err := utils.GetNetworkId()
	cattr := ctx.Uint64(utils.GetFlagName(utils.ContractAttrFlag))
	if err != nil {
		return err
	}
	if networkId == config.NETWORK_ID_SOLO_NET {
		gasPrice = 0
	}

	cversion := fmt.Sprintf("%s", version)

	if ctx.IsSet(utils.GetFlagName(utils.ContractPrepareDeployFlag)) {
		preResult, err := utils.PrepareDeployContract(store, code, name, cversion, author, email, desc, cattr)
		if err != nil {
			return fmt.Errorf("PrepareDeployContract error:%s", err)
		}
		if preResult.State == 0 {
			return fmt.Errorf("Contract pre-deploy failed\n")
		}
		fmt.Printf("Contract pre-deploy successfully\n")
		fmt.Printf("Gas consumed:%d\n", preResult.Gas)
		return nil
	}

	signer, err := cmdcom.GetAccount(ctx)
	if err != nil {
		return fmt.Errorf("Get signer account error:%s", err)
	}

	txHash, err := utils.DeployContract(gasPrice, gasLimit, signer, store, code, name, cversion, author, email, desc, cattr)
	if err != nil {
		return fmt.Errorf("DeployContract error:%s", err)
	}
	var c []byte
	if cattr == 0 {
		c, _ = common.HexToBytes(code)
	} else {
		c = []byte(code)
	}
	address := types.AddressFromVmCode(c)
	fmt.Printf("Deploy contract:\n")
	fmt.Printf("  Contract Address:%s\n", address.ToHexString())
	fmt.Printf("  TxHash:%s\n", txHash)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './zeepin info status %s' to query transaction status\n", txHash)
	return nil
}

func invokeCodeContract(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if !ctx.IsSet(utils.GetFlagName(utils.ContractCodeFileFlag)) {
		fmt.Printf("Missing code or name argument\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}

	codeFile := ctx.String(utils.GetFlagName(utils.ContractCodeFileFlag))
	if "" == codeFile {
		return fmt.Errorf("Please specific code file")
	}
	codeStr, err := ioutil.ReadFile(codeFile)
	if err != nil {
		return fmt.Errorf("Read code:%s error:%s", codeFile, err)
	}
	code := strings.TrimSpace(string(codeStr))
	c, err := common.HexToBytes(code)
	if err != nil {
		return fmt.Errorf("hex to bytes error:%s", err)
	}
	cattr := ctx.Uint64(utils.GetFlagName(utils.ContractAttrFlag))

	if ctx.IsSet(utils.GetFlagName(utils.ContractPrepareInvokeFlag)) {
		var preResult *cstates.PreExecResult
		var err error
		if cattr == 0 {
			preResult, err = utils.PrepareInvokeCodeEmbeddedContract(c)
		} else {
			//cmethod := ctx.String(utils.GetFlagName(utils.ContractMethodFlag))
			//paramType := ctx.Uint64(utils.GetFlagName(utils.ContractParamTypeFlag))
			//preResult, err = utils.PrepareInvokeWASMVMContract(contractAddr, cmethod, wasmvm.ParamType(paramType), 1, params)
		}
		if err != nil {
			return fmt.Errorf("PrepareInvokeCodeEmbeddedContract error:%s", err)
		}
		if preResult.State == 0 {
			return fmt.Errorf("Contract pre-invoke failed\n")
		}
		fmt.Printf("Contract pre-invoke successfully\n")
		fmt.Printf("Gas consumed:%d\n", preResult.Gas)

		rawReturnTypes := ctx.String(utils.GetFlagName(utils.ContractReturnTypeFlag))
		if rawReturnTypes == "" {
			fmt.Printf("Return:%s (raw value)\n", preResult.Result)
			return nil
		}
		values, err := utils.ParseReturnValue(preResult.Result, rawReturnTypes)
		if err != nil {
			return fmt.Errorf("parseReturnValue values:%+v types:%s error:%s", values, rawReturnTypes, err)
		}
		switch len(values) {
		case 0:
			fmt.Printf("Return: nil\n")
		case 1:
			fmt.Printf("Return:%+v\n", values[0])
		default:
			fmt.Printf("Return:%+v\n", values)
		}
		return nil
	}
	gasPrice := ctx.Uint64(utils.GetFlagName(utils.TransactionGasPriceFlag))
	gasLimit := ctx.Uint64(utils.GetFlagName(utils.TransactionGasLimitFlag))
	networkId, err := utils.GetNetworkId()
	if err != nil {
		return err
	}
	if networkId == config.NETWORK_ID_SOLO_NET {
		gasPrice = 0
	}

	invokeTx, err := httpcom.NewSmartContractTransaction(gasPrice, gasLimit, c, byte(cattr))
	if err != nil {
		return err
	}

	signer, err := cmdcom.GetAccount(ctx)
	if err != nil {
		return fmt.Errorf("Get signer account error:%s", err)
	}

	err = utils.SignTransaction(signer, invokeTx)
	if err != nil {
		return fmt.Errorf("SignTransaction error:%s", err)
	}
	tx, err := invokeTx.IntoImmutable()
	if err != nil {
		return err
	}
	txHash, err := utils.SendRawTransaction(tx)
	if err != nil {
		return fmt.Errorf("SendTransaction error:%s", err)
	}

	fmt.Printf("TxHash:%s\n", txHash)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './zeepin info status %s' to query transaction status\n", txHash)
	return nil
}

func invokeContract(ctx *cli.Context) error {
	SetRpcPort(ctx)
	if !ctx.IsSet(utils.GetFlagName(utils.ContractAddrFlag)) {
		fmt.Printf("Missing contract address argument.\n")
		cli.ShowSubcommandHelp(ctx)
		return nil
	}
	contractAddrStr := ctx.String(utils.GetFlagName(utils.ContractAddrFlag))
	contractAddr, err := common.AddressFromHexString(contractAddrStr)
	if err != nil {
		return fmt.Errorf("Invalid contract address error:%s", err)
	}

	paramsStr := ctx.String(utils.GetFlagName(utils.ContractParamsFlag))
	params, err := utils.ParseParams(paramsStr)
	if err != nil {
		return fmt.Errorf("parseParams error:%s", err)
	}

	paramData, _ := json.Marshal(params)
	fmt.Printf("Invoke:%x Params:%s\n", contractAddr[:], paramData)

	attr := ctx.Uint64(utils.GetFlagName(utils.ContractAttrFlag))

	if ctx.IsSet(utils.GetFlagName(utils.ContractPrepareInvokeFlag)) {
		var preResult *cstates.PreExecResult
		var err error
		if attr == 0 {
			preResult, err = utils.PrepareInvokeEmbeddedContract(contractAddr, params)
		} else {
			cmethod := ctx.String(utils.GetFlagName(utils.ContractMethodFlag))
			paramType := ctx.Uint64(utils.GetFlagName(utils.ContractParamTypeFlag))
			preResult, err = utils.PrepareInvokeWASMVMContract(contractAddr, cmethod, wasmvm.ParamType(paramType), 1, params, byte(attr))
		}
		if err != nil {
			if attr == 0 {
				return fmt.Errorf("PrepareInvokeEmbeddedSmartContact error:%s", err)
			} else {
				return fmt.Errorf("PrepareInvokeWASMVMSmartContact error:%s", err)
			}
		}
		if preResult.State == 0 {
			return fmt.Errorf("Contract invoke failed\n")
		}
		fmt.Printf("Contract invoke successfully\n")
		fmt.Printf("  Gaslimit:%d\n", preResult.Gas)

		rawReturnTypes := ctx.String(utils.GetFlagName(utils.ContractReturnTypeFlag))
		if rawReturnTypes == "" {
			fmt.Printf("  Return:%s (raw value)\n", preResult.Result)
			return nil
		}
		values, err := utils.ParseReturnValue(preResult.Result, rawReturnTypes)
		if err != nil {
			return fmt.Errorf("parseReturnValue values:%+v types:%s error:%s", values, rawReturnTypes, err)
		}
		switch len(values) {
		case 0:
			fmt.Printf("  Return: nil\n")
		case 1:
			fmt.Printf("  Return:%+v\n", values[0])
		default:
			fmt.Printf("  Return:%+v\n", values)
		}
		return nil
	}
	signer, err := cmdcom.GetAccount(ctx)
	if err != nil {
		return fmt.Errorf("Get signer account error:%s", err)
	}
	gasPrice := ctx.Uint64(utils.GetFlagName(utils.TransactionGasPriceFlag))
	gasLimit := ctx.Uint64(utils.GetFlagName(utils.TransactionGasLimitFlag))
	networkId, err := utils.GetNetworkId()
	if err != nil {
		return err
	}
	if networkId == config.NETWORK_ID_SOLO_NET {
		gasPrice = 0
	}
	var txHash string
	if attr == 0 {
		txHash, err = utils.InvokeEmbeddedContract(gasPrice, gasLimit, signer, contractAddr, params)
	} else {
		cmethod := ctx.String(utils.GetFlagName(utils.ContractMethodFlag))
		paramType := ctx.Uint64(utils.GetFlagName(utils.ContractParamTypeFlag))
		txHash, err = utils.InvokeWasmVMContract(gasPrice, gasLimit, signer, 1, contractAddr, cmethod, wasmvm.ParamType(paramType), params)
	}
	if err != nil {
		return fmt.Errorf("Invoke Embedded contract error:%s", err)
	}

	fmt.Printf("  TxHash:%s\n", txHash)
	fmt.Printf("\nTip:\n")
	fmt.Printf("  Using './zeepin info status %s' to query transaction status\n", txHash)
	return nil
}
