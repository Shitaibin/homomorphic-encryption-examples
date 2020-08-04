package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/ldsec/lattigo/bfv"
	"github.com/pkg/errors"
)

// 本同态加密样例采用的默认参数
var defaultParams = bfv.DefaultParams[bfv.PN13QP218]

// 同态加密执行器
// var evaluator bfv.Evaluator

// 用于处理同态加密明文的序列化和反序列化
var encoder bfv.Encoder

func init() {
	defaultParams.T = 0x3ee0001

	// evaluator = bfv.NewEvaluator(defaultParams)
	encoder = bfv.NewEncoder(defaultParams)
}

// 银行的私钥和公钥, 加密解密器
var (
	SK        *bfv.SecretKey
	PK        *bfv.PublicKey
	encryptor bfv.Encryptor
	decryptor bfv.Decryptor
)

// SkPkToString
func SkPkToString() (string, string, error) {
	skb, err := SK.MarshalBinary()
	if err != nil {
		return "", "", err
	}
	pkb, err := SK.MarshalBinary()
	if err != nil {
		return "", "", err
	}
	return string(skb), string(pkb), nil
}

type Bank struct {
	BankID string
	SK     string
	PK     string
}

func NewBank(bid string) (*Bank, error) {
	kgen := bfv.NewKeyGenerator(defaultParams)
	SK, PK = kgen.GenKeyPair()
	encryptor = bfv.NewEncryptorFromPk(defaultParams, PK)
	decryptor = bfv.NewDecryptor(defaultParams, SK)

	// 转换成结构体
	sk, pk, err := SkPkToString()
	if err != nil {
		return nil, errors.WithMessage(err, "New Bank")
	}
	logs.Info("Bank = %v, len(SK) = %d, len(PK) = %d", bid, len(sk), len(pk))

	return &Bank{
		BankID: bid,
		SK:     sk,
		PK:     pk,
	}, nil
}
