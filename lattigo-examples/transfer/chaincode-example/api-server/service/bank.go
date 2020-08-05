package service

import (
	"api-server/models"
	"bufio"
	"encoding/base64"
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/ldsec/lattigo/bfv"
	"github.com/pkg/errors"
)

// 本同态加密样例采用的默认参数
var defaultParams = bfv.DefaultParams[bfv.PN13QP218]

// 同态加密执行器
// var evaluator bfv.Evaluator

// 用于处理同态加密明文的序列化和反序列化
var encoder bfv.Encoder

// 银行的私钥和公钥, 加密解密器
type He struct {
	BankID    string
	SK        *bfv.SecretKey
	PK        *bfv.PublicKey
	encryptor bfv.Encryptor
	decryptor bfv.Decryptor
}

// 存储每家银行的He
// 数据并不进行持久化，启动后需要调用创建银行接口生成
var BanksHe map[string]*He

func init() {
	BanksHe = make(map[string]*He)

	defaultParams.T = 0x3ee0001

	// evaluator = bfv.NewEvaluator(defaultParams)
	encoder = bfv.NewEncoder(defaultParams)
}

// 获取银行密钥信息
func GetBank(bid string) *He {
	return BanksHe[bid]
}

// SkPkToString
func SkPkToString(SK *bfv.SecretKey, PK *bfv.PublicKey) (string, string, error) {
	if SK == nil {
		return "", "", NoBankError
	}

	skb, err := SK.MarshalBinary()
	if err != nil {
		return "", "", err
	}
	pkb, err := PK.MarshalBinary()
	if err != nil {
		return "", "", err
	}
	return string(skb), string(pkb), nil
}

func NewBank(bid string) (*models.Bank, error) {
	var (
		SK        *bfv.SecretKey
		PK        *bfv.PublicKey
		encryptor bfv.Encryptor
		decryptor bfv.Decryptor
	)

	kgen := bfv.NewKeyGenerator(defaultParams)
	SK, PK = kgen.GenKeyPair()
	encryptor = bfv.NewEncryptorFromPk(defaultParams, PK)
	decryptor = bfv.NewDecryptor(defaultParams, SK)

	logs.Info("Generate bank %s keys success.", bid)

	// 转换成结构体
	sk, pk, err := SkPkToString(SK, PK)
	if err != nil {
		return nil, errors.WithMessage(err, "New Bank")
	}
	logs.Info("Bank = %v, len(SK) = %d, len(PK) = %d", bid, len(sk), len(pk))

	args := packArgs([]string{bid, pk})
	req := channel.Request{
		ChaincodeID: ChainCodeName,
		Fcn:         "AddBankPublicKey",
		Args:        args,
	}

	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := CLI.cc.Execute(req, reqPeers)
	if err != nil {
		return nil, errors.WithMessage(err, "invoke chaincode error")
	}

	logs.Info("Invoke chaincode response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)

	// 银行密钥信息保存到map
	he := &He{
		BankID:    bid,
		SK:        SK,
		PK:        PK,
		encryptor: encryptor,
		decryptor: decryptor,
	}
	BanksHe[bid] = he

	return &models.Bank{
		BankID:    bid,
		TxID:      resp.TransactionID,
		ValidCode: resp.TxValidationCode,
		Message:   "success",
	}, nil
}

func SaveBankKeys(bid, fp string) error {
	logs.Info("Save bank %s keys to %s", bid, fp)

	bank := GetBank(bid)
	sk, pk, err := SkPkToString(bank.SK, bank.PK)
	if err != nil {
		return errors.WithMessage(err, "Get key failed")
	}
	base64sk := base64.StdEncoding.EncodeToString([]byte(sk))
	base64pk := base64.StdEncoding.EncodeToString([]byte(pk))

	// 打开文件
	f, err := os.Create(fp)
	// f, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE, 0666) //打开文件
	if err != nil {
		return errors.WithMessage(err, "open file error")
	}
	defer f.Close()

	// 准备待写文件
	skTitle := "----------------Secret Key----------------\n"
	pkTitle := "----------------Public Key----------------\n"
	lines := []string{skTitle, base64sk, "\n", pkTitle, base64pk}

	// 写文件
	logs.Info("Ready to writer %s", fp)
	wf := bufio.NewWriter(f)
	defer wf.Flush()
	for _, l := range lines {
		_, err = wf.WriteString(l)
		if err != nil {
			return errors.WithMessage(err, "write error")
		}
	}

	logs.Info("Write keys file done.")

	return nil
}
