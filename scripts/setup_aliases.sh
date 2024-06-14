#!/bin/bash

SH=$(echo $SHELL)

if [[ "$SH" == "/bin/bash"  ]]; then 
    if [ -f $HOME/.bash_aliases ]; then 
	    echo "alias $(cat ./BINARY_NAME)='$HOME/.local/share/$(cat ./BINARY_NAME)-$(cat ./VERSION)/$(cat ./BINARY_NAME)'" >> $HOME/.bash_aliases;
      source $HOME/.bash_aliases;
    else 
	    echo "alias $(cat ./BINARY_NAME)='$HOME/.local/share/$(cat ./BINARY_NAME)-$(cat ./VERSION)/$(cat ./BINARY_NAME)'" >> $HOME/.bashrc;
      source $HOME/.bashrc;
    fi 
fi 


if [[ "$SH" == "/bin/zsh"  ]]; then 
    if [ -f $HOME/.zsh_aliases ]; then 
	    echo "alias $(cat ./BINARY_NAME)='$HOME/.local/share/$(cat ./BINARY_NAME)-$(cat ./VERSION)/$(cat ./BINARY_NAME)'" >> $HOME/.zsh_aliases;
      source $HOME/.zsh_aliases; 
    else 
	    echo "alias $(cat ./BINARY_NAME)='$HOME/.local/share/$(cat ./BINARY_NAME)-$(cat ./VERSION)/$(cat ./BINARY_NAME)'" >> $HOME/.zshrc;
      source $HOME/.zshrc;
    fi 
fi 

