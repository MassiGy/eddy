#!/bin/bash

if [  ! -d  $HOME/.config/$(cat ./BINARY_NAME)-$(cat ./VERSION) ]; then
    echo "Creating $HOME/.config/$(cat ./BINARY_NAME)-$(cat ./VERSION) directory";
    mkdir $HOME/.config/$(cat ./BINARY_NAME)-$(cat ./VERSION); 

    cp ./BINARY_NAME $HOME/.config/$(cat ./BINARY_NAME)-$(cat ./VERSION);
    cp ./VERSION $HOME/.config/$(cat ./BINARY_NAME)-$(cat ./VERSION);
fi 