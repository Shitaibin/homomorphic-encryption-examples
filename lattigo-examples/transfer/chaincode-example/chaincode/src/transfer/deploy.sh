#!/bin/bash

# 更新vendor
go mod vendor

# copy to fabric-samples/chaincodes
rm -rf /Users/shitaibin/go/src/github.com/hyperledger/fabric/fabric-samples/chaincode/transfer
cp -r ../transfer /Users/shitaibin/go/src/github.com/hyperledger/fabric/fabric-samples/chaincode/transfer

# enter docker container and deploy chaincode
# import byfn/scripts/script.sh
docker exec cli /opt/gopath/src/github.com/chaincode/transfer/script.sh $CHANNEL_NAME $CLI_DELAY $LANGUAGE $CLI_TIMEOUT $VERBOSE $NO_CHAINCODE
if [ $? -ne 0 ]; then
  echo "ERROR !!!! Deploy Transfer chaincode failed"
  exit 1
fi