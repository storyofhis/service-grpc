#!/bin/bash
echo 'export GOPATH=$HOME/Go' >> $HOME/.bashrc
source $HOME/.bashrc

echo 'export PATH=$PATH:$GOPATH/bin' >> $HOME/.bashrc
source $HOME/.bashrc