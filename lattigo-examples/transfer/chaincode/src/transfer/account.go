package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

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

	if data == nil {
		return nil, fmt.Errorf("no account: %s %s", bankID, accID)
	}

	var acc Account
	if err = json.Unmarshal(data, &acc); err != nil {
		log.Critical(err.Error())
		return nil, fmt.Errorf("unmarshal account error: %s", err.Error())
	}
	return &acc, nil
}
