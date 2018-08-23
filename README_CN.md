
<h1 align="center">ZEEPIN</h1>
<h4 align="center">Version 0.1.0 </h4>

[English](README.md) | [中文](README_CN.md) | [한글](README_KO.md)


[Zeepin Chain Whitepaper EN](https://www.zeepin.io/pdfs/Zeepin%20Chain%20Tech%20WP%20V1.0%20EN.pdf) | [Zeepin Chain Whitepaper CH](https://www.zeepin.io/pdfs/Zeepin%20Chain_WP%20CN%20V1.0.pdf)


欢迎查看zeepin的源码库!

Zeepin Chain是一条去中心化的文创及娱乐资产公链，通过区块链构建标准化基础设施，为创意人群提供高效工作的解决方案，帮助创意组织提高创新效率，促进文创产业开放透明、公平高效的价值流通。同时Zeepin Chain还将打造区块链数字娱乐资产发行平台，为全球文娱资产代币化提供区块链技术支持及落地场景的建设。Zeepin Chain公链作为一条行业基础链，拥有整合第三方娱乐资产和系统的能力，建立一个自由的交易市场和兑换平台。

Zeepin Chain构建了完整的区块链技术框架，采用GBFT－POS共识机制（星际共识），提供具备图灵完备性的虚拟机作为智能合约的执行环境，为应用架构提供自定义脚步控制支持。支持Java、C#、Python、Javascript等编程语言开发的脚本，虚拟机都可以通过api与链进行集成交互。

zeepin致力于创建一个可自由配置、易扩展、高性能的区块链底层基础设施，让部署区块链环境及开发dApp变得更加的简单。zeepin v0.1.0版本根据文创行业需求基于本体1.0核心框架进行定制开发，目前代码处于快速迭代开发、测试中，欢迎更多的开发者加入到zeepin技术社区中来！


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
    * [查询转账结果TxHash](#查询转账结果txhash)
    * [查询账户余额示例](#查询账户余额示例)
    * [查询解绑的Gala示例](#查询解绑的gala示例)
    * [提取解绑的Gala示例](#提取解绑的gala示例)
    * [查询当前区块高度示例](#查询当前区块高度示例)
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

先[创建zeepin钱包](#创建zeepin钱包)后运行zeepin

   ```
	./zeepin --enableconsensus --rest --restport=20334 --ws --wsport=20335 --rpcport=20336 --nodeport=20338
   ```

查看更多命令参数：

   ```
	--loglevel=0 日志参数
	--password=xxx 钱包密码
	./zeepin --help
   ```


#### 创建zeepin钱包


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
查询TxHash结果数据中的数值应除以10000,因为ZPT和Gala的精度为4位数；


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


### 查询当前区块高度示例

```shell
./zeepin info curblockheight
```

显示当前区块高度：

```
CurrentBlockHeight:1693
```

### 查询区块信息示例

```shell
./zeepin info block 1693
```

执行显示：
```
{
   "Hash": "ea0720ca0ad21b408136bf233c4a9e11be26e8d26cee676c23d8463abbc0200a",
   "Size": 1011,
   "Header": {
      "Version": 0,
      "PrevBlockHash": "56f92d7953d9adec9afa8a029310992e15f23ad9308ef27cbe6766a3cb9245d1",
      "TransactionsRoot": "0000000000000000000000000000000000000000000000000000000000000000",
      "BlockRoot": "8ae95b3b5773532804e115acc38e890f64ae26e2242bdc66ba3b486e4ad2a2d9",
      "Timestamp": 1534287845,
      "Height": 1693,
      "ConsensusData": 13335926600357978572,
      "ConsensusPayload": "7b226c6561646572223a31312c227672665f76616c7565223a2242476162765173586334755a7641336e454549656a7853723177566849764833723330613864685471794c573872634b2b3872534246374a67476a4b46553741684238446c505a765a5542757271595364525a544857303d222c227672665f70726f6f66223a225a486457346e756b4443704e616a567454386b47344367442b3539425548324a37384c3571797951507561347335616767467856636b62705475486158304236596d4e676e664f42364d3047434a68743243697346773d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a313638372c226e65775f636861696e5f636f6e666967223a6e756c6c7d",
      "NextBookkeeper": "ZC3Fmgr3oS56Rg9vxZeVo2mwMMcTvb44rA",
      "Bookkeepers": [
         "02127d484040ddb9bba2668530f3d8674877af456e7d3d3f13d290f4b804b90645",
         "0301ddda1e31cb38c2b11a6772ba0c0f74fb0edcd1c92b490ba62803b2a054bc04",
         "0397380e865977ff77123c8f7f89d8d6f6dd3baef854a2253c7544200baf0d9372",
         "02b6da8f38f622b34a5cc0193b01fd7c9fbcfc631f867a81540410e337f7f12944",
         "023ecc0e7103a522e0dab399e40bd53be123b8f7eccfbf3a86ad84224aeed4b132",
         "02f8c7c0bdc8db04d33a78d662a0b3d584942e778bd6ff54df947e4403b54340a5"
      ],
      "SigData": [
         "3a88a57915f1219472a55f1bcb43837ef351f869511b8723cedb2a0f584d568c57ffd6758698da67310258b5e962ec0ed4a3408635ab05fa264da0fc41e24666",
         "abaecf69bdb839a178a321b86a2f35f0e0ee7a89bca3548ebdfa61c76e332d1fe721d76641c3dc8c11f2aee5fdb3bccc2203612b1c2b3f9573845a8ab1c07ce1",
         "4e52ac5729202edb72c802365451b1d25a7c475a98914b67d0c78d7e93ceb07a89cb73d29155bd3df9df70256a478ca92dcce8253539824256860f7b92803844",
         "005f286a0e232298294781bbc97580bd470a79f3f805c6417403836fdda8551a659c4a4b4f34e4d181accf3c9951f6510b462baad85a77b8da3b48c0c8e9c69f",
         "420cb8b17da30cbb9dc05fd3d18ecc0ce62110f3393b5c58f37e80197313103b674785f6e4069b78df713c9f855cde5c3302384478ec55478b206d63cdaeb074",
         "fe7f93c35a186440de6f77410693bfffae9c34a07d2bd0087ac4a91bb1d882872c0d6d205c7291a5f69915feec60f502b4cd40d679bec5783b5572b01ffb33cf"
      ],
      "Hash": "ea0720ca0ad21b408136bf233c4a9e11be26e8d26cee676c23d8463abbc0200a"
   },
   "Transactions": []
}
```

## 官方社区

### 官方网站

- https://zeepin.io/
- https://twitter.com/ZeepinChain
- https://medium.com/@zeepin
- https://www.reddit.com/r/ZEEPIN/
- https://www.facebook.com/ZeepinChain/

### Telegram

- https://t.me/zeepin
- https://t.me/ZeepinNews


## 许可证

zeepin遵守GNU Lesser General Public License, 版本3.0。 详细信息请查看项目根目录下的LICENSE文件。
