# Native Contract API : Param
* [Introduction](#introduction)
* [Contract Method](#contract-method)

## Introduction
本文档是全局参数管理合约的使用说明文档
## Contract Method

### ParamInit
初始化合约，在创世区块中被调用。

method: init

args: nil

return: bool

#### example
```
    init := sstates.Contract{
		Address: ParamContractAddress,
		Method:  "init",
	}
```
### TransferAdmin
合约管理员更换，需要由当前管理员调用。

method: transferAdmin

args: smartcontract/service/native/global_params.Admin

return: bool

#### example
```
    var destinationAdmin global_params.Admin
	address, _ := common.AddressFromBase58("TA4knXiWFZ8K4W3e5fAnoNntdc5G3qMT7C")
	copy(destinationAdmin[:], address[:])
	adminBuffer := new(bytes.Buffer)
	if err := destinationAdmin.Serialize(adminBuffer); err != nil {
		fmt.Println("Serialize admins struct error.")
		os.Exit(1)
	}
	contract := &sstates.Contract{
		Address: genesis.ParamContractAddress,
		Method:  "transferAdmin",
		Args:    adminBuffer.Bytes(),
	}
```

### AcceptAdmin
接受合约Admin权限
method: acceptAdmin

args: smartcontract/service/native/global_params.Admin

return: bool

#### example
```
    var destinationAdmin global_params.Admin
	address, _ := common.AddressFromBase58("TA4knXiWFZ8K4W3e5fAnoNntdc5G3qMT7C")
	copy(destinationAdmin[:], address[:])
	adminBuffer := new(bytes.Buffer)
	if err := destinationAdmin.Serialize(adminBuffer); err != nil {
		fmt.Println("Serialize admin struct error.")
		os.Exit(1)
	}

	contract := &sstates.Contract{
		Address: genesis.ParamContractAddress,
		Method:  "acceptAdmin",
		Args:    adminBuffer.Bytes(),
	}
```
### SetOperator
管理员设置合约操作员
method: setOperator

args: smartcontract/service/native/global_params.Admin

return: bool
#### example
```
    var destinationOperator global_params.Admin
	address, _ := common.AddressFromBase58("TA4knXiWFZ8K4W3e5fAnoNntdc5G3qMT7C")
	copy(destinationOperator[:], address[:])
	adminBuffer := new(bytes.Buffer)
	if err := destinationOperator.Serialize(adminBuffer); err != nil {
		fmt.Println("Serialize admin struct error.")
		os.Exit(1)
	}

	contract := &sstates.Contract{
		Address: genesis.ParamContractAddress,
		Method:  "acceptAdmin",
		Args:    adminBuffer.Bytes(),
	}
```

### SetGlobalParam
操作员设置全局参数，该参数设置并不会理解生效。
method: setGlobalParam

args: smartcontract/service/native/global_params.Params

return: bool

#### example
```
    params := new(global_params.Params)
	*params = make(map[string]string)
	for i := 0; i < 3; i++ {
		k := "key-test" + strconv.Itoa(i) + "-" + key
		v := "value-test" + strconv.Itoa(i) + "-" + value
		(*params) = append(*params, &global_params.Param{k,v})
	}
	paramsBuffer := new(bytes.Buffer)
	if err := params.Serialize(paramsBuffer); err != nil {
		fmt.Println("Serialize params struct error.")
		os.Exit(1)
	}

	contract := &sstates.Contract{
		Address: genesis.ParamContractAddress,
		Method:  "setGlobalParam",
		Args:    paramsBuffer.Bytes(),
	}
```

### GetGlobalParam
获取全局参数，该函数将会返回 smartcontract/service/native/global_params.Params

method: getGlobalParam

args: smartcontract/service/native/global_params.ParamNameList

return: array

#### example
```
    nameList := new(global_params.ParamNameList)
	for i := 0; i < 3; i++ {
		k := "key-test" + strconv.Itoa(i) + "-" + key
		(*nameList) = append(*nameList, k)
	}
	nameListBuffer := new(bytes.Buffer)
	if err := nameList.Serialize(nameListBuffer); err != nil {
		fmt.Println("Serialize ParamNameList struct error.")
		os.Exit(1)
	}
	contract := &sstates.Contract{
		Address: genesis.ParamContractAddress,
		Method:  "getGlobalParam",
		Args:    nameListBuffer.Bytes(),
	}
```

### CreateSnapshot
管理员使参数设置生效。

method: createSnapshot

args: nil

return: bool

#### example
```
    contract := &sstates.Contract{
		Address: genesis.ParamContractAddress,
		Method:  "createSnapshot",
	}
```
