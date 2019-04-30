package myroulette

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdkimpl/object"
	"blockchain/smcsdk/sdkimpl/sdkhelper"
	"blockchain/smcsdk/utest"
	"fmt"
)

var (
	contractName       = "Roulette" //contract name
	contractMethods    = []string{"SetPublicKey(types.PubKey)", "SetSettings(string)", "SetRecFeeInfo(string)", "WithdrawFunds(string,types.Address,bn.Number)", "PlaceBet(string,int64,[]byte,[]byte,types.Address)", "SettleBet([]byte,int64)", "RefundBets([]byte,int64)"}
	contractInterfaces = []string{}
	orgID              = "orgNUjCm1i8RcoW2kVTbDw4vKW6jzfMxewJHjkhuiduhjuikjuyhnnjkuhujk111"
)

//TestObject This is a struct for test
type TestObject struct {
	obj *MyRoulette
}

//FuncRecover recover panic by Assert
func FuncRecover(err *types.Error) {
	if rerr := recover(); rerr != nil {
		if _, ok := rerr.(types.Error); ok {
			err.ErrorCode = rerr.(types.Error).ErrorCode
			err.ErrorDesc = rerr.(types.Error).ErrorDesc
			fmt.Println(err)
		} else {
			panic(rerr)
		}
	}
}

//NewTestObject This is a function
func NewTestObject(sender sdk.IAccount) *TestObject {
	return &TestObject{&MyRoulette{sdk: utest.UTP.ISmartContract}}
}

//transfer This is a method of TestObject
func (t *TestObject) transfer(balance bn.Number) *TestObject {
	contract := t.obj.sdk.Message().Contract()
	utest.Transfer(t.obj.sdk.Message().Sender(), t.obj.sdk.Helper().GenesisHelper().Token().Name(), contract.Account(), balance)
	t.obj.sdk = sdkhelper.OriginNewMessage(t.obj.sdk, contract, t.obj.sdk.Message().MethodID(), t.obj.sdk.Message().(*object.Message).OutputReceipts())
	return t
}

//setSender This is a method of TestObject
func (t *TestObject) setSender(sender sdk.IAccount) *TestObject {
	t.obj.sdk = utest.SetSender(sender.Address())
	return t
}

//run This is a method of TestObject
func (t *TestObject) run() *TestObject {
	t.obj.sdk = utest.ResetMsg()
	return t
}

//InitChain This is a method of TestObject
func (t *TestObject) InitChain() {
	utest.NextBlock(1)
	t.obj.InitChain()
	utest.Commit()
	return
}

//SetPublicKey This is a method of TestObject
func (t *TestObject) SetPublicKey(newPublicKey types.PubKey) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SetPublicKey(newPublicKey)
	utest.Commit()
	return
}

//SetSettings This is a method of TestObject
func (t *TestObject) SetSettings(newSettingsStr string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SetSettings(newSettingsStr)
	utest.Commit()
	return
}

//SetRecFeeInfo This is a method of TestObject
func (t *TestObject) SetRecFeeInfo(recFeeInfoStr string) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SetRecFeeInfo(recFeeInfoStr)
	utest.Commit()
	return
}

//WithdrawFunds This is a method of TestObject
func (t *TestObject) WithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.WithdrawFunds(tokenName, beneficiary, withdrawAmount)
	utest.Commit()
	return
}

//PlaceBet This is a method of TestObject
func (t *TestObject) PlaceBet(betInfoJson string, commitLastBlock int64, commit, signData []byte, refAddress types.Address) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.PlaceBet(betInfoJson, commitLastBlock, commit, signData, refAddress)
	utest.Commit()
	return
}

//SettleBet This is a method of TestObject
func (t *TestObject) SettleBet(reveal []byte, settleCount int64) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.SettleBet(reveal, settleCount)
	utest.Commit()
	return
}

//RefundBets This is a method of TestObject
func (t *TestObject) RefundBets(commit []byte, refundCount int64) (err types.Error) {
	err.ErrorCode = types.CodeOK
	defer FuncRecover(&err)
	utest.NextBlock(1)
	t.obj.RefundBets(commit, refundCount)
	utest.Commit()
	return
}
