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


# ZEEPIN同步节点的部署

## 获取zeepin
### 从release获取
- 从[下载页面](https://github.com/zeepin/zeepinChain/releases)获取

## 服务器部署
### 选择网络
zeepin的运行支持以下方式

* TestNet同步节点部署
* MainNet同步节点部署
* MainNet竞选节点部署

参数：--networkid：1为默认主网；2为测网；3为单机运行；

#### TestNet同步节点部署

运行zeepin

   ```
	./zeepin --networkid 2
   ```

#### MainNet同步节点部署

运行zeepin

   ```
	./zeepin
   ```

#### MainNet竞选节点部署

先通过zeepin cli [创建zeepin钱包](#创建zeepin钱包)后运行zeepin

   ```
	./zeepin --enableconsensus --rest --restport=20334 --ws --wsport=20335 --rpcport=20336 --nodeport=20338
   ```

查看更多命令参数：

   ```
	--loglevel=0 日志参数
	--password=xxx 钱包密码
	./zeepin --help
   ```


