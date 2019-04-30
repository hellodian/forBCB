package myroulette

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

//_setPublicKey This is a method of MyRoulette
func (mr *MyRoulette) _setPublicKey(v types.PubKey) {
	mr.sdk.Helper().StateHelper().McSet("/publicKey", &v)
}

//_publicKey This is a method of MyRoulette
func (mr *MyRoulette) _publicKey() types.PubKey {

	return *mr.sdk.Helper().StateHelper().McGetEx("/publicKey", new(types.PubKey)).(*types.PubKey)
}

//_clrPublicKey This is a method of MyRoulette
func (mr *MyRoulette) _clrPublicKey() {
	mr.sdk.Helper().StateHelper().McClear("/publicKey")
}

//_chkPublicKey This is a method of MyRoulette
func (mr *MyRoulette) _chkPublicKey() bool {
	return mr.sdk.Helper().StateHelper().Check("/publicKey")
}

//_McChkPublicKey This is a method of MyRoulette
func (mr *MyRoulette) _McChkPublicKey() bool {
	return mr.sdk.Helper().StateHelper().McCheck("/publicKey")
}

//_setLockedInBets This is a method of MyRoulette
func (mr *MyRoulette) _setLockedInBets(k string, v bn.Number) {
	mr.sdk.Helper().StateHelper().McSet(fmt.Sprintf("/lockedInBets/%v", k), &v)
}

//_lockedInBets This is a method of MyRoulette
func (mr *MyRoulette) _lockedInBets(k string) bn.Number {
	temp := bn.N(0)
	return *mr.sdk.Helper().StateHelper().McGetEx(fmt.Sprintf("/lockedInBets/%v", k), &temp).(*bn.Number)
}

//_clrLockedInBets This is a method of MyRoulette
func (mr *MyRoulette) _clrLockedInBets(k string) {
	mr.sdk.Helper().StateHelper().McClear(fmt.Sprintf("/lockedInBets/%v", k))
}

//_chkLockedInBets This is a method of MyRoulette
func (mr *MyRoulette) _chkLockedInBets(k string) bool {
	return mr.sdk.Helper().StateHelper().Check(fmt.Sprintf("/lockedInBets/%v", k))
}

//_McChkLockedInBets This is a method of MyRoulette
func (mr *MyRoulette) _McChkLockedInBets(k string) bool {
	return mr.sdk.Helper().StateHelper().McCheck(fmt.Sprintf("/lockedInBets/%v", k))
}

//_setSetting This is a method of MyRoulette
func (mr *MyRoulette) _setSetting(v *Setting) {
	mr.sdk.Helper().StateHelper().McSet("/setting", v)
}

//_setting This is a method of MyRoulette
func (mr *MyRoulette) _setting() *Setting {

	return mr.sdk.Helper().StateHelper().McGetEx("/setting", new(Setting)).(*Setting)
}

//_clrSetting This is a method of MyRoulette
func (mr *MyRoulette) _clrSetting() {
	mr.sdk.Helper().StateHelper().McClear("/setting")
}

//_chkSetting This is a method of MyRoulette
func (mr *MyRoulette) _chkSetting() bool {
	return mr.sdk.Helper().StateHelper().Check("/setting")
}

//_McChkSetting This is a method of MyRoulette
func (mr *MyRoulette) _McChkSetting() bool {
	return mr.sdk.Helper().StateHelper().McCheck("/setting")
}

//_setRecFeeInfo This is a method of MyRoulette
func (mr *MyRoulette) _setRecFeeInfo(v []RecFeeInfo) {
	mr.sdk.Helper().StateHelper().McSet("/recFeeInfo", &v)
}

//_recFeeInfo This is a method of MyRoulette
func (mr *MyRoulette) _recFeeInfo() []RecFeeInfo {

	return *mr.sdk.Helper().StateHelper().McGetEx("/recFeeInfo", new([]RecFeeInfo)).(*[]RecFeeInfo)
}

//_clrRecFeeInfo This is a method of MyRoulette
func (mr *MyRoulette) _clrRecFeeInfo() {
	mr.sdk.Helper().StateHelper().McClear("/recFeeInfo")
}

//_chkRecFeeInfo This is a method of MyRoulette
func (mr *MyRoulette) _chkRecFeeInfo() bool {
	return mr.sdk.Helper().StateHelper().Check("/recFeeInfo")
}

//_McChkRecFeeInfo This is a method of MyRoulette
func (mr *MyRoulette) _McChkRecFeeInfo() bool {
	return mr.sdk.Helper().StateHelper().McCheck("/recFeeInfo")
}

//_setRoundInfo This is a method of MyRoulette
func (mr *MyRoulette) _setRoundInfo(k string, v *RoundInfo) {
	mr.sdk.Helper().StateHelper().Set(fmt.Sprintf("/roundInfo/%v", k), v)
}

//_roundInfo This is a method of MyRoulette
func (mr *MyRoulette) _roundInfo(k string) *RoundInfo {

	return mr.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/roundInfo/%v", k), new(RoundInfo)).(*RoundInfo)
}

//_chkRoundInfo This is a method of MyRoulette
func (mr *MyRoulette) _chkRoundInfo(k string) bool {
	return mr.sdk.Helper().StateHelper().Check(fmt.Sprintf("/roundInfo/%v", k))
}

//_setBetInfo This is a method of MyRoulette
func (mr *MyRoulette) _setBetInfo(k1 string, k2 string, v *BetInfo) {
	mr.sdk.Helper().StateHelper().Set(fmt.Sprintf("/betInfo/%v/%v", k1, k2), v)
}

//_betInfo This is a method of MyRoulette
func (mr *MyRoulette) _betInfo(k1 string, k2 string) *BetInfo {

	return mr.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/betInfo/%v/%v", k1, k2), new(BetInfo)).(*BetInfo)
}

//_chkBetInfo This is a method of MyRoulette
func (mr *MyRoulette) _chkBetInfo(k1 string, k2 string) bool {
	return mr.sdk.Helper().StateHelper().Check(fmt.Sprintf("/betInfo/%v/%v", k1, k2))
}
