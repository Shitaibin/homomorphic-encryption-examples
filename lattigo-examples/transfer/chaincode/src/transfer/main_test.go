package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"gotest.tools/assert"
)

func TestTransferChainCode_Init(t *testing.T) {
	// 创建stub
	var cc = new(TransferChainCode)
	stub := shim.NewMockStub("TransferChainCode", cc)

	// 调用Init
	res := stub.MockInit("1", [][]byte{[]byte("init")})
	assert.Equal(t, int32(shim.OK), res.Status, "Init failed: %s", string(res.Message))

	for _, bid := range banks {
		for _, acid := range accs {
			if _, err := cc.GetAccount(stub, bid, acid); err != nil {
				t.Error(err.Error())
			}
		}
	}
}
