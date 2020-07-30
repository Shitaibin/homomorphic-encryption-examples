package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func main() {
	err := shim.Start(new(TransferChainCode))
	if err != nil {
		log.Criticalf("Error starting Simple chaincode: %s", err.Error())
	}
}

const (
	BankObjType = "bank"
)

type TransferChainCode struct{}

func (t *TransferChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	accounts := NewBankAccounts()
	for _, ac := range accounts {
		if err := t.PutAccount(stub, ac); err != nil {
			return shim.Error(err.Error())
		}
	}
	return shim.Success([]byte("ok"))
}

func (t *TransferChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()

	switch f {
	case "SetAccountBalance":
		if err := t.SetAccountBalance(stub, args[1:]); err != nil {
			return shim.Error(err.Error())
		}
	}
	return shim.Success([]byte("ok"))
}

// SetAccountBalance TODO 设置用户余额
func (t *TransferChainCode) SetAccountBalance(stub shim.ChaincodeStubInterface, args []string) error {
	log.Infof("set balance [%s] [%s]", args[0], args[1])

	// 需要读取account，并更新
	return nil
}

// 利用bankID和account.ID创建复合键，把账户对象保存到db
func (t *TransferChainCode) PutAccount(stub shim.ChaincodeStubInterface, acc *Account) error {

	ck, err := stub.CreateCompositeKey(BankObjType, []string{acc.BankID, acc.ID})
	if err != nil {
		log.Criticalf(err.Error())
		return fmt.Errorf("create key error: %s", err.Error())
	}

	data, err := json.Marshal(acc)
	if err != nil {
		log.Criticalf(err.Error())
		return fmt.Errorf("marshal account error: %s", err.Error())
	}

	if err = stub.PutState(ck, data); err != nil {
		log.Critical(err.Error())
		return fmt.Errorf("save account to db error: %s", err.Error())
	}

	return nil
}

// 利用bankID和account.ID创建复合键，从db读取账户
func (t *TransferChainCode) GetAccount(stub shim.ChaincodeStubInterface, bankID string, accID string) (*Account, error) {

	ck, err := stub.CreateCompositeKey(BankObjType, []string{bankID, accID})
	if err != nil {
		log.Criticalf(err.Error())
		return nil, fmt.Errorf("create key error: %s", err.Error())
	}

	data, err := stub.GetState(ck)
	if err != nil {
		log.Critical(err.Error())
		return nil, fmt.Errorf("get account from db error: %s", err.Error())
	}

	var acc Account
	if err = json.Unmarshal(data, &acc); err != nil {
		log.Critical(err.Error())
		return nil, fmt.Errorf("unmarshal account error: %s", err.Error())
	}
	return &acc, nil
}
