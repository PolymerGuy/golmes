# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=./golmes
BINART_FOLDER_NAME=./
BINARY_UNIX=$(BINARY_NAME)_unix
ABQ_MOCK_BINARY = ./examples/dummyApp/abq
INPUT_FILE = examples/example_two_args.yml


all: test build run
build: 
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_NAME).go
test:
	$(GOBUILD) -o $(ABQ_MOCK_BINARY) $(ABQ_MOCK_BINARY).go  
	mv $(ABQ_MOCK_BINARY) $(BINART_FOLDER_NAME)
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(ABQ_MOCK_BINARY) $(ABQ_MOCK_BINARY).go 
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_NAME).go
	$(BINARY_NAME) run $(INPUT_FILE)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v 
