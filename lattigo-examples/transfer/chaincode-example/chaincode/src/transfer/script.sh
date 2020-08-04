#!/bin/bash


. import /Users/shitaibin/go/src/github.com/hyperledger/fabric/fabric-samples/first-network/byfn.sh

echo "Install Transfer chaincode as mycc-2.0"
VERSION=2.0
CC_SRC_PATH=github.com/chaincode/transfer
installChaincode 0 1



