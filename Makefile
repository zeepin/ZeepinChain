GOFMT=gofmt
GC=go build
VERSION := $(shell git describe --abbrev=4 --always --tags)
BUILD_NODE_PAR = -ldflags "-X github.com/imZhuFei/zeepin/common/config.Version=$(VERSION)" #-race

ARCH=$(shell uname -m)
DBUILD=docker build
DRUN=docker run
DOCKER_NS ?= zeepin
DOCKER_TAG=$(ARCH)-$(VERSION)
ZPT_CFG_IN_DOCKER=config.json
WALLET_FILE=wallet.dat

#SRC_FILES = $(shell git ls-files | grep -e .go$ | grep -v _test.go)
TOOLS=./tools
ABI=$(TOOLS)/abi
NATIVE_ABI_SCRIPT=./cmd/abi/native_abi_script

zeepin: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o zeepin main.go
 
sigsvr: $(SRC_FILES) abi 
	$(GC)  $(BUILD_NODE_PAR) -o sigsvr sigsvr.go
	@if [ ! -d $(TOOLS) ];then mkdir -p $(TOOLS) ;fi
	@mv sigsvr $(TOOLS)

ztools: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o zeepinChainTools ztools.go
ztest: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o ztest ztest.go

ztest-cross: ztest-linux ztest-windows ztest-darwin
ztest-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o ztest-linux-amd64 ztest.go
ztest-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o ztest-windows-amd64.exe ztest.go
ztest-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o ztest-darwin-amd64 ztest.go

abi: 
	@if [ ! -d $(ABI) ];then mkdir -p $(ABI) ;fi
	@cp $(NATIVE_ABI_SCRIPT)/*.json $(ABI)

tools: sigsvr abi

all: zeepin tools

zeepin-cross: zeepin-windows zeepin-linux zeepin-darwin

zeepin-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o zeepin-windows-amd64.exe main.go

zeepin-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o zeepin-linux-amd64 main.go

zeepin-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o zeepin-darwin-amd64 main.go

tools-cross: tools-windows tools-linux tools-darwin

tools-windows: abi 
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o sigsvr-windows-amd64.exe sigsvr.go
	@if [ ! -d $(TOOLS) ];then mkdir -p $(TOOLS) ;fi
	@mv sigsvr-windows-amd64.exe $(TOOLS)

tools-linux: abi 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o sigsvr-linux-amd64 sigsvr.go
	@if [ ! -d $(TOOLS) ];then mkdir -p $(TOOLS) ;fi
	@mv sigsvr-linux-amd64 $(TOOLS)

tools-darwin: abi 
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o sigsvr-darwin-amd64 sigsvr.go
	@if [ ! -d $(TOOLS) ];then mkdir -p $(TOOLS) ;fi
	@mv sigsvr-darwin-amd64 $(TOOLS)

all-cross: zeepin-cross tools-cross abi

format:
	$(GOFMT) -w main.go

$(WALLET_FILE):
	@if [ ! -e $(WALLET_FILE) ]; then $(error Please create wallet file first) ; fi

docker/payload: docker/build/bin/zeepin docker/Dockerfile $(ZPT_CFG_IN_DOCKER) $(WALLET_FILE)
	@echo "Building zeepin payload"
	@mkdir -p $@
	@cp docker/Dockerfile $@
	@cp docker/build/bin/zeepin $@
	@cp -f $(ZPT_CFG_IN_DOCKER) $@/config.json
	@cp -f $(WALLET_FILE) $@
	@tar czf $@/config.tgz -C $@ config.json $(WALLET_FILE)
	@touch $@

docker/build/bin/%: Makefile
	@echo "Building zeepin in docker"
	@mkdir -p docker/build/bin docker/build/pkg
	@$(DRUN) --rm \
		-v $(abspath docker/build/bin):/go/bin \
		-v $(abspath docker/build/pkg):/go/pkg \
		-v $(GOPATH)/src:/go/src \
		-w /go/src/github.com/zeepin/zeepin \
		golang:1.9.5-stretch \
		$(GC)  $(BUILD_NODE_PAR) -o docker/build/bin/zeepin main.go
	@touch $@

docker: Makefile docker/payload docker/Dockerfile 
	@echo "Building zeepin docker"
	@$(DBUILD) -t $(DOCKER_NS)/zeepin docker/payload
	@docker tag $(DOCKER_NS)/zeepin $(DOCKER_NS)/zeepin:$(DOCKER_TAG)
	@touch $@

clean:
	rm -rf *.8 *.o *.out *.6 *exe
	rm -rf zeepin zeepin-* tools docker/payload docker/build
	rm -rf ztest  ztest-*

