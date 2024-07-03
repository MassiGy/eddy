BINARY_NAME:=$(shell cat ./BINARY_NAME)
AUTHOR:="Massiles Ghernaout"

TAG=$(shell git tag | tail -1)
ifneq ($(TAG),)
    VERSION:=${TAG}
else 
    VERSION:=$(shell cat ./VERSION)
endif


ifeq ($(OS),Windows_NT)
    MACHINE = windows
    ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
        ARCH = amd64
    endif
    ifeq ($(PROCESSOR_ARCHITECTURE),x86)
        ARCH = 386
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        MACHINE = linux
    endif
    ifeq ($(UNAME_S),Darwin)
        MACHINE = darwin
    endif
    UNAME_P := $(shell arch)
    ifeq ($(UNAME_P),x86_64)
        ARCH = amd64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        ARCH = 386
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        ARCH = arm
    endif
endif




# [ PRODUCTION ]
info: 
	@echo "Project: ${BINARY_NAME}@${VERSION}"
	@echo "Author: ${AUTHOR}"
	@echo "Current OS: ${MACHINE}"
	@echo "Current architecture: ${ARCH}"

binary: 

	rm -rf bin/*
	@echo "Building the executable..."

ifeq ($(MACHINE),Windows_NT)
		GOOS=${MACHINE} GOARCH=${ARCH} \
		go build -ldflags \
		"-X main.version=${VERSION} -X main.binary_name=${BINARY_NAME}.exe -extldflags=-static"\
		-o bin/${BINARY_NAME}.exe \
		cmd/*.go 
else 
		GOOS=${MACHINE} GOARCH=${ARCH} \
        go build -ldflags \
        "-X main.version=${VERSION} -X main.binary_name=${BINARY_NAME} -extldflags=-static -s -w"\
		-o bin/${BINARY_NAME} \
		cmd/*.go 
endif



# [ RELEASES ]

win32_x86_release: 
	rm -rf releases/win32_x86/*

	GOOS=windows GOARCH=386 \
	go build -ldflags \
	"-X main.version=${VERSION} -X main.binary_name=${BINARY_NAME}.exe -extldflags=-static -s -w" \
	-o releases/win32_x86/${BINARY_NAME}.exe \
	cmd/*.go 

	zip -r releases/win32_x86.zip releases/win32_x86/*


win32_amd64_release: 
	rm -rf releases/win32_amd64/*

	GOOS=windows GOARCH=amd64 \
	go build -ldflags \
	"-X main.version=${VERSION} -X main.binary_name=${BINARY_NAME}.exe -extldflags=-static -s -w" \
	-o releases/win32_amd64/${BINARY_NAME}.exe \
	cmd/*.go 

	zip -r releases/win32_amd64.zip releases/win32_amd64/*


linux_x86_release:
	rm -rf releases/linux_x86/*

	GOOS=linux GOARCH=386 \
	go build -ldflags \
	"-X main.version=${VERSION} -X main.binary_name=${BINARY_NAME} -extldflags=-static -s -w" \
	-o releases/linux_x86/${BINARY_NAME} \
	cmd/*.go 

	zip -r releases/linux_x86.zip releases/linux_x86/*

linux_amd64_release:
	rm -rf releases/linux_amd64/*

	GOOS=linux GOARCH=amd64 \
	go build -ldflags \
	"-X main.version=${VERSION} -X main.binary_name=${BINARY_NAME} -extldflags=-static -s -w" \
	-o releases/linux_amd64/${BINARY_NAME} \
	cmd/*.go 

	zip -r releases/linux_amd64.zip releases/linux_amd64/*

# https://stackoverflow.com/questions/65881808/how-can-i-cross-compile-to-darwin-386-from-linux
# darwin/32bit support was dropped since go 1.15, if you want to get this release, use a go version prior to 1.14
osx_x86_release: 
	rm -rf releases/osx_x86/*

	GOOS=darwin GOARCH=386 \
	go build -ldflags \
	"-X main.version=${VERSION} -X main.binary_name=${BINARY_NAME} -extldflags=-static -s -w"\
	-o releases/osx_x86/${BINARY_NAME} \
	cmd/*.go 

	zip -r releases/osx_x86.zip releases/osx_x86/*

osx_amd64_release: 
	rm -rf releases/osx_amd64/*

	GOOS=darwin GOARCH=amd64 \
	go build -ldflags \
	"-X main.version=${VERSION} -X main.binary_name=${BINARY_NAME} -extldflags=-static -s -w" \
	-o releases/osx_amd64/${BINARY_NAME} \
	cmd/*.go 

	zip -r releases/osx_amd64.zip releases/osx_amd64/*

releases: win32_x86_release win32_amd64_release linux_x86_release linux_amd64_release osx_amd64_release

# [ DEVELOPEMENT ]
clean:
	rm -rf bin/*
run: 
	./bin/${BINARY_NAME}
runsrc:
	GOOS=${MACHINE} GOARCH=${ARCH} ENV=dev DEBUG=true go run ./cmd/*
