package main

import (
	"sync"

	"github.com/ldsec/lattigo/bfv"
)

// 银行ID
const BANK001 = "bank_001"
const BANK002 = "bank_002"

// 账户ID
const ACCOUNT001 = "acc_001"
const ACCOUNT002 = "acc_002"

type Bank struct {
	Name     string              // 银行名称
	Lock     sync.Mutex          // 账户的锁, TODO 并发时再加锁处理
	Accounts map[string]*Account // 银行内所有账户
}

// 创建银行，包含2个默认账号，但不用户不设置余额
// 由链下银行设置余额上链
func NewBank(name string, params *bfv.Parameters) *Bank {

	log.Infof("创建银行 [%v]", name)
	defer log.Infof("创建银行 [%v] 完成", name)

	if params == nil {
		params = defaultParams
	}

	b := &Bank{
		Name:     name,
		Lock:     sync.Mutex{},
		Accounts: make(map[string]*Account),
	}

	// 生成默认账号
	for _, ac := range NewDefaultAccounts() {
		log.Infof("[%s] add account [%s]", b.Name, ac.ID)
		b.Accounts[ac.ID] = ac
	}

	return b
}

// 完成行内2个账号的转账
func (b *Bank) LocalTransfer(acc1, acc2 string, amount *bfv.Ciphertext) error {
	log.Infof("[%s] [%s] 向 [%s] [%s] 转账", b.Name, acc1, b.Name, acc2)
	defer log.Info("转账完毕")

	// TODO amount 合法性检查，比如小于0；acc1余额不足等

	bal1 := b.Accounts[acc1].Balance
	bal2 := b.Accounts[acc2].Balance

	// 转账：同态加密运算
	newBal1 := evaluator.SubNew(bal1, amount)
	newBal2 := evaluator.AddNew(bal2, amount)

	// 转账设置回去
	b.Accounts[acc1].Balance = newBal1
	b.Accounts[acc2].Balance = newBal2

	return nil
}

// 检查用户余额是否匹配
// 匹配返回nil
// 进行同态解密的明文，然后解码得到原始数据
// TODO 转移到链下
// func (b *Bank) CheckAccountBalance(accID string, expectBalance uint64) error {
// 	log.Infof("[%s] check [%s] balance, epxect balance is [%d]", b.Name, accID, expectBalance)

// 	cepBal := b.Accounts[accID].Balance

// 	bal := b.He.encoder.DecodeUint(b.He.decryptor.DecryptNew(cepBal))[0]
// 	if bal != expectBalance {
// 		return fmt.Errorf("account [%s] [%s] balance is %d, expect balance is %d", b.Name, accID, bal, expectBalance)
// 	}

// 	log.Infof("[%s] [%s] balance is [%d], it's correct\n", b.Name, accID, bal)
// 	return nil
// }

// 设置用户余额
func (b *Bank) SetAccountBalance(accID string, balance *bfv.Ciphertext) error {
	log.Infof("[%s] set [%s] balance to [%d]", b.Name, accID, balance)

	b.Accounts[accID].Balance = balance
	return nil
}

func (b *Bank) AddAccountBalance(accID string, amount *bfv.Ciphertext) error {
	log.Infof("增加 [%s] [%s]", b.Name, accID)
	defer log.Info("余额修改完成")

	// TODO amount检查，非负

	bal := b.Accounts[accID].Balance

	// 转账：同态加密运算
	newBal := evaluator.AddNew(bal, amount)

	// 转账设置回去
	b.Accounts[accID].Balance = newBal
	return nil
}

func (b *Bank) SubAccountBalance(accID string, amount *bfv.Ciphertext) error {
	log.Infof("减少 [%s] [%s] ", b.Name, accID)
	defer log.Info("余额修改完成")

	// TODO amount检查，非负

	bal := b.Accounts[accID].Balance

	// TODO 检查余额是否不足以进行转账

	// 转账：同态加密运算
	newBal := evaluator.SubNew(bal, amount)

	// 转账设置回去
	b.Accounts[accID].Balance = newBal
	return nil
}

type Account struct {
	ID      string          // 银行账号
	Balance *bfv.Ciphertext // 同态加密过的余额
}

func NewDefaultAccounts() []*Account {
	a1 := &Account{ID: ACCOUNT001}
	a2 := &Account{ID: ACCOUNT002}
	return []*Account{a1, a2}
}

// TODO 同态加密密钥信息转移到链下
type HeInfo struct {
	sk        *bfv.SecretKey // 私钥
	pk        *bfv.PublicKey // 公钥
	encryptor bfv.Encryptor  // 明文加密的密文
	decryptor bfv.Decryptor  // 解密密文
	encoder   bfv.Encoder    // 数据编码到明文
}

func NewHeInfo(params *bfv.Parameters) *HeInfo {
	// 同态加密的公私钥，私钥创建decryptor和encryptor
	kgen := bfv.NewKeyGenerator(params)
	sk, pk := kgen.GenKeyPair()
	// 使用公钥创建encryptor，如果是使用私钥创建，私钥每次对相同数据加密出来的数据是不一样的
	encryptor := bfv.NewEncryptorFromPk(params, pk)
	decryptor := bfv.NewDecryptor(params, sk)
	encoder := bfv.NewEncoder(params)

	return &HeInfo{sk, pk, encryptor, decryptor, encoder}
}
