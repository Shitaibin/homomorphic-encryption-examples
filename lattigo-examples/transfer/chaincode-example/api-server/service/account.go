package service

import "github.com/pkg/errors"

func SetAccountBalance(BankID string, AccountID string, Balance float64) error {
	if Balance < 0 {
		return errors.New("Balance is negative")
	}

	// todo
	// 	  1. 请求参数为：BankID、AccountID、Balance（明文）
	//    2. Server利用该银行的公钥对Amount进行加密得到余额密文CipherBalance
	//    3. Server调用`SetAccountBalance`把CipherBalance上链
	//    4. 响应为：BankID、AccountID、Status（成功或失败）
	return nil
}

func GetAccountBalance(BankID string, AccountID string) (float64, error) {
	// todo
	return 0, errors.New("mock")
}
