.PHONY: deps build


deps:
	go get ./...

build: deps
	@mkdir -p release/
	go build -o release/oscar main/oscar.go
	CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -a -o release/oscar_linux main/oscar.go

examples:
	go run main/oscar.go run example/base64.lua example/httpbin.lua example/fail.lua example/failInit.lua example/common.lua -e example/env.ini -l example/wrapper.lua -j dev/report.json  --html-report dev/examples-report
