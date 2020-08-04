package main

import "github.com/hyperledger/fabric/core/chaincode/shim"

// log 全局logger
var log *shim.ChaincodeLogger

func init() {
	if log == nil {
		log = shim.NewLogger("TransferCC")
	}
}
