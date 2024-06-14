VERSION:=$(shell cat ./VERSION)
BINARY_NAME:=$(shell cat ./BINARY_NAME)
CONFIG_DIR:=${HOME}/.config/${BINARY_NAME}-${VERSION}
SHARED_DIR:=${HOME}/.local/share/${BINARY_NAME}-${VERSION}
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

make_bin_shared: 
	cp $(shell pwd)/bin/${BINARY_NAME} ${SHARED_DIR}/${BINARY_NAME};
	@echo "The ${BINARY_NAME} binary file can be found in ${SHARED_DIR}";
	@echo ""
	bash ./scripts/setup_aliases.sh;
	@echo "Added ${BINARY_NAME} aliases to .bashrc|.bash_aliases|.zshrc|.zsh_aliases";
	@echo "you can run the program by using this command: ${BINARY_NAME}"

setup: 
	@echo $(git tag | tail -1) > VERSION 
	@echo ""
	@echo "Setting up the config and local shared directories, and the appropriate files."
	bash ./scripts/setup_config_dir.sh 
	bash ./scripts/setup_shared_dir.sh 

	bash ./scripts/setup_about_file.sh

install: setup make_bin_shared  

rm_local_bin: 
	rm -rf ${SHARED_DIR}/${BINARY_NAME}

uninstall: rm_local_bin 
	rm -rf ${CONFIG_DIR}
	rm -rf ${SHARED_DIR}