
build:
	go build -o gladia -v ./cmd/*.go

dist:
	GOOS=linux GOARCH=arm GOARM=7 go build -o dist/gladia-linux-armhf -v ./cmd/*.go
	GOOS=linux GOARCH=arm64 go build -o dist/gladia-linux-arm64 -v ./cmd/*.go
	GOOS=linux GOARCH=amd64 go build -o dist/gladia-linux-x86_64 -v ./cmd/*.go
	GOOS=linux GOARCH=386 go build -o dist/gladia-linux-i386 -v ./cmd/*.go
	GOOS=windows GOARCH=amd64 go build -o dist/gladia-windows-x86_64.exe -v ./cmd/*.go
	GOOS=darwin GOARCH=amd64 go build -o dist/gladia-darwin-x86_64 -v ./cmd/*.go
	GOOS=darwin GOARCH=arm64 go build -o dist/gladia-darwin-arm64 -v ./cmd/*.go

dev:
	go run -x ./cmd/*.go -audio-file split_infinity.wav
watch-dev:
	go run -x ./cmd/*.go -audio-file split_infinity.wav

test:
	go test -race -v ./...
watch-test:
	reflex -t 50ms -s -- sh -c 'gotest -race -v ./...'

bench:
	go test -benchmem -benchtime=10000000x -bench=. ./...
watch-bench:
	reflex -t 50ms -s -- sh -c 'go test -benchmem -benchtime=10000000x -bench=. ./...'

coverage:
	go test -v -coverprofile=cover.out -covermode=atomic ./...
	go tool cover -html=cover.out -o cover.html

# tools
tools:
	go install github.com/cespare/reflex@latest
	go install github.com/rakyll/gotest@latest
	go install github.com/psampaz/go-mod-outdated@latest
	go install github.com/jondot/goweight@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go get -t -u golang.org/x/tools/cmd/cover
	go get -t -u github.com/sonatype-nexus-community/nancy@latest
	go mod tidy

lint:
	golangci-lint run --timeout 60s --max-same-issues 50 ./...
lint-fix:
	golangci-lint run --timeout 60s --max-same-issues 50 --fix ./...

audit: tools
	go mod tidy
	go list -json -m all | nancy sleuth

outdated: tools
	go mod tidy
	go list -u -m -json all | go-mod-outdated -update -direct

weight: tools
	goweight
