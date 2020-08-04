package service

import (
	"strconv"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"

	"github.com/astaxie/beego/logs"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/pkg/errors"
)

func Transfer(FromBankID, FromAccountID, ToBankID, ToAccountID string, Amount float64) (fab.TransactionID, peer.TxValidationCode, error) {
	// todo
	//    1. 请求参数为：FromBankID、FromAccountID、ToBankID、ToAccountID、Amount
	//    2. Server调用链码`Transfer`进行链上用户余额转账
	//    3. 响应为：FromBankID、FromAccountID、ToBankID、ToAccountID、Amount、Status（成功或失败）

	amountStr := strconv.FormatFloat(Amount, 'f', -1, 64)
	args := packArgs([]string{FromBankID, FromAccountID, ToBankID, ToAccountID, amountStr})
	req := channel.Request{
		ChaincodeID: ChainCodeName,
		Fcn:         "Transfer",
		Args:        args,
	}

	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := CLI.cc.Execute(req, reqPeers)
	if err != nil {
		return "", 0, errors.WithMessage(err, "Transfer invoke error")
	}

	logs.Info("Invoke chaincode response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)

	return resp.TransactionID, resp.TxValidationCode, nil
}
