# Binary name
BINARY=go-bet

# Builds the project
build:
	go build -o ${BINARY}
	go test -v

# Installs our project: copies binaries
install:
	go install
release:
	# Build for windows
	go clean
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ${BINARY}_windows_386.exe

	# Build for linux
	go clean
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o  ${BINARY}_linux_386

	# Build for mac
	go clean
	go build -o ${BINARY}_darwin_arm64

	go clean

# Cleans our projects: deletes binaries
clean:
	go clean

.PHONY:  clean build