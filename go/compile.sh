mkdir -p dist

# Linux ARM (32-bit)
GOOS=linux GOARCH=arm GOARM=7 go build -o dist/gladia-linux-armv7

# Linux ARM (64-bit)
GOOS=linux GOARCH=arm64 go build -o dist/gladia-linux-arm64

# Linux 64-bit
GOOS=linux GOARCH=amd64 go build -o dist/gladia-linux-amd64

# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o dist/gladia-windows-amd64.exe

# macOS amd 64-bit
GOOS=darwin GOARCH=amd64 go build -o dist/gladia-darwin-amd64

# macOS arm 64-bit
GOOS=darwin GOARCH=arm64 go build -o dist/gladia-darwin-arm64
