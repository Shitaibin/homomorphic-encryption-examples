package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/ldsec/lattigo/bfv"
	"github.com/pkg/errors"
)

func SetAccountBalance(BankID string, AccountID string, Balance uint64) (fab.TransactionID, peer.TxValidationCode, error) {
	if Balance < 0 {
		return "", 0, errors.New("Balance is negative")
	}

	// todo
	// 	  1. 请求参数为：BankID、AccountID、Balance（明文）
	//    2. Server利用该银行的公钥对Amount进行加密得到余额密文CipherBalance
	//    3. Server调用`SetAccountBalance`把CipherBalance上链
	//    4. 响应为：BankID、AccountID、Status（成功或失败）

	plain := bfv.NewPlaintext(defaultParams)
	encoder.EncodeUint([]uint64{Balance}, plain)
	cipBal := encryptor.EncryptNew(plain)
	binBal, err := cipBal.MarshalBinary()
	if err != nil {
		return "", 0, errors.WithMessage(err, "marshal cipher balance error")
	}

	args := packArgs([]string{BankID, AccountID, string(binBal)})
	req := channel.Request{
		ChaincodeID: ChainCodeName,
		Fcn:         "SetAccountBalance",
		Args:        args,
	}

	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := CLI.cc.Execute(req, reqPeers)
	if err != nil {
		return "", 0, errors.WithMessage(err, "invoke chaincode error")
	}

	logs.Info("Invoke chaincode response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)

	return resp.TransactionID, resp.TxValidationCode, nil
}

func GetAccountBalance(BankID string, AccountID string) (uint64, error) {
	// todo
	//    1. 请求参数为：BankID、AccountID
	//    2. Server调用链码`QueryAccountBalance`，获取用户余额，结果为同态加密的用户余额CipherBalance

	args := packArgs([]string{BankID, AccountID})
	req := channel.Request{
		ChaincodeID: ChainCodeName,
		Fcn:         "QueryAccountBalance",
		Args:        args,
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer0Org1)
	resp, err := CLI.cc.Query(req, reqPeers)
	if err != nil {
		return 0, errors.WithMessage(err, "query chaincode error")
	}

	//    3. Server利用该银行的私钥对CipherBalance进行解密，获得Balance
	gotCipBal := &bfv.Ciphertext{}
	err = gotCipBal.UnmarshalBinary(resp.Payload)
	if err != nil {
		return 0, errors.WithMessage(err, "unmarshal balance error")
	}
	gotPt := decryptor.DecryptNew(gotCipBal)
	gotBal := encoder.DecodeUint(gotPt)[0]

	//    4. 响应为：BankID、AccountID、Balance
	return gotBal, nil
}
