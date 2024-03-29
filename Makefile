#
# 查看支持的平台
#go tool dist list

# 编译到 Linux
.PHONY: build-linux-amd64
build-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/vacproxy-linux-amd64/vacproxy-linux-amd64
.PHONY: build-linux-arm64
build-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./build/vacproxy-linux-arm64/vacproxy-linux-arm64

# 编译到 macOS
.PHONY: build-darwin-amd64
build-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/vacproxy-darwin-amd64/vacproxy-darwin-amd64
.PHONY: build-darwin-arm64
build-darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./build/vacproxy-darwin-arm64/vacproxy-darwin-arm64

# 编译到 windows
.PHONY: build-windows-amd64
build-windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/vacproxy-windows-amd64/vacproxy-windows-amd64.exe
.PHONY: build-windows-arm64
build-windows-arm64:
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./build/vacproxy-windows-arm64/vacproxy-windows-arm64.exe

# 编译到 全部平台
.PHONY: build-all
build-all:
	make clean
	mkdir -p ./build
	mkdir -p ./build/vacproxy-linux-amd64
	cp config.yaml ./build/vacproxy-linux-amd64
	make build-linux-amd64
	zip -j -q -r -o ./build/vacproxy-linux-amd64.zip ./build/vacproxy-linux-amd64

	mkdir -p ./build/vacproxy-linux-arm64
	cp config.yaml ./build/vacproxy-linux-arm64
	make build-linux-arm64
	zip -j -q -r -o ./build/vacproxy-linux-arm64.zip ./build/vacproxy-linux-arm64

	mkdir -p ./build/vacproxy-darwin-amd64
	cp config.yaml ./build/vacproxy-darwin-amd64
	make build-darwin-amd64
	zip -j -q -r -o ./build/vacproxy-darwin-amd64.zip ./build/vacproxy-darwin-amd64

	mkdir -p ./build/vacproxy-darwin-arm64
	cp config.yaml ./build/vacproxy-darwin-arm64
	make build-darwin-arm64
	zip -j -q -r -o ./build/vacproxy-darwin-arm64.zip ./build/vacproxy-darwin-arm64

	mkdir -p ./build/vacproxy-windows-amd64
	cp config.yaml ./build/vacproxy-windows-amd64
	make build-windows-amd64
	zip -j -q -r -o ./build/vacproxy-windows-amd64.zip ./build/vacproxy-windows-amd64

	mkdir -p ./build/vacproxy-windows-arm64
	cp config.yaml ./build/vacproxy-windows-arm64
	make build-windows-arm64
	zip -j -q -r -o ./build/vacproxy-windows-arm64.zip ./build/vacproxy-windows-arm64

.PHONY: clean
clean:
	rm -rf ./build
