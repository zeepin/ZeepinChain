#  候选节点部署流程 

1、候选节点系统配置最低要求：  
    最低4C16G&SSD500G的硬件配置。具体可参考[智品节点招募细则]() 

2、创建目录，下载最新版本zeepin二进制文件至目录下
```
mkdir -p /data/gopath/mainnet
cd /data/gopath/mainnet
```
运行命令: 
```
curl  <https://XXXXXXXXXX>  | sh`
```
   或者可以自行在这里下载：<https://XXXXXXXXXXX/releases> 

3、把生成的钱包文件 wallet.dat拷贝到zeepin二进制目录下，执行命令 ./zeepin account list -v  记录Public key和钱包地址。注意：不能用同一个钱包开启2条链。

4、后台运行 `./zeepin --rest --enableconsensus`  
    ※如果希望指定rest端口可以通过 `--restport 20334`指定，默认为20334端口)
   对应的防火墙策略：20334、20338 面向all。如果是基于云服务，同时还要确保云服务的网络安全配置上打开了20338,20334端口。

   如果后台运行,可以构建如下的shell脚本，命名为start.sh：
```
#! /bin/bash
./zeepin --rest --enableconsensus<<eof
your password
eof >log &
```
然后用以下命令启动
nohup ./start.sh >/dev/null 2>log &

5、验证节点运行
    可以在浏览器中输入, http://your server ip:30334/api/v1/block/height 查看当前节点的区块高度，是否和主网上XXXXXXX  高度保持一致。