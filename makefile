# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get
BINARY_NAME=color-blender

build_linux: test
	$(GOBUILD) -o $(BINARY_NAME) -v
build_win: test
	$(GOBUILD) -o $(BINARY_NAME).exe -v
test:
	$(GOTEST) -v -cover ./...
coverage:
	-$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOTOOL) cover -html=coverage.out
