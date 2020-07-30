package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	BankObjType = "bank"
)

type TransferChainCode struct{}

func (t *TransferChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	bank1 := NewBank(BANK001, nil)
	bank2 := NewBank(BANK002, nil)
	t.PutBank(stub, bank1)
	t.PutBank(stub, bank2)
	return shim.Success([]byte("ok"))
}

func (t *TransferChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success([]byte("ok"))
}

func main() {
	err := shim.Start(new(TransferChainCode))
	if err != nil {
		log.Criticalf("Error starting Simple chaincode: %s", err.Error())
	}
}

// 利用Name创建复合键，把银行对象保存到db
func (t *TransferChainCode) PutBank(stub shim.ChaincodeStubInterface, bank *Bank) error {

	ck, err := stub.CreateCompositeKey(BankObjType, []string{bank.Name})
	if err != nil {
		log.Criticalf(err.Error())
		return fmt.Errorf("create key error: %s", err.Error())
	}

	data, err := json.Marshal(bank)
	if err != nil {
		log.Criticalf(err.Error())
		return fmt.Errorf("marshal bank error: %s", err.Error())
	}

	if err = stub.PutState(ck, data); err != nil {
		log.Critical(err.Error())
		return fmt.Errorf("save bank to db error: %s", err.Error())
	}

	return nil
}

// 利用Name创建复合键，从db读取银行对象
func (t *TransferChainCode) GetBank(stub shim.ChaincodeStubInterface, bankName string) (*Bank, error) {

	ck, err := stub.CreateCompositeKey(BankObjType, []string{bankName})
	if err != nil {
		log.Criticalf(err.Error())
		return nil, fmt.Errorf("create key error: %s", err.Error())
	}

	data, err := stub.GetState(ck)
	if err != nil {
		log.Critical(err.Error())
		return nil, fmt.Errorf("get bank from db error: %s", err.Error())
	}

	var bank Bank
	if err = json.Unmarshal(data, &bank); err != nil {
		log.Critical(err.Error())
		return nil, fmt.Errorf("unmarshal bank error: %s", err.Error())
	}
	return &bank, nil
}
