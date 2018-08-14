
<h1 align="center">zeepin </h1>
<h4 align="center">Version 0.1 </h4>

[English](README.md) | [中文](README_CN.md) | [한글](README_KO.md)


[Zeepin Chain Whitepaper EN](https://www.zeepin.io/pdfs/Zeepin%20Chain%20Tech%20WP%20V1.0%20EN.pdf) | [Zeepin Chain Whitepaper CH](https://www.zeepin.io/pdfs/Zeepin%20Chain_WP%20CN%20V1.0.pdf)


欢迎查看zeepin的源码库!

Zeepin Chain是一条去中心化的文创及娱乐资产公链，通过区块链构建标准化基础设施，为创意人群提供高效工作的解决方案，帮助创意组织提高创新效率，促进文创产业开放透明、公平高效的价值流通。同时Zeepin Chain还将打造区块链数字娱乐资产发行平台，为全球文娱资产代币化提供区块链技术支持及落地场景的建设。Zeepin Chain公链作为一条行业基础链，拥有整合第三方娱乐资产和系统的能力，建立一个自由的交易市场和兑换平台。

Zeepin Chain构建了完整的区块链技术框架，采用GBFT－POS共识机制（星际共识），提供具备图灵完备性的虚拟机作为智能合约的执行环境，为应用架构提供自定义脚步控制支持。支持Java、C#、Python、Javascript等编程语言开发的脚本，虚拟机都可以通过api与链进行集成交互。

zeepin致力于创建一个可自由配置、高性能、可扩展的区块链底层基础设施，让部署及调用去中心化应用变得更加的简单。目前代码还处于内部测试阶段，但处于快速的迭代开发中，欢迎及希望更多的开发者加入到zeepin中来！


## 目录

* [获取zeepin](#获取zeepin)
    * [从release获取](#从release获取)
* [服务器部署](#服务器部署)
    * [选择网络](#选择网络)
        * [TestNet同步节点部署](#testnet同步节点部署)
        * [MainNet同步节点部署](#mainnet同步节点部署)
        * [MainNet竞选节点部署](#mainnet竞选节点部署)
    * [创建zeepin钱包](#创建zeepin钱包)
    * [ZPT转账调用示例](#zpt转账调用示例)
    * [查询转账结果TxHash示](#查询转账结果txhash)
    * [查询账户余额示例](#查询账户余额示例)
    * [查询解绑的Gala示例](#查询解绑的gala示例)
    * [提取解绑的Gala示例](#提取解绑的gala示例)
    * [查询区块信息示例](#查询区块信息示例)
* [官方社区](#官方社区)
    * [官方网站](#官方网站)
    * [Telegram](#telegram)
* [许可证](#许可证)


## 获取zeepin
### 从release获取
- 从[下载页面](https://github.com/zeepin/zeepinChain/releases)获取

## 服务器部署
### 选择网络
zeepin的运行支持以下方式

* TestNet同步节点部署
* MainNet同步节点部署
* MainNet竞选节点部署


#### TestNet同步节点部署

运行zeepin

   ```
	./zeepin --networkid 2
   ```

#### MainNet同步节点部署

运行zeepin

   ```
	./zeepin --networkid 1
   ```

#### MainNet竞选节点部署

运行zeepin

   ```
	./zeepin --enableconsensus --rest --restport=20334 --ws --wsport=20335 --rpcport=20336 --nodeport=20338 --loglevel=0
   ```

查看更多命令参数：

   ```
	./zeepin --help
   ```


#### 创建zeepin钱包

创建zeepin钱包

   ```
	./zeepin account add -d
   ```
   接着输入密码即可创建钱包，请务必保存好钱包密码和私钥
   
   
### ZPT转账调用示例
   - from: 转出地址； - to: 转入地址； - amount: 转出资产数量；

```shell
  ./zeepin asset transfer --from ZJohWxMxiMWHczSCV5ZUybZEf5jh9VQE5G  --to ZTA8f5U8Zd1gELkje7fJmYmdmMMm1WxPyv --amount 5
```

执行完输出：

```
Transfer ZPT
  From:ZJohWxMxiMWHczSCV5ZUybZEf5jh9VQE5G
  To:ZTA8f5U8Zd1gELkje7fJmYmdmMMm1WxPyv
  Amount:5
  TxHash:78cd8097ed58cf06b8f7d7591f05657afb9f32686ef88956c246b3a4146f6ec2

Tip:
  Using './zeepin info status 78cd8097ed58cf06b8f7d7591f05657afb9f32686ef88956c246b3a4146f6ec2' to query transaction status
```
可以通过这个TxHash查询转账交易的结果，等待至少一个区块时间即可。


### 查询转账结果TxHash

```shell
./zeepin info status 78cd8097ed58cf06b8f7d7591f05657afb9f32686ef88956c246b3a4146f6ec2
```
查询结果：
```shell
Transaction states:
{
   "TxHash": "78cd8097ed58cf06b8f7d7591f05657afb9f32686ef88956c246b3a4146f6ec2",
   "State": 1,
   "GasConsumed": 10000000,
   "Notify": [
      {
         "ContractAddress": "0100000000000000000000000000000000000000",
         "States": [
            "transfer",
            "ZJohWxMxiMWHczSCV5ZUybZEf5jh9VQE5G",
            "ZTA8f5U8Zd1gELkje7fJmYmdmMMm1WxPyv",
            50000
         ]
      },
      {
         "ContractAddress": "0200000000000000000000000000000000000000",
         "States": [
            "transfer",
            "ZJohWxMxiMWHczSCV5ZUybZEf5jh9VQE5G",
            "ZC3Fmgr3oS56Rg9vxZeVo2mwMMcUiYGcPp",
            10000000
         ]
      }
   ]
}
```
查询TxHash结果数据中的数值应除去10000,因为ZPT和Gala的精度为4位数；


### 查询账户余额示例

```shell
./zeepin asset balance ZJohWxMxiMWHczSCV5ZUybZEf5jh9VQE5G
```
查询结果：
```shell
  BalanceOf:ZJohWxMxiMWHczSCV5ZUybZEf5jh9VQE5G
  ZPT:100098945
  GALA:100026953.517
```


### 查询解绑的Gala示例

```shell
./zeepin asset unboundgala ZJohWxMxiMWHczSCV5ZUybZEf5jh9VQE5G
```

查询结果：
```shell
  Unbound GALA:
  Account:ZJohWxMxiMWHczSCV5ZUybZEf5jh9VQE5G
  GALA:10129.2123
```

### 提取解绑的Gala示例

```shell
./zeepin asset withdrawgala ZJohWxMxiMWHczSCV5ZUybZEf5jh9VQE5G
```

### 查询区块信息示例

```shell
./zeepin info block 1
```


## 官方社区

### 官方网站

- https://zeepin.io/
- https://twitter.com/ZeepinChain
- https://medium.com/@zeepin
- https://www.reddit.com/r/ZEEPIN/
- https://www.facebook.com/ZeepinChain/

### Telegram

- https://t.me/ZeepinNews
- https://t.me/zeepin


## 许可证

zeepin遵守GNU Lesser General Public License, 版本3.0。 详细信息请查看项目根目录下的LICENSE文件。
