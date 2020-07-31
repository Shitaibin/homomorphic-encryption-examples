package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/ldsec/lattigo/bfv"
	"github.com/stretchr/testify/assert"
)

func initWithCheck(t *testing.T) (*TransferChainCode, *shim.MockStub) {
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

	return cc, stub
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

func TestTransferChainCode_Init(t *testing.T) {
	cc, stub := initWithCheck(t)
	if cc == nil || stub == nil {
		t.Errorf("cc = %v, stub = %v", cc, stub)
	}
}

func TestTransferChainCode_SetAccountBalance(t *testing.T) {
	// 生成加密后的amount
	he := NewTestHeInfo()
	plain := bfv.NewPlaintext(defaultParams)
	bal := uint64(100)
	he.encoder.EncodeUint([]uint64{bal}, plain)
	cipBal := he.encryptor.EncryptNew(plain)
	binBal, err := cipBal.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal cipher balance error: %s", err.Error())
	}

	// 创建链码并调用接口，合成链码调用参数
	args := [][]byte{[]byte("SetAccountBalance"),
		[]byte(BANK001),
		[]byte(ACCOUNT001),
		binBal}
	_, stub := initWithCheck(t)
	res := stub.MockInvokeWithSignedProposal("1", args, nil)
	assert.Equal(t, int32(shim.OK), res.Status, res.Message)

	// 读取账户余额，检查是否正确
	args = [][]byte{[]byte("QueryAccountBalance"),
		[]byte(BANK001),
		[]byte(ACCOUNT001)}
	res = stub.MockInvokeWithSignedProposal("1", args, nil)
	assert.Equal(t, int32(shim.OK), res.Status, res.Message)
	assert.NotNil(t, res.Payload)
	assert.NotEmpty(t, res.Payload)

	gotCipBal := &bfv.Ciphertext{}
	err = gotCipBal.UnmarshalBinary(res.Payload)
	if err != nil {
		t.Fatalf(err.Error())
	}
	gotPt := he.decryptor.DecryptNew(gotCipBal)
	gotBal := he.encoder.DecodeUint(gotPt)[0]
	if gotBal != bal {
		t.Errorf("balance not match, want = %v, got = %v", bal, gotBal)
	}
}
