VERSION=$(shell git describe --tags --always --dirty)

build:
	go build -ldflags="-X 'main.AppVersion=$(VERSION)'" -o dash ./cmd/dash