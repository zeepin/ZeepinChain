
<h1 align="center">zeepin </h1>
<p align="center" class="version">Version 1.0.0 </p>

[![GoDoc](https://godoc.org/github.com/zeepin/ZeepinChain?status.svg)](https://godoc.org/github.com/zeepin/ZeepinChain)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeepin/ZeepinChain)](https://goreportcard.com/report/github.com/zeepin/ZeepinChain)
[![Travis](https://travis-ci.org/zeepin/ZeepinChain.svg?branch=master)](https://travis-ci.org/zeepin/ZeepinChain)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/zeepin/ZeepinChain?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)


English | [中文](install_CN.md) 
## Build development environment
The requirements to build zeepin are:

- Golang version 1.9 or later
- Glide (a third party package management tool)
- Properly configured Go language environment
- Golang supported operating system

## Deployment|Get zeepin
### Get from source code

Clone the zeepin repository into the appropriate $GOPATH/src/github.com/zeepin directory.

```
$ git clone https://github.com/zeepin/ZeepinChain.git
```
or
```
$ go get github.com/zeepin/ZeepinChain
```
Fetch the dependent third party packages with glide.

```
$ cd $GOPATH/src/github.com/zeepin/ZeepinChain
$ glide install
```

Build the source code with make.

```
$ make all
```

After building the source code sucessfully, you should see two executable programs:

- `zeepin`: the node program/command line program for node control
- `tools/sigsvr`: (optional)zeepin Signature Server - sigsvr is a rpc server for signing transactions for some special requirement. Detail docs can be reference at [link](sigsvr.md).

### get from release
- You can download latest zeepin binary file with  `curl https://dev.zeepin.io/ZeepinChain_install | sh` .

- You can download other version at [release page](https://github.com/zeepin/ZeepinChain/releases).