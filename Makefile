VERSION:=$(shell cat ./VERSION)
BINARY_NAME:=$(shell cat ./BINARY_NAME)
AUTHOR:="Massiles Ghernaout"

info: 
	@echo "Project: ${BINARY_NAME}@${VERSION}"
	@echo "Author: ${AUTHOR}"

clean:
	rm -rf bin/*

bin: 
	go build -o bin/${BINARY_NAME} cmd/main.go 

binstatic:
	rm -rf bin/*
	@echo "Building a static executable..."
	CGO_ENABLED=0 go build -a -tags netgo,osusergo -ldflags "-X main.version=${VERSION} -X main.binary_name=${BINARY_NAME} -extldflags '-static -s -w'" -o bin/${BINARY_NAME} cmd/main.go

run: 
	./bin/${BINARY_NAME}

runsrc:
	ENV=dev DEBUG=true go run ./cmd/main.go

