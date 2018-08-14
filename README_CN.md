
<h1 align="center">zeepin </h1>
<h4 align="center">Version 0.1 </h4>

[English](README.md) | [中文](README_CN.md) | [한글](README_KO.md)


欢迎查看zeepin的源码库!

Zeepin Chain是一条去中心化的文创及娱乐资产公链，通过区块链构建标准化基础设施，为创意人群提供高效工作的解决方案，帮助创意组织提高创新效率，促进文创产业开放透明、公平高效的价值流通。同时Zeepin Chain还将打造区块链数字娱乐资产发行平台，为全球文娱资产代币化提供区块链技术支持及落地场景的建设。Zeepin Chain公链作为一条行业基础链，拥有整合第三方娱乐资产和系统的能力，建立一个自由的交易市场和兑换平台。

Zeepin Chain构建了完整的区块链技术框架，采用GBFT－POS共识机制（星际共识），提供具备图灵完备性的虚拟机作为智能合约的执行环境，为应用架构提供自定义脚步控制支持。支持Java、C#、Python、Javascript等编程语言开发的脚本，虚拟机都可以通过api与链进行集成交互。

zeepin致力于创建一个可自由配置、高性能、可扩展的区块链底层基础设施，让部署及调用去中心化应用变得更加的简单。目前代码还处于内部测试阶段，但处于快速的迭代开发中，欢迎及希望更多的开发者加入到zeepin中来！


## 目录

* [构建开发环境](#构建开发环境)
* [获取zeepin](#获取zeepin)
    * [从release获取](#从release获取)
* [服务器部署](#服务器部署)
    * [选择网络](#选择网络)
        * [主网同步节点部署](#主网同步节点部署)
        * [公开测试网同步节点部署](#公开测试网同步节点部署)
    * [运行](#运行)
    * [ZPT转账调用示例](#zpt转账调用示例)
* [开源社区](#开源社区)
    * [网站](#网站)
    * [Discord开发者社区](#discord开发者社区)
* [许可证](#许可证)

## 构建开发环境
成功编译zeepin需要以下准备：

* Golang版本在1.9及以上
* 安装第三方包管理工具glide
* 正确的Go语言开发环境
* Golang所支持的操作系统

## 获取zeepin
### 从源码获取
克隆zeepin仓库到 **$GOPATH/src/github.com/zeepin** 目录

```shell
$ git clone https://github.com/zeepin/ZeepinChain.git
```
或者
```shell
$ go get github.com/zeepin/ZeepinChain
```

用第三方包管理工具glide拉取依赖库

````shell
$ cd $GOPATH/src/github.com/zeepin/ZeepinChain
$ glide install
````

用make编译源码

```shell
$ make all
```

成功编译后会生成两个可以执行程序

* `zeepin`: 节点程序/以命令行方式提供的节点控制程序
* `tools/sigsvr`: (可选)签名服务 - sigsvr是一个签名服务的server以满足一些特殊的需求。详细的文档可以在[这里](./docs/specifications/sigsvr_CN.md)参考

### 从release获取
- 你可以通过命令 ` curl https://dev.zeepin.io/ZeepinChain_install | sh ` 获取最新的zeepin版本
- 你也可以从[下载页面](https://github.com/zeepin/zeepin/releases)获取.

## 服务器部署
### 选择网络
zeepin的运行支持以下方式

* 主网同步节点部署
* 公开测试网同步节点部署

#### 主网同步节点部署

直接启动zeepin

   ```
	./zeepin --networkid 1
   ```

#### 公开测试网同步节点部署

直接启动zeepin

   ```
	./zeepin --networkid 2
   ```

#### 单机部署配置

在单机上创建一个目录，在目录下存放以下文件：
- 节点程序 + 节点控制程序 `zeepin`
- 钱包文件`wallet.dat`

使用命令 `$ ./zeepin --testmode --networkid 3` 即可启动单机版的测试网络。

单机配置的例子如下：
- 目录结构

    ```shell
    $ tree
    └── node
        ├── zeepin
        └── wallet.dat
    ```

#### 多机部署配置

注意：当你想搭建一个私网去运行GBFT模式的zeepin时，你必须要使用 --config参数去指定一个具体的配置文件，并且要使用 --networkid去定义一个
网络标识（不能等于1/2/3），否则，zeepin节点将会自动连接到主网上。

##### GBFT部署方法

多机环境下，根据GBFT共识算法的要求，最少需要7个节点（共识节点）完成部署。

我们可以通过修改配置文件[`config-gbft.json`](./docs/specifications/config-gbft.json)进行快速部署，配置文件的说明请点击[这里](./docs/specifications/config_CN.md)。

1. 生成七个钱包文件，每个钱包文件包含一个账户，共七个账户，分别作为每个节点的记账人。生成账户和钱包的命令为：
	```
	./zeepin account add -d -w wallet.dat
	Use default setting '-t ecdsa -b 256 -s SHA256withECDSA' 
		signature algorithm: ecdsa 
		curve: P-256 
		signature scheme: SHA256withECDSA 
	Password:
	Re-enter Password:

	Index: 1
	Label: 
	Address: AXkDGfr9thEqWmCKpTtQYaazJRwQzH48eC
	Public key: 03d7d8c0c4ca2d2bc88209db018dc0c6db28380d8674aff86011b2a6ca32b512f9
	Signature scheme: SHA256withECDSA

	Create account successfully.
	```
	使用-w参数指定生成的钱包文件名

2. 修改`config-gbft.json`，将上一步生成的七个账户的公钥、address分别填入`config-gbft.json`中的peers配置的1-7项

3. 将相关文件复制到目标主机，包括：
    - 配置文件`config-gbft.json`
    - 节点程序`zeepin`
    - 钱包文件

4. 设置每个节点网络连接的端口号（推荐不做修改，使用默认端口配置）
    - `NodePort`为的P2P连接端口号（默认20338）
    - `HttpJsonPort`和`HttpLocalPort`为RPC端口号（默认为20336，20337）

5. 种子节点配置
    - 在7个主机中选出至少一个做种子节点，并将种子节点地址分别填写到每个配置文件的`SeelList`中，格式为`种子节点IP地址 + 种子节点NodePort`



多机部署配置完成，每个节点目录结构如下:

   ```shell
	$ ls
	config.json  zeepin  wallet.dat
   ```
### 运行

以任意顺序运行每个节点node程序， 并在出现`Password:`提示后输入节点的钱包密码。

如果你想启动一个共识节点（比如在自己的私网上），必须使用 --enableconsensus参数。如果你想自己运行一个私网，需要使用 --networkid参数指定私网节
点之间通信的网络id（不能等于1/2/3），并且要使用 --config参数指定配置文件。

例如:
   ```
    $ ./zeepin --enableconsensus --networkid 4 --config ./config.json
    $ - Input your wallet password
   ```

了解更多请运行 `./zeepin --help`


### ZPT转账调用示例
   - from: 转出地址； - to: 转入地址； - amount: 资产转移数量；
      from参数可以不指定，如果不指定则使用默认账户。

```shell
  ./zeepin asset transfer  --to=TA4Xe9j8VbU4m3T1zEa1uRiMTauiAT88op --amount=10
```

执行完后会输出：

```
Transfer ZPT
From:TA6edvwgNy3c1nBHgmFj8KrgQ1JCJNhM3o
To:TA4Xe9j8VbU4m3T1zEa1uRiMTauiAT88op
Amount:10
TxHash:10dede8b57ce0b272b4d51ab282aaf0988a4005e980d25bd49685005cc76ba7f
```
其中TxHash是转账交易的交易HASH，可以通过这个HASH查询转账交易的直接结果。
出于区块链出块时间的限制，提交的转账请求不会马上执行，需要等待至少一个区块时间，等待记账节点打包交易。

### 查询转账结果示例

--hash:指定查询的转账交易hash
```shell
./zeepin asset status --hash=10dede8b57ce0b272b4d51ab282aaf0988a4005e980d25bd49685005cc76ba7f
```
查询结果：
```shell
Transaction:transfer success
From:TA6edvwgNy3c1nBHgmFj8KrgQ1JCJNhM3o
To:TA4Xe9j8VbU4m3T1zEa1uRiMTauiAT88op
Amount:10
```

### 查询账户余额示例

--address:账户地址

```shell
./zeepin asset balance --address=TA4Xe9j8VbU4m3T1zEa1uRiMTauiAT88op
```
查询结果：
```shell
BalanceOf:TA4Xe9j8VbU4m3T1zEa1uRiMTauiAT88op
ZPT:10
GALA:0
GALAApprove:0
```


## 开源社区

### 网站

- https://zeepin.io/

### Discord开发者社区

- https://discord.gg/

## 许可证

zeepin遵守GNU Lesser General Public License, 版本3.0。 详细信息请查看项目根目录下的LICENSE文件。
