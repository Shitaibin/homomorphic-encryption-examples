package main

import (
	"github.com/ldsec/lattigo/bfv"
	"github.com/pkg/errors"
)

// 银行ID
const BANK001 = "bank_001"
const BANK002 = "bank_002"

// 账户ID
const ACCOUNT001 = "acc_001"
const ACCOUNT002 = "acc_002"

var banks = []string{BANK001, BANK002}
var accs = []string{ACCOUNT001, ACCOUNT002}

type Account struct {
	ID      string // 银行账号
	BankID  string // 银行ID
	Balance []byte // 同态加密过的余额，并且进行了序列化
	// Balance *bfv.Ciphertext // 同态加密过的余额
}

// 创建2个银行，分别包含2个默认账号，但不用户不设置余额
// 由链下银行设置余额上链
func NewBankAccounts() []*Account {

	var accounts []*Account
	for _, b := range banks {
		for _, id := range accs {
			accounts = append(accounts, &Account{ID: id, BankID: b})
		}
	}
	return accounts
}

// Bank 银行信息
type Bank struct {
	pkByte []byte // 公钥用于生成encryptor加密明文

	// 增加缓存Bank，减少获取pk、encryptor的计算，demo中可不做性能优化
	// pk        *bfv.PublicKey // 反序列化而得
	// encryptor bfv.Encryptor  // 由公钥生成
}

func (b *Bank) EncryptAmountNew(amount uint64) (*bfv.Ciphertext, error) {
	// 生成公钥和encryptor
	pk := &bfv.PublicKey{}
	if err := pk.UnmarshalBinary(b.pkByte); err != nil {
		return nil, errors.WithMessage(err, "EncryptNew")
	}
	encryptor := bfv.NewEncryptorFromPk(defaultParams, pk)

	// 金额编码 & 加密
	plain := bfv.NewPlaintext(defaultParams)
	encoder.EncodeUint([]uint64{amount}, plain)
	return encryptor.EncryptNew(plain), nil
}
