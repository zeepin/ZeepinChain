<h1 align="center">EXCHANGE_GUIDELINES_CN</h1>
<h4 align="center">Version 0.1.0 </h4>

[English](EXCHANGE_GUIDELINES.md) | [中文](EXCHANGE_GUIDELINES_CN.md) | [한글](EXCHANGE_GUIDELINES_KO.md)


ZEEPIN区块链资产分为：

* 原生资产：
  * ZPT：主币、治理、资产等
  * Gala：交易转账、合约部署、系统燃料等
  
* 合约资产：
  * 通过合约产生的token
  * 如MOX等

交易所和ZEEPIN区块链对接时，主要是处理这两种类型资产的：充值、提现、交易查询、账户创建等操作。



## 1、ZEEPIN同步节点的部署

### 获取zeepin
#### 从release获取
- 从[下载页面](https://github.com/zeepin/zeepinChain/releases)获取

## 服务器部署
#### MainNet同步节点部署

目录结构如下

   ```
	$ tree -L 1
	.
	├── zeepin
	└── wallet.dat
   ```



运行zeepin节点

   ```
	./zeepin
   ```
   

默认关闭websocket和rest端口，可以配置以下参数启动端口：


   ```
RESTFUL OPTIONS:
  --rest            Enable restful api server
  --restport value  Restful server listening port (default: 20334)

WEB SOCKET OPTIONS:
  --ws            Enable websocket server
  --wsport value  Ws server listening port (default: 20335)
   
   ```
  
zeepin -h 查看更多命令，如参数：

   ```
	--loglevel=0 日志参数
   ```


## 2、ZEEPIN CLI客户端使用

### 安全策略

强制要求交易所使用白名单和防火墙，隔绝外部请求，否则会有重大安全隐患。

ZEEPIN CLI 自身不提供远程开关钱包功能，打开钱包时也没有验证过程。因此，安全策略由交易所根据自身情况制定。由于钱包要一直保持打开状态以便处理用户的提现，因此，从安全角度考虑，钱包必须运行在独立的服务器上，并参考下表配置好端口防火墙。

|   port type   | Mainnet default port |
| ------------- | -------------------- |
| Rest Port     | 20334                |
| Websorcket    | 20335                |
| Json RPC port | 20336                |
| Node port     | 20338                |


### 创建钱包

交易所需要创建一个在线钱包管理用户充值地址。钱包是用来存储账户（公钥和私钥）、合约地址等信息，是用户持有资产的最重要的凭证，一定要保管好钱包文件和钱包密码，不要丢失或泄露。 交易所不需要为每个地址创建一个钱包文件，通常一个钱包文件可以存储用户所有充值地址。也可以使用一个冷钱包（离线钱包）作为更安全的存储方式。


创建zeepin钱包


   ```
	./zeepin account add -d
   ```
   
接着输入密码，执行完输出：

```
Use default setting '-t ecdsa -b 256 -s SHA256withECDSA'
	signature algorithm: ecdsa
	curve: P-256
	signature scheme: SHA256withECDSA
Password:
Re-enter Password:

Index: 1
Label:
Address: ZT047K36grEi5H6BF7gLb2Z0JwBFMQRRCU
Public key: 02c7fed64a315c664034bae1257f45c9fdf8c24033f0904ce7b47b0090232323
Signature scheme: SHA256withECDSA

```

zeepin钱包统一为Z开头，

请务必保存好钱包密码和私钥！
钱包地址大小写敏感，请务必注意！

zeepin钱包私钥生成算法和NEO一致，同一个私钥对应的ZPT和NEO的公钥地址不相同。


###  生成充值地址

####  充值地址有两种生成方式：

- 动态生成：用户创建账户时通过 Java SDK 实现动态创建 ZPT/Gala 地址并返回（优点：自动维护 / 缺点：备份不太方便）

- 批量创建：交易所批量生成，用户创建账户后分配给用户 ZPT 地址（优点：方便备份 / 缺点：定期人工生成）

  
  批量创建方法： ZEEPIN CLI 执行：
  
  ```
  ./zeepin account add -d -n [数量]  -w [钱包文件名]
  
  ```
  
  -d 默认值为 1，即调用默认设置
  -n 批量创建的地址数量
  -w 指定钱包文件，默认为wallet.dat
  
  例如要一次创建10个zeepin地址:

```
	$ ./zeepin account add -d -n 10 -w walletEx.dat
	
	Use default setting '-t ecdsa -b 256 -s SHA256withECDSA'
		signature algorithm: ecdsa
		curve: P-256
		signature scheme: SHA256withECDSA
		
	Password:
	Re-enter Password:

	Index: 1
	Label:
	Address: ZDahjrxYu2vFUPQMDZx6bRtF7RDJs4Xkxb
	Public key: 03e38e4944da8a7bba7a12e88a679b0b0063e2362b4b82cf6bdfdc49b5c595418f
	Signature scheme: SHA256withECDSA

	Index: 2
	Label:
	Address: ZW7gQeK9o4Axs6Dg17E7yiBsUnnsRaGWbF
	Public key: 03c86673df828cc90f03497dd72eb785daf5f0b1561e4d8250eea05114a6652343
	Signature scheme: SHA256withECDSA

	Index: 3
	Label:
	Address: ZEzF5pBkfamfHBSDYvLSJbatM4iSLsZF3c
	Public key: 02f10b9250391a309e26481faa3b7c5f041dd0a7f62277cdd81527d8defce63e20
	Signature scheme: SHA256withECDSA
	
	....... .......

	Index: 10
	Label:
	Address: ZQm7vGM5ezFdJhUigbLgv2ekzvxpCp3Eqv
	Public key: 036220266424bcf381cfae8233e6422a9761ae05cd31036ede89301a940af3b353
	Signature scheme: SHA256withECDSA

	Create account successfully.

```






