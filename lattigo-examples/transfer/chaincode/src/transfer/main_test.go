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

	if _, err := cc.GetBank(stub, BANK001); err != nil {
		t.Error(err.Error())
	}
	if bank, err := cc.GetBank(stub, BANK002); err != nil {
		t.Error(err.Error())
	} else if bank == nil {
		t.Error("bank is nil")
	}
}
