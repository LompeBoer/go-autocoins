BIN_DIR=./bin
BIN=go-autocoins
BIN_WINDOWS=go-autocoins.exe
WINDOWS_OS=windows
LINUX_OS=linux
MAC_OS=darwin
ARCH=amd64
CMD=cmd/autocoins/*.go

.PHONY: build bin-dir clean run test release

build: bin-dir
	go build -o $(BIN_DIR)/$(BIN) $(CMD); \

build-darwin: bin-dir
	GOOS=$(MAC_OS) GOARCH=$(ARCH) go build -o $(BIN_DIR)/$(BIN) $(CMD); \
	tar -czvf $(BIN_DIR)/$(BIN).$(MAC_OS)-$(ARCH).tar.gz $(BIN_DIR)/$(BIN); \
	rm $(BIN_DIR)/$(BIN); \

build-windows: bin-dir
	GOOS=$(WINDOWS_OS) GOARCH=$(ARCH) go build -o $(BIN_DIR)/$(BIN_WINDOWS) $(CMD); \
	zip -9 -y $(BIN_DIR)/$(BIN).$(WINDOWS_OS)-$(ARCH).zip $(BIN_DIR)/$(BIN_WINDOWS); \
	rm $(BIN_DIR)/$(BIN_WINDOWS); \

build-linux: bin-dir
	GOOS=$(LINUX_OS) GOARCH=$(ARCH) go build -o $(BIN_DIR)/$(BIN) $(CMD); \
	tar -czvf $(BIN_DIR)/$(BIN).$(LINUX_OS)-$(ARCH).tar.gz $(BIN_DIR)/$(BIN); \
	rm $(BIN_DIR)/$(BIN); \

bin-dir:
	mkdir -p $(BIN_DIR)

no-bin-dir:
	rm -rf $(BIN_DIR)

clean:
	rm -rf $(BIN_DIR)

run:
	go run cmd/autocoins/*.go

test:
	go test -v ./...

release: build
	VERSION=$$($(BIN_DIR)/$(BIN) --version); \
	git tag -s -F CHANGELOG.md $$VERSION; \
	git tag -v $$VERSION;
