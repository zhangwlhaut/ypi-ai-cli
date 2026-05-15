.PHONY: build install test clean

BINARY=yapi
GO=go

build:
	$(GO) build -o $(BINARY) ./cmd/yapi/

install:
	$(GO) install ./cmd/yapi/

test:
	$(GO) test ./...

clean:
	rm -f $(BINARY)
