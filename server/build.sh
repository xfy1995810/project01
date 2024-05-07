#!/bin/bash
#

if [[ $1 == "clean" ]];then
    rm -rf embed/dist/*
    exit 0
fi

go build -trimpath -ldflags "-w"
upx -9 -k  dcss
rm -f dcss.*
