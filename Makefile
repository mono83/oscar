.PHONY: deps release test cross

clean: ## Clears environment
	@echo $(shell date +'%H:%M:%S') "\033[0;33mRemoving old release\033[0m"
	@mkdir -p release
	@rm -rf ./assets.go
	@rm -rf release/*

test: ## Runs unit tests
	@echo $(shell date +'%H:%M:%S') "\033[0;32mRunning unit tests\033[0m"
	@CGO_ENABLED=0 go test -tags=http ./...

deps: ## Download required dependencies
	@echo $(shell date +'%H:%M:%S') "\033[0;32mDownloading dependencies\033[0m"
	@go get github.com/stretchr/testify/assert
	@go get ./...

release: clean deps test
	@mkdir -p release/
	go build -o release/oscar main/oscar.go
	go build -o release/oscar_linux main/oscar.go

cross: clean deps test ## Builds cross-OS binaries and run tests
	@echo $(shell date +'%H:%M:%S') "\033[0;32mCompiling Linux version\033[0m"
	@GOOS="linux" GOARCH="amd64" go build -o release/oscar-linux64 main/oscar.go
	@echo $(shell date +'%H:%M:%S') "\033[0;32mCompiling MacOS version\033[0m"
	@GOOS="darwin" GOARCH="amd64" go build -o release/oscar-darwin64 main/oscar.go
	@echo $(shell date +'%H:%M:%S') "\033[0;32mCompiling Windows version\033[0m"
	@GOOS="windows" GOARCH="386" go build -o release/oscar.exe main/oscar.go

examples:
	go run main/oscar.go run example/base64.lua example/httpbin.lua example/fail.lua example/failInit.lua example/common.lua -e example/env.ini -l example/wrapper.lua -j dev/report.json  --html-report dev/examples-report
