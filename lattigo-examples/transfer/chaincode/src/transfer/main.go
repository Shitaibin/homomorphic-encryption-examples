package main

import (
	"encoding/json"
	"fmt"

	"github.com/ldsec/lattigo/bfv"

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
		return t.SetAccountBalance(stub, args)
	case "Transfer":
		return t.Transfer(stub, args)
	case "QueryAccountBalance":
		return t.QueryAccountBalance(stub, args)
	}
	return shim.Success([]byte("ok"))
}

// SetAccountBalance 设置用户余额
// args[0]: bankID, string
// args[1]: accountID, string
// args[2]: cipherBalance, *bfv.Ciphertext
func (t *TransferChainCode) SetAccountBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 3 {
		return shim.Error(fmt.Sprintf("need 3 args, got %v", len(args)))
	}

	bankID, accountID := args[0], args[1]
	log.Infof("set balance [%s] [%s]", bankID, accountID)

	var err error

	// 解析参数
	cipBal := &bfv.Ciphertext{}
	err = cipBal.UnmarshalBinary([]byte(args[2]))
	if err != nil {
		msg := fmt.Sprintf("unmarshal cipher balance error, bank = %v, account = %v, error = %v",
			bankID, accountID, err.Error())
		return shim.Error(msg)
	}

	// 需要读取account余额
	acc, err := t.GetAccount(stub, bankID, accountID)
	if err != nil {
		msg := fmt.Sprintf("get account error, bank = %v, account = %v, error = %v",
			bankID, accountID, err.Error())
		return shim.Error(msg)
	}

	// 设置account新余额
	data, err := cipBal.MarshalBinary()
	if err != nil {
		shim.Error(err.Error())
	}
	acc.Balance = data

	// 保存account
	if err = t.PutAccount(stub, acc); err != nil {
		msg := fmt.Sprintf("save account error, bank = %v, account = %v, error = %v",
			bankID, accountID, err.Error())
		return shim.Error(msg)
	}

	return shim.Success([]byte("ok"))
}

// todo 完成2个账户间的转账
func (t *TransferChainCode) Transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success([]byte("ok"))
}

// QueryAccountBalance 查询某个账户的余额
// args[0]: bankID, string
// args[1]: accountID, string
func (t *TransferChainCode) QueryAccountBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(fmt.Sprintf("need 2 args, got %v", len(args)))
	}

	bankID, accountID := args[0], args[1]
	log.Infof("QueryAccountBalance [%s] [%s]", bankID, accountID)

	var err error
	acc, err := t.GetAccount(stub, bankID, accountID)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(acc.Balance)
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
