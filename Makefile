GOPATH := $(shell pwd)

run: build
	./bin/sup

build: $(shell find . -name '*.go') $(shell find . -name '*.js')
	go install sup
	mkdir -p bin
	go build -o bin/sup src/cmd/sup/*.go

test: $(shell find . -name '*.go') $(shell find . -name '*.js')
	go test sup

lint: 
	go get github.com/golang/lint/golint
	$(GOPATH)/bin/golint *.go

