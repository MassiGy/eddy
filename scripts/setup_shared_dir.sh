#!/bin/bash

if [ ! -d $HOME/.local/share/$(cat ./BINARY_NAME)-$(cat ./VERSION) ]; then
    echo "Creating the $HOME/.local/share/$(cat ./BINARY_NAME)-$(cat ./VERSION) directory";
    mkdir $HOME/.local/share/$(cat ./BINARY_NAME)-$(cat ./VERSION);

    cp ./BINARY_NAME $HOME/.local/share/$(cat ./BINARY_NAME)-$(cat ./VERSION);
    cp ./VERSION $HOME/.local/share/$(cat ./BINARY_NAME)-$(cat ./VERSION);

fi