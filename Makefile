# The name of the executable (default is current directory name)
TARGET := `basename ${PWD}`
.DEFAULT_GOAL: build

ARCH ?= amd64
GOARCH ?= ${ARCH}
GOARM ?=

# These will be provided to the target
VERSION := 1.0.0
BUILD := `git rev-parse HEAD`

HUBNAME=opny

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-s -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: deps api build clean install uninstall run docker/build docker/push run-agent

all: deps api/client install
	@echo "Done"
deps:
	go get -u -f ./...

	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u github.com/go-openapi/runtime
	go get -u golang.org/x/net/context
	go get -u golang.org/x/net/context/ctxhttp

api:
	protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=google/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. api/api.proto
	protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. api/api.proto
	protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. api/api.proto

api/client: api
	swagger generate client -t ./api -f api/api.swagger.json

build:
	ARCH=${ARCH} GOARCH=${GOARCH} GOARM=${GOARM} go build $(LDFLAGS) -o $(TARGET)

clean:
	@rm -f $(TARGET)

install: build
	@go install $(LDFLAGS)

uninstall: clean
	@rm -f $$(which ${TARGET})

run: install
	@$(TARGET) $@

docker/build: build
	@docker build . -t $(HUBNAME)/$(TARGET)

docker/push: docker/build
	@docker push $(HUBNAME)/$(TARGET)

docker/clean:
	@docker rmi $(docker images | grep ${TARGET} | awk '{print $1}')
