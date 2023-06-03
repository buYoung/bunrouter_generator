fileName = bunrouter_generator
addFlags = -ldflags "-s -w"
outputPath = ./bin

build:
	GOOS=darwin GOARCH=arm64 go build $(addFlags) -o $(outputPath)/$(fileName) main.go
	GOOS=windows GOARCH=amd64 go build $(addFlags) -o $(outputPath)/$(fileName).exe main.go