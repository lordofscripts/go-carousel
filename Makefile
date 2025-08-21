GO=go
GOFLAGS = -v
BIN_OUT=$(GOBIN)/goCarousel
MAIN=cmd/*.go

.PHONY: clean build

build:
	$(GO) build -tags logx $(GOFLAGS) -o ${BIN_OUT} ${MAIN}

release:
	$(GO) build $(GOFLAGS) -o ${BIN_OUT} ${MAIN}

clean:
	go clean

run:
	go run -race  $MAIN

lint: 
	@gofmt -l . | grep ".*\.go"

test:
	go test tests/*test.go	