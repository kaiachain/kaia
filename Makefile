# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

GO ?= latest
GOPATH := $(or $(GOPATH), $(shell go env GOPATH))
GORUN = env GOPATH=$(GOPATH) GO111MODULE=on go run

BIN = $(shell pwd)/build/bin
BUILD_PARAM?=install
FIXTURE_TEST_REGEX?=^(TestExecutionSpecBlockTestSuite|TestExecutionSpecStateTestSuite|TestKaiaSpecState|TestRLP|TestTransaction|TestVM)$$

OBJECTS=kcn ken kscn ksen kbn kgen homi

.PHONY: all submodules test clean ${OBJECTS}

all: ${OBJECTS}

${OBJECTS}:
ifeq ($(USE_ROCKSDB), 1)
	$(GORUN) build/ci.go ${BUILD_PARAM} -tags rocksdb ./cmd/$@
else
	$(GORUN) build/ci.go ${BUILD_PARAM} ./cmd/$@
endif

abigen:
	$(GORUN) build/ci.go ${BUILD_PARAM} ./cmd/abigen
	@echo "Done building."
	@echo "Run \"$(BIN)/abigen\" to launch abigen."

test:
	$(GORUN) build/ci.go test

test-seq:
	$(GORUN) build/ci.go test -p 1

test-datasync test-networks test-node test-storage test-tests:
	$(GORUN) build/ci.go test -p 1 ./$(patsubst test-%,%,$@)/...

test-eestfixture:
	$(GORUN) build/ci.go test -ensure-fixtures -p 1 -run '$(FIXTURE_TEST_REGEX)' ./tests/...

test-others:
	$(GORUN) build/ci.go test -p 1 -exclude github.com/kaiachain/kaia/datasync,github.com/kaiachain/kaia/networks,github.com/kaiachain/kaia/node,github.com/kaiachain/kaia/storage,github.com/kaiachain/kaia/tests

cover:
	$(GORUN) build/ci.go cover -p 1 -coverprofile=coverage.out
	go tool cover -func=coverage.out -o coverage_report.txt
	go tool cover -html=coverage.out -o coverage_report.html
	@echo "Two coverage reports coverage_report.txt and coverage_report.html are generated."

lint:
	$(GORUN) build/ci.go lint

lint-try:
	$(GORUN) build/ci.go lint-try

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(BIN)/* build/_workspace/src/

# The devtools target installs tools required for 'go generate'.
# You need to put $BIN (or $GOPATH/bin) in your PATH to use 'go generate'.

submodules:
	git submodule update --init

devtools:
	env GOFLAGS= GOBIN= go install golang.org/x/tools/cmd/stringer@latest
	env GOFLAGS= GOBIN= go install github.com/go-bindata/go-bindata/...@latest
	env GOFLAGS= GOBIN= go install mvdan.cc/gofumpt@latest
	env GOFLAGS= GOBIN= go install golang.org/x/tools/cmd/goimports@latest
	env GOFLAGS= GOBIN= go install github.com/fjl/gencodec@latest
	env GOFLAGS= GOBIN= go install github.com/golang/mock/mockgen@latest
	env GOFLAGS= GOBIN= go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	env GOFLAGS= GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
