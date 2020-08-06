package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/ldsec/lattigo/bfv"
	"github.com/pkg/errors"
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

	log.Info("TransferChainCode init done.")

	return shim.Success([]byte("ok"))
}

func (t *TransferChainCode) Invoke(stub shim.ChaincodeStubInterface) (resp pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			resp = shim.Error(fmt.Sprintf("catch panic: %v", err.(error).Error()))
		}
	}()

	f, args := stub.GetFunctionAndParameters()

	log.Infof("Invoke func = %v, args = %v", f, len(args))

	switch f {
	case "AddBankPublicKey":
		return t.AddBankPublicKey(stub, args)
	case "SetAccountBalance":
		return t.SetAccountBalance(stub, args)
	case "QueryAccountBalance":
		return t.QueryAccountBalance(stub, args)
	case "Transfer":
		return t.Transfer(stub, args)
	}
	return shim.Success([]byte("ok"))
}

// SetAccountBalance 设置用户余额
// args[0]: bankID, string
// args[1]: accountID, string
// args[2]: cipherBalance, *bfv.Ciphertext
func (t *TransferChainCode) SetAccountBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 3 {
		return shim.Error(fmt.Sprintf("SetAccountBalance need 3 args, got %v", len(args)))
	}

	bankID, accountID := args[0], args[1]
	log.Infof("SetAccountBalance [%s] [%s]", bankID, accountID)

	var err error

	var msg string = "transfer success"
	defer func() {
		log.Debug(msg)
	}()

	// 解析参数
	cipBal := &bfv.Ciphertext{}
	err = cipBal.UnmarshalBinary([]byte(args[2]))
	if err != nil {
		msg = fmt.Sprintf("unmarshal cipher balance error, bank = %v, account = %v, error = %v",
			bankID, accountID, err.Error())
		return shim.Error(msg)
	}

	// 需要读取account余额
	acc, err := t.GetAccount(stub, bankID, accountID)
	if err != nil {
		msg = fmt.Sprintf("get account error, bank = %v, account = %v, error = %v",
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
		msg = fmt.Sprintf("save account error, bank = %v, account = %v, error = %v",
			bankID, accountID, err.Error())
		return shim.Error(msg)
	}

	return shim.Success([]byte("ok"))
}

// AddBankPublicKey 银行上传同态机密公钥
// args[0]: bankID, string
// args[0]: pkByte, []byte, 公钥序列化的字节码
func (t *TransferChainCode) AddBankPublicKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(fmt.Sprintf("AddBankPublicKey need 2 args, got %v", len(args)))
	}

	bankID, pkByte := args[0], []byte(args[1])
	// 不查询银行，直接Add
	if err := PutBank(stub, bankID, pkByte); err != nil {
		return shim.Error(errors.WithMessage(err, "AddBankPublicKey").Error())
	}
	return shim.Success([]byte(fmt.Sprintf("Add public key for bank [%v] success.", bankID)))
}

// Transfer 完成2个账户间的转账
// args[0]: fromBankID, string
// args[1]: fromAccountID, string
// args[2]: toBankID, string
// args[3]: toAccountID, string
// args[4]: amount, uint64, SHOULD be plain text
func (t *TransferChainCode) Transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 5 {
		return shim.Error(fmt.Sprintf("AddBankPublicKey need 5 args, got %v", len(args)))
	}

	var msg string = "transfer success"
	defer func() {
		log.Info(msg)
	}()

	fromBankID, fromAccountID, toBankID, toAccountID := args[0], args[1], args[2], args[3]
	// 获取2个用户的余额，
	from, err := t.GetAccount(stub, fromBankID, fromAccountID)
	if err != nil {
		msg = fmt.Sprintf("read from account error = %v", err.Error())
		return shim.Error(msg)
	}

	to, err := t.GetAccount(stub, toBankID, toAccountID)
	if err != nil {
		msg = fmt.Sprintf("read to account error = %v", err.Error())
		return shim.Error(msg)
	}

	// 解析amount
	amount, err := strconv.ParseUint(args[4], 10, 64)
	if err != nil {
		msg = fmt.Sprintf("parse amount error = %v", err.Error())
		return shim.Error(msg)
	}

	// 转账参数合法性检查
	if from == nil {
		return shim.Error("from is nil")
	}
	if from.Balance == nil || len(from.Balance) <= 0 {
		return shim.Error("from.Balance is nil or empty")
	}

	if to == nil {
		return shim.Error("to is nil")
	}
	if to.Balance == nil || len(to.Balance) <= 0 {
		return shim.Error("to.Balance is nil or empty")
	}

	// 获取2个账户余额
	fromBal, err := unmarshalBal(from.Balance)
	if err != nil {
		msg = fmt.Sprintf("unmarshal from account [%s - %s] balance error = %v",
			from.BankID, from.ID, err.Error())
		return shim.Error(msg)
	}
	toBal, err := unmarshalBal(to.Balance)
	if err != nil {
		msg := fmt.Sprintf("unmarshal to account [%s - %s] balance error = %v",
			to.BankID, to.ID, err.Error())
		return shim.Error(msg)
	}

	// 每家银行的余额需要分别使用自己的密钥计算
	fromBank, err := GetBank(stub, fromBankID)
	if err != nil {
		msg = errors.WithMessage(err, "Transfer error").Error()
		return shim.Error(msg)
	}
	toBank, err := GetBank(stub, toBankID)
	if err != nil {
		msg = errors.WithMessage(err, "Transfer error").Error()
		return shim.Error(msg)
	}

	fromAmount, err := fromBank.EncryptAmountNew(amount)
	if err != nil {
		msg = errors.WithMessage(err, "Transfer encrypt from amount error").Error()
		return shim.Error(msg)
	}
	toAmount, err := toBank.EncryptAmountNew(amount)
	if err != nil {
		msg = errors.WithMessage(err, "Transfer encrypt from amount error").Error()
		return shim.Error(msg)
	}

	// 调用2个银行密钥对amount加密，然后计算
	evaluator.Sub(fromBal, fromAmount, fromBal)
	evaluator.Add(toBal, toAmount, toBal)
	fromData, err := fromBal.MarshalBinary()
	if err != nil {
		msg = fmt.Sprintf("marshal from balance error = %v", err.Error())
		return shim.Error(msg)
	}
	toData, err := toBal.MarshalBinary()
	if err != nil {
		msg = fmt.Sprintf("marshal to balance error = %v", err.Error())
		return shim.Error(msg)
	}

	// 把新余额保存到用户
	from.Balance = fromData
	to.Balance = toData

	if err := t.PutAccount(stub, from); err != nil {
		msg = fmt.Sprintf("save from account error, bank = %v, account = %v, error = %v",
			from.BankID, from.ID, err.Error())
		return shim.Error(msg)
	}
	if err := t.PutAccount(stub, to); err != nil {
		msg = fmt.Sprintf("save to account error, bank = %v, account = %v, error = %v",
			to.BankID, to.ID, err.Error())
		return shim.Error(msg)
	}

	// 发布转账完成的事件，包含转账的2个用户
	err = stub.SetEvent(ChainCodeEventName_Transfer, NewMarshaledTransferEvent(fromBankID, fromAccountID, toBankID, toAccountID, args[4]))
	if err != nil {
		msg = errors.WithMessage(err, "set event error").Error()
		return shim.Error(msg)
	}

	return shim.Success([]byte("ok"))
}

func unmarshalBal(data []byte) (*bfv.Ciphertext, error) {
	gotCipBal := &bfv.Ciphertext{}
	err := gotCipBal.UnmarshalBinary(data)
	return gotCipBal, err
}

// QueryAccountBalance 查询某个账户的余额
// args[0]: bankID, string
// args[1]: accountID, string
func (t *TransferChainCode) QueryAccountBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error(fmt.Sprintf("QueryAccountBalance need 2 args, got %v", len(args)))
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
