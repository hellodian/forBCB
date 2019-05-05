package excellencies

import (
	"fmt"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdk/bn"
)

//_setSecretSigner This is a method of Excellencies
func (e *Excellencies) _setSecretSigner(v types.PubKey) {
	e.sdk.Helper().StateHelper().McSet("/secretSigner", &v)
}

//_secretSigner This is a method of Excellencies
func (e *Excellencies) _secretSigner() types.PubKey {

	return *e.sdk.Helper().StateHelper().McGetEx("/secretSigner", new(types.PubKey)).(*types.PubKey)
}

//_clrSecretSigner This is a method of Excellencies
func (e *Excellencies) _clrSecretSigner() {
	e.sdk.Helper().StateHelper().McClear("/secretSigner")
}

//_chkSecretSigner This is a method of Excellencies
func (e *Excellencies) _chkSecretSigner() bool {
	return e.sdk.Helper().StateHelper().Check("/secretSigner")
}

//_McChkSecretSigner This is a method of Excellencies
func (e *Excellencies) _McChkSecretSigner() bool {
	return e.sdk.Helper().StateHelper().McCheck("/secretSigner")
}

//_setLockedInBets This is a method of Excellencies
func (e *Excellencies) _setLockedInBets(k string, v bn.Number) {
	e.sdk.Helper().StateHelper().McSet(fmt.Sprintf("/lockedInBets/%v", k), &v)
}

//_lockedInBets This is a method of Excellencies
func (e *Excellencies) _lockedInBets(k string) bn.Number {
	temp := bn.N(0)
	return *e.sdk.Helper().StateHelper().McGetEx(fmt.Sprintf("/lockedInBets/%v", k), &temp).(*bn.Number)
}

//_clrLockedInBets This is a method of Excellencies
func (e *Excellencies) _clrLockedInBets(k string) {
	e.sdk.Helper().StateHelper().McClear(fmt.Sprintf("/lockedInBets/%v", k))
}

//_chkLockedInBets This is a method of Excellencies
func (e *Excellencies) _chkLockedInBets(k string) bool {
	return e.sdk.Helper().StateHelper().Check(fmt.Sprintf("/lockedInBets/%v", k))
}

//_McChkLockedInBets This is a method of Excellencies
func (e *Excellencies) _McChkLockedInBets(k string) bool {
	return e.sdk.Helper().StateHelper().McCheck(fmt.Sprintf("/lockedInBets/%v", k))
}

//_setSettings This is a method of Excellencies
func (e *Excellencies) _setSettings(v *Settings) {
	e.sdk.Helper().StateHelper().McSet("/settings", v)
}

//_settings This is a method of Excellencies
func (e *Excellencies) _settings() *Settings {

	return e.sdk.Helper().StateHelper().McGetEx("/settings", new(Settings)).(*Settings)
}

//_clrSettings This is a method of Excellencies
func (e *Excellencies) _clrSettings() {
	e.sdk.Helper().StateHelper().McClear("/settings")
}

//_chkSettings This is a method of Excellencies
func (e *Excellencies) _chkSettings() bool {
	return e.sdk.Helper().StateHelper().Check("/settings")
}

//_McChkSettings This is a method of Excellencies
func (e *Excellencies) _McChkSettings() bool {
	return e.sdk.Helper().StateHelper().McCheck("/settings")
}

//_setRecFeeInfo This is a method of Excellencies
func (e *Excellencies) _setRecFeeInfo(v []RecFeeInfo) {
	e.sdk.Helper().StateHelper().McSet("/recFeeInfo", &v)
}

//_recFeeInfo This is a method of Excellencies
func (e *Excellencies) _recFeeInfo() []RecFeeInfo {

	return *e.sdk.Helper().StateHelper().McGetEx("/recFeeInfo", new([]RecFeeInfo)).(*[]RecFeeInfo)
}

//_clrRecFeeInfo This is a method of Excellencies
func (e *Excellencies) _clrRecFeeInfo() {
	e.sdk.Helper().StateHelper().McClear("/recFeeInfo")
}

//_chkRecFeeInfo This is a method of Excellencies
func (e *Excellencies) _chkRecFeeInfo() bool {
	return e.sdk.Helper().StateHelper().Check("/recFeeInfo")
}

//_McChkRecFeeInfo This is a method of Excellencies
func (e *Excellencies) _McChkRecFeeInfo() bool {
	return e.sdk.Helper().StateHelper().McCheck("/recFeeInfo")
}

//_setRoundInfo This is a method of Excellencies
func (e *Excellencies) _setRoundInfo(k string, v *RoundInfo) {
	e.sdk.Helper().StateHelper().Set(fmt.Sprintf("/roundInfo/%v", k), v)
}

//_roundInfo This is a method of Excellencies
func (e *Excellencies) _roundInfo(k string) *RoundInfo {

	return e.sdk.Helper().StateHelper().GetEx(fmt.Sprintf("/roundInfo/%v", k), new(RoundInfo)).(*RoundInfo)
}

//_chkRoundInfo This is a method of Excellencies
func (e *Excellencies) _chkRoundInfo(k string) bool {
	return e.sdk.Helper().StateHelper().Check(fmt.Sprintf("/roundInfo/%v", k))
}

//_setPoolAmount This is a method of Excellencies
func (e *Excellencies) _setPoolAmount(k string, v bn.Number) {
	e.sdk.Helper().StateHelper().McSet(fmt.Sprintf("/poolAmount/%v", k), &v)
}

//_poolAmount This is a method of Excellencies
func (e *Excellencies) _poolAmount(k string) bn.Number {
	temp := bn.N(0)
	return *e.sdk.Helper().StateHelper().McGetEx(fmt.Sprintf("/poolAmount/%v", k), &temp).(*bn.Number)
}

//_clrPoolAmount This is a method of Excellencies
func (e *Excellencies) _clrPoolAmount(k string) {
	e.sdk.Helper().StateHelper().McClear(fmt.Sprintf("/poolAmount/%v", k))
}

//_chkPoolAmount This is a method of Excellencies
func (e *Excellencies) _chkPoolAmount(k string) bool {
	return e.sdk.Helper().StateHelper().Check(fmt.Sprintf("/poolAmount/%v", k))
}

//_McChkPoolAmount This is a method of Excellencies
func (e *Excellencies) _McChkPoolAmount(k string) bool {
	return e.sdk.Helper().StateHelper().McCheck(fmt.Sprintf("/poolAmount/%v", k))
}

//_setSlipper This is a method of Excellencies
func (e *Excellencies) _setSlipper(k string, v SlipperInfo) {
	e.sdk.Helper().StateHelper().McSet(fmt.Sprintf("/slipper/%v", k), &v)
}

//_slipper This is a method of Excellencies
func (e *Excellencies) _slipper(k string) SlipperInfo {

	return *e.sdk.Helper().StateHelper().McGetEx(fmt.Sprintf("/slipper/%v", k), new(SlipperInfo)).(*SlipperInfo)
}

//_clrSlipper This is a method of Excellencies
func (e *Excellencies) _clrSlipper(k string) {
	e.sdk.Helper().StateHelper().McClear(fmt.Sprintf("/slipper/%v", k))
}

//_chkSlipper This is a method of Excellencies
func (e *Excellencies) _chkSlipper(k string) bool {
	return e.sdk.Helper().StateHelper().Check(fmt.Sprintf("/slipper/%v", k))
}

//_McChkSlipper This is a method of Excellencies
func (e *Excellencies) _McChkSlipper(k string) bool {
	return e.sdk.Helper().StateHelper().McCheck(fmt.Sprintf("/slipper/%v", k))
}

//_setPoker This is a method of Excellencies
func (e *Excellencies) _setPoker(v []*Poker) {
	e.sdk.Helper().StateHelper().Set("/poker", v)
}

//_poker This is a method of Excellencies
func (e *Excellencies) _poker() []*Poker {

	return *e.sdk.Helper().StateHelper().GetEx("/poker", new([]*Poker)).(*[]*Poker)
}

//_chkPoker This is a method of Excellencies
func (e *Excellencies) _chkPoker() bool {
	return e.sdk.Helper().StateHelper().Check("/poker")
}

