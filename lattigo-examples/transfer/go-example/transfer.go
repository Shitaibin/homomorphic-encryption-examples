package main

import "log"

// 银行ID
const BANK001 = "bank_001"
const BANK002 = "bank_002"

// 账户ID
const ACCOUNT001 = "acc_001"
const ACCOUNT002 = "acc_002"

func main() {
	localTransfer()
	crossBankTransfer()
}

// localTransfer 行内转账
func localTransfer() {
	log.Println("开始行内转账测试")
	defer log.Println("行内转账测试完成")

	bank1 := NewBank(BANK001, nil)

	// 行内转账
	if err := bank1.LocalTransfer(ACCOUNT001, ACCOUNT002, 25); err != nil {
		log.Printf("local transfer error: %s", err.Error())

		// 失败金额检查
		if err := bank1.CheckAccountBalance(ACCOUNT001, 100); err != nil {
			log.Printf("check error: %s", err.Error())
		}
		if err := bank1.CheckAccountBalance(ACCOUNT002, 100); err != nil {
			log.Printf("check error: %s", err.Error())
		}

		return
	}

	if err := bank1.CheckAccountBalance(ACCOUNT001, 75); err != nil {
		log.Printf("check error: %s", err.Error())
	}
	if err := bank1.CheckAccountBalance(ACCOUNT002, 125); err != nil {
		log.Printf("check error: %s", err.Error())
	}
}

func crossBankTransfer() {
	log.Println("开始跨行转账测试")
	defer log.Println("行内跨行测试完成")

	bank1 := NewBank(BANK001, nil)
	bank2 := NewBank(BANK002, nil)

	var amount uint64 = 30
	if err := bank1.SubAccountBalance(ACCOUNT001, amount); err != nil {
		log.Printf("扣余额失败：[%s] [%s] [%v]", bank1.Name, ACCOUNT001, amount)
		return
	}
	if err := bank2.AddAccountBalance(ACCOUNT002, amount); err != nil {
		log.Printf("扣余额失败：[%s] [%s] [%v]", bank2.Name, ACCOUNT002, amount)
	}

	// 成功余额检查
	if err := bank1.CheckAccountBalance(ACCOUNT001, 70); err != nil {
		log.Printf("check error: %s", err.Error())
	}
	if err := bank1.CheckAccountBalance(ACCOUNT002, 130); err != nil {
		log.Printf("check error: %s", err.Error())
	}

	// TODO 增加原子性和错误处理
}
