package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/ldsec/lattigo/bfv"
	"github.com/stretchr/testify/assert"
)

func initWithCheck(t *testing.T) *shim.MockStub {
	// 创建stub
	var cc = new(TransferChainCode)
	stub := shim.NewMockStub("TransferChainCode", cc)

	// 调用Init
	res := stub.MockInit("fake_tx_id_1", [][]byte{[]byte("init")})
	assert.Equal(t, int32(shim.OK), res.Status, "Init failed: %s", string(res.Message))

	// 检查2个银行的所有账户是否都创建好了
	for _, bid := range banks {
		for _, acid := range accs {
			if _, err := cc.GetAccount(stub, bid, acid); err != nil {
				t.Error(err.Error())
			}
		}
	}

	return stub
}

func NewTestHeInfo() *HeInfo {
	// 同态加密的公私钥，私钥创建decryptor和encryptor
	kgen := bfv.NewKeyGenerator(defaultParams)
	sk, pk := kgen.GenKeyPair()
	// 使用公钥创建encryptor，如果是使用私钥创建，私钥每次对相同数据加密出来的数据是不一样的
	encryptor := bfv.NewEncryptorFromPk(defaultParams, pk)
	decryptor := bfv.NewDecryptor(defaultParams, sk)
	encoder := bfv.NewEncoder(defaultParams)

	return &HeInfo{sk, pk, encryptor, decryptor, encoder}
}

// 同态加密密钥信息转移到链下
type HeInfo struct {
	sk        *bfv.SecretKey // 私钥
	pk        *bfv.PublicKey // 公钥
	encryptor bfv.Encryptor  // 明文加密的密文
	decryptor bfv.Decryptor  // 解密密文
	encoder   bfv.Encoder    // 数据编码到明文
}

func TestTransferChainCode_Init(t *testing.T) {
	stub := initWithCheck(t)
	if stub == nil {
		t.Errorf("stub = %v", stub)
	}
}

func TestTransferChainCode_SetAccountBalance(t *testing.T) {
	testSetGetAccountBalance(t)
}

func testSetGetAccountBalance(t *testing.T) {
	he := NewTestHeInfo()
	stub := initWithCheck(t)
	balance := uint64(64)

	// 设置账户金额并进行余额正确性检查
	testSetAccountBalance(t, stub, BANK001, ACCOUNT001, he, balance)
	testCheckAccountBalance(t, stub, BANK001, ACCOUNT001, he, balance)
}

func testSetAccountBalance(t *testing.T, stub *shim.MockStub, bankID string, accountID string, he *HeInfo, bal uint64) {
	// 生成加密后的amount
	plain := bfv.NewPlaintext(defaultParams)
	he.encoder.EncodeUint([]uint64{bal}, plain)
	cipBal := he.encryptor.EncryptNew(plain)
	binBal, err := cipBal.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal cipher balance error: %s", err.Error())
	}

	// 创建链码并调用接口，合成链码调用参数
	args := [][]byte{[]byte("SetAccountBalance"),
		[]byte(bankID),
		[]byte(accountID),
		binBal}
	res := stub.MockInvokeWithSignedProposal("1", args, nil)
	assert.Equal(t, int32(shim.OK), res.Status, res.Message)
}

// 读取账户余额，检查是否正确
func testCheckAccountBalance(t *testing.T, stub *shim.MockStub, bankID string, accountID string, he *HeInfo, expBal uint64) {
	// 读取账户余额
	args := [][]byte{[]byte("QueryAccountBalance"),
		[]byte(bankID),
		[]byte(accountID)}
	res := stub.MockInvokeWithSignedProposal("1", args, nil)
	assert.Equal(t, int32(shim.OK), res.Status, res.Message)
	assert.NotNil(t, res.Payload)
	assert.NotEmpty(t, res.Payload)

	// 余额解密和正确性检查
	gotCipBal := &bfv.Ciphertext{}
	err := gotCipBal.UnmarshalBinary(res.Payload)
	if err != nil {
		t.Fatalf(err.Error())
	}
	gotPt := he.decryptor.DecryptNew(gotCipBal)
	gotBal := he.encoder.DecodeUint(gotPt)[0]
	if gotBal != expBal {
		t.Errorf("[%s - %s] balance not match, want = %v, got = %v", bankID, accountID, expBal, gotBal)
	} else {
		t.Logf("[%s - %s] balance is correct, balance = %v", bankID, accountID, gotBal)
	}
}

func testAddBankPK(t *testing.T, stub *shim.MockStub, bankID string, pk *bfv.PublicKey) {
	pkByte, err := pk.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	// 创建链码并调用接口，合成链码调用参数
	args := [][]byte{[]byte("AddBankPublicKey"),
		[]byte(bankID),
		pkByte}
	res := stub.MockInvokeWithSignedProposal("1", args, nil)
	assert.Equal(t, int32(shim.OK), res.Status, res.Message)
}

func TestTransferChainCode_AddBankPublicKey(t *testing.T) {
	bank1He := NewTestHeInfo()
	stub := initWithCheck(t)

	testAddBankPK(t, stub, BANK001, bank1He.pk)
	bank, err := GetBank(stub, BANK001)
	assert.NoError(t, err)
	assert.NotNil(t, bank)
	b, err := bank1He.pk.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, bank.pkByte, b, "pk byte not equal")
}

func TestTransferChainCode_Transfer(t *testing.T) {
	bank1He := NewTestHeInfo()
	bank2He := NewTestHeInfo()
	stub := initWithCheck(t)
	balance := uint64(100)

	// 给账户设置余额
	testSetAccountBalance(t, stub, BANK001, ACCOUNT001, bank1He, balance)
	testSetAccountBalance(t, stub, BANK002, ACCOUNT002, bank2He, 200)
	testCheckAccountBalance(t, stub, BANK001, ACCOUNT001, bank1He, balance)
	testCheckAccountBalance(t, stub, BANK002, ACCOUNT002, bank2He, 200)

	// 银行公钥上链
	testAddBankPK(t, stub, BANK001, bank1He.pk)
	testAddBankPK(t, stub, BANK002, bank2He.pk)

	// 转账30
	amount := "30"

	args := [][]byte{[]byte("Transfer"),
		[]byte(BANK001),
		[]byte(ACCOUNT001),
		[]byte(BANK002),
		[]byte(ACCOUNT002),
		[]byte(amount),
	}
	res := stub.MockInvokeWithSignedProposal("1", args, nil)
	assert.Equal(t, int32(shim.OK), res.Status, res.Message)
	// assert.Nil(t, res.Payload)

	// 检查余额是否正确
	testCheckAccountBalance(t, stub, BANK001, ACCOUNT001, bank1He, 70)
	testCheckAccountBalance(t, stub, BANK002, ACCOUNT002, bank2He, 230)
}
