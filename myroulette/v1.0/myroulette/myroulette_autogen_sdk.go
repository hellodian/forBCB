package myroulette

import (
	"blockchain/smcsdk/sdk"
)

//SetSdk This is a method of MyRoulette
func (mr *MyRoulette) SetSdk(sdk sdk.ISmartContract) {
	mr.sdk = sdk
}

//GetSdk This is a method of MyRoulette
func (mr *MyRoulette) GetSdk() sdk.ISmartContract {
	return mr.sdk
}
