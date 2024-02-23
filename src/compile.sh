mkdir -p dist

# Linux ARM (32-bit)
GOOS=linux GOARCH=arm GOARM=7 go build -o dist/gladia-linux-armhf

# Linux ARM (64-bit)
GOOS=linux GOARCH=arm64 go build -o dist/gladia-linux-arm64

# Linux 64-bit (x86_64)
GOOS=linux GOARCH=amd64 go build -o dist/gladia-linux-x86_64

# Linux 32-bit (x86)
GOOS=linux GOARCH=386 go build -o dist/gladia-linux-i386

# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o dist/gladia-windows-x86_64.exe

# macOS amd 64-bit (x86_64)
GOOS=darwin GOARCH=amd64 go build -o dist/gladia-darwin-x86_64

# macOS arm 64-bit
GOOS=darwin GOARCH=arm64 go build -o dist/gladia-darwin-arm64
