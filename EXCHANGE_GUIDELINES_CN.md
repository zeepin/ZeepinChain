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



## ZEEPIN同步节点的部署

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


1、创建zeepin钱包


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
请务必保存好钱包密码和私钥,zeepin钱包统一为Z开头。


2、运行zeepin节点

   ```
	./zeepin
   ```

默认关闭websocket和rest端口，可以配置以下参数启动端口：

RESTFUL OPTIONS:
  --rest            Enable restful api server
  --restport value  Restful server listening port (default: 20334)

WEB SOCKET OPTIONS:
  --ws            Enable websocket server
  --wsport value  Ws server listening port (default: 20335)
  
  
zeepin -h 查看更多命令，如参数：

   ```
	--loglevel=0 日志参数
   ```


