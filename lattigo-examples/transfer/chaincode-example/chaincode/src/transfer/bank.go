package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
)

func GetBank(stub shim.ChaincodeStubInterface, bankID string) (*Bank, error) {
	pkByte, err := stub.GetState(bankID)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprint("GetBank read pkByte error with id", bankID))
	}

	bank := &Bank{pkByte: pkByte}
	return bank, nil
}

func PutBank(stub shim.ChaincodeStubInterface, bankID string, pkByte []byte) error {
	if err := stub.PutState(bankID, pkByte); err != nil {
		return errors.WithMessage(err, "PutBank")
	}
	return nil
}
