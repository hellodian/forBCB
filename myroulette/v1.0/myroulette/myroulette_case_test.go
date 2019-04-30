package myroulette

import (
	"blockchain/algorithm"
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/utest"
	"common/keys"
	"common/kms"
	"encoding/hex"
	"fmt"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tmlibs/common"
	"gopkg.in/check.v1"
	"io/ioutil"
	"math"
	"testing"
)

//Test This is a function
func Test(t *testing.T) { check.TestingT(t) }

const (
	ownerName = "local_owner"
	password  = "12345678"
)

var (
	cdc = amino.NewCodec()
)

func init() {
	crypto.RegisterAmino(cdc)
	crypto.SetChainId("local")
	kms.InitKMS("./.keystore", "local_mode", "", "", "0x1003")
	kms.GenPrivKey(ownerName, []byte(password))
}

//MySuite This is a struct
type MySuite struct{}

var _= check.Suite(&MySuite{})

//TestMyRoulette_SetPublicKey This is a method of MySuite
func (mysuit *MySuite) TestMyRoulette_SetPublicKey(c *check.C) () {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	pubKey, _ := kms.GetPubKey(ownerName, []byte(password))

	account := utest.NewAccount(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1000000000))

	var tests = []struct {
		account sdk.IAccount
		pubKey  []byte
		desc    string
		code    uint32
	}{
		{contractOwner, pubKey, "--正常流程--", types.CodeOK},
		{contractOwner, []byte("0xff"), "--异常流程--公钥长度不正确--", types.ErrInvalidParameter},
		{account, pubKey, "--异常流程--非owner调用--", types.ErrNoAuthorization},
	}

	for _, item := range tests {
		utest.AssertError(test.run().setSender(item.account).SetPublicKey(item.pubKey), item.code)
	}
}

//TestMyRoulette_SetSettings This is a method of MySuite
func (mysuit *MySuite) TestMyRoulette_SetSettings(c *check.C) () {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	setting := Setting{}
	setting.MaxProfit = 2E12
	setting.MaxLimit = 2E10
	setting.MinLimit = 1E8
	setting.FeeRatio = 50
	setting.FeeMiniNum = 300000
	setting.SendToCltRatio = 100
	setting.BetExpirationBlocks = 250
	setting.TokenNames = []string{test.obj.sdk.Helper().GenesisHelper().Token().Name()}
	resBytes1, _ := jsoniter.Marshal(setting)

	setting.MaxLimit = 2E9
	setting.MinLimit = 2E10
	resBytes2, _ := jsoniter.Marshal(setting)

	setting.MaxLimit = 2E10
	setting.MinLimit = 2E8
	setting.TokenNames = []string{}
	resBytes3, _ := jsoniter.Marshal(setting)

	setting.TokenNames = []string{test.obj.sdk.Helper().GenesisHelper().Token().Name()}
	setting.MaxLimit = 0
	resBytes4, _ := jsoniter.Marshal(setting)

	setting.MaxLimit = 2E10
	setting.MinLimit = -1
	resBytes5, _ := jsoniter.Marshal(setting)

	setting.MinLimit = 2E8
	setting.MaxProfit = math.MinInt64
	resBytes6, _ := jsoniter.Marshal(setting)

	setting.MaxProfit = 2E12
	setting.FeeMiniNum = -1
	resBytes7, _ := jsoniter.Marshal(setting)

	setting.FeeMiniNum = 300000
	setting.FeeRatio = -1
	resBytes8, _ := jsoniter.Marshal(setting)

	setting.FeeRatio = 1001
	resBytes9, _ := jsoniter.Marshal(setting)

	setting.FeeRatio = 50
	setting.SendToCltRatio = -1
	resBytes10, _ := jsoniter.Marshal(setting)

	setting.SendToCltRatio = 1001
	resBytes11, _ := jsoniter.Marshal(setting)

	setting.SendToCltRatio = 100
	setting.BetExpirationBlocks = -1
	resBytes12, _ := jsoniter.Marshal(setting)

	var tests = []struct {
		account  sdk.IAccount
		settings []byte
		desc     string
		code     uint32
	}{
		{contractOwner, resBytes1, "--正常流程--", types.CodeOK},
		{contractOwner, resBytes2, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes3, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes4, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes5, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes6, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes7, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes8, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes9, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes10, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes11, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes12, "--异常流程--", types.ErrInvalidParameter},
		{accounts[0], resBytes1, "--异常流程--", types.ErrNoAuthorization},
	}

	test.run().setSender(contractOwner).InitChain()
	for _, item := range tests {
		utest.AssertError(test.run().setSender(item.account).SetSettings(string(item.settings)), item.code)
	}
}

//TestMyRoulette_SetRecFeeInfo This is a method of MySuite
func (mysuit *MySuite) TestMyRoulette_SetRecFeeInfo(c *check.C) () {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	recFeeInfo := make([]RecFeeInfo, 0)
	resBytes2, _ := jsoniter.Marshal(recFeeInfo)
	item := RecFeeInfo{
		RecFeeRatio: 500,
		RecFeeAddr:  "local9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE",
	}
	recFeeInfo = append(recFeeInfo, item)
	resBytes1, _ := jsoniter.Marshal(recFeeInfo)

	item1 := RecFeeInfo{
		RecFeeRatio: 501,
		RecFeeAddr:  "local9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE",
	}
	recFeeInfo = append(recFeeInfo, item1)
	resBytes3, _ := jsoniter.Marshal(recFeeInfo)

	recFeeInfo = append(recFeeInfo[:1], recFeeInfo[2:]...)
	item2 := RecFeeInfo{
		RecFeeRatio: 450,
		RecFeeAddr:  "lo9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE",
	}
	recFeeInfo = append(recFeeInfo, item2)
	resBytes4, _ := jsoniter.Marshal(recFeeInfo)

	recFeeInfo = append(recFeeInfo[:1], recFeeInfo[2:]...)
	item3 := RecFeeInfo{
		RecFeeRatio: 500,
		RecFeeAddr:  test.obj.sdk.Helper().BlockChainHelper().CalcAccountFromName(contractName,""),
	}
	recFeeInfo = append(recFeeInfo, item3)
	//resBytes5, _ := jsoniter.Marshal(recFeeInfo)

	recFeeInfo = append(recFeeInfo[:1], recFeeInfo[2:]...)
	item4 := RecFeeInfo{
		RecFeeRatio: -1,
		RecFeeAddr:  "local9ge366rtqV9BHqNwn7fFgA8XbDQmJGZqE",
	}
	recFeeInfo = append(recFeeInfo, item4)
	resBytes6, _ := jsoniter.Marshal(recFeeInfo)

	var tests = []struct {
		account sdk.IAccount
		infos   []byte
		desc    string
		code    uint32
	}{
		{contractOwner, resBytes1, "--正常流程--", types.CodeOK},
		{contractOwner, resBytes2, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes3, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes4, "--异常流程--", types.ErrInvalidAddress},
		//{contractOwner, resBytes5, "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, resBytes6, "--异常流程--", types.ErrInvalidParameter},
		{accounts[0], resBytes1, "--异常流程--", types.ErrNoAuthorization},
	}

	for _, item := range tests {
		utest.AssertError(test.run().setSender(item.account).SetRecFeeInfo(string(item.infos)), item.code)
	}
}

//TestMyRoulette_WithdrawFunds This is a method of MySuite
func (mysuit *MySuite) TestMyRoulette_WithdrawFunds(c *check.C) () {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	genesisToken := test.obj.sdk.Helper().GenesisHelper().Token()
	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	contractAccount := utest.UTP.Helper().ContractHelper().ContractOfName(contractName).Account()

	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)

	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), contractAccount, bn.N(1E11))
	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	test.run().setSender(contractOwner).InitChain()

	var tests = []struct {
		account        sdk.IAccount
		tokenName      string
		beneficiary    types.Address
		withdrawAmount bn.Number
		desc           string
		code           uint32
	}{
		{contractOwner, genesisToken.Name(), contractOwner.Address(), bn.N(1E10), "--正常流程--", types.CodeOK},
		{contractOwner, genesisToken.Name(), accounts[0].Address(), bn.N(1E10), "--正常流程--", types.CodeOK},
		{contractOwner, genesisToken.Name(), contractOwner.Address(), bn.N(1E15), "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, genesisToken.Name(), contractOwner.Address(), bn.N(-1), "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, genesisToken.Name(), contractAccount, bn.N(1E10), "--异常流程--", types.ErrInvalidParameter},
		{contractOwner, "xt", contractOwner.Address(), bn.N(1E10), "--异常流程--", types.ErrInvalidParameter},
		{accounts[0], genesisToken.Name(), contractOwner.Address(), bn.N(1E10), "--异常流程--", types.ErrNoAuthorization},
	}

	for _, item := range tests {
		utest.AssertError(test.run().setSender(item.account).WithdrawFunds(item.tokenName, item.beneficiary, item.withdrawAmount), item.code)
	}
}

//TestMyRoulette_PlaceBet This is a method of MySuite
func (mysuit *MySuite) TestMyRoulette_PlaceBet(c *check.C) () {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	//test.setSender(contractOwner).InitChain()
	//TODO


	contract := utest.UTP.Message().Contract()
	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)

	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), contract.Account(), bn.N(1E11))
	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 6)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	commitLastBlock, pubKey, _, commit, signData, _ := PlaceBetHelper(100)
	//utest.AssertError(err, types.CodeOK)

	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetPublicKey(pubKey[:]), types.CodeOK)

	betData := []BetData{{1, []int64{},bn.N(1000000000)}}
	betData1 := []BetData{{TWODIGITCOMBINATION, []int64{1,2},bn.N(1000000000)}}
	betData2 := []BetData{{THREEDIGITCOMBINATION, []int64{1,2,3},bn.N(1000000000)}}
	betData3 := []BetData{{FOURDIGITCOMBINATION, []int64{1,2,4,5},bn.N(1000000000)}}
	betData4 := []BetData{{1, []int64{},bn.N(1000000000)},{TWODIGITCOMBINATION, []int64{1,2},bn.N(1000000000)},{THREEDIGITCOMBINATION, []int64{1,2,3},bn.N(1000000000)},{FOURDIGITCOMBINATION, []int64{1,2,4,5},bn.N(1000000000)}}
	betData5 := []BetData{{THREEDIGITCOMBINATION, []int64{0,37,2},bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	betDataJsonBytes1, _ := jsoniter.Marshal(betData1)
	betDataJsonBytes2, _ := jsoniter.Marshal(betData2)
	betDataJsonBytes3, _ := jsoniter.Marshal(betData3)
	betDataJsonBytes4, _ := jsoniter.Marshal(betData4)
	betDataJsonBytes5, _ := jsoniter.Marshal(betData5)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes2), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes3), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[4]).transfer(bn.N(4000000000)).PlaceBet(string(betDataJsonBytes4), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[5]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes5), commitLastBlock, commit, signData[:], ""), types.CodeOK)
}

func Load(keystorePath string, password, fingerprint []byte) (acct *keys.Account, err types.Error) {
	if keystorePath == "" {
		common.PanicSanity("Cannot loads account because keystorePath not set")
	}

	walBytes, mErr := ioutil.ReadFile(keystorePath)
	if mErr != nil {
		err.ErrorCode = types.ErrInvalidParameter
		err.ErrorDesc = "account does not exist"
		return
	}

	jsonBytes, mErr := algorithm.DecryptWithPassword(walBytes, password, fingerprint)
	if mErr != nil {
		err.ErrorCode = types.ErrInvalidParameter
		err.ErrorDesc = fmt.Sprintf("the password is wrong err info : %s", mErr)
		return
	}

	acct = new(keys.Account)
	mErr = cdc.UnmarshalJSON(jsonBytes, acct)
	if mErr != nil {
		err.ErrorCode = types.ErrInvalidParameter
		err.ErrorDesc = fmt.Sprintf("UnmarshalJSON is wrong err info : %s", mErr)
		return
	}

	acct.KeystorePath = keystorePath
	err.ErrorCode = types.CodeOK
	return
}

//TestMyRoulette_SettleBet is a method of MySuite
func (mysuit *MySuite) TestMyRoulette_SettleBet(c *check.C) () {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	//test.setSender(contractOwner).InitChain()
	//TODO


	contract := utest.UTP.Message().Contract()
	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)

	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), contract.Account(), bn.N(1E11))
	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 5)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	commitLastBlock, pubKey, reveal, commit, signData, _ := PlaceBetHelper(100)
	//utest.AssertError(err, types.CodeOK)

	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetPublicKey(pubKey[:]), types.CodeOK)

	betData := []BetData{{1, []int64{},bn.N(1000000000)}}
	betData1 := []BetData{{TWODIGITCOMBINATION, []int64{1,2},bn.N(1000000000)}}
	betData2 := []BetData{{THREEDIGITCOMBINATION, []int64{1,2,3},bn.N(1000000000)}}
	betData3 := []BetData{{FOURDIGITCOMBINATION, []int64{1,2,4,5},bn.N(1000000000)}}
	betData4 := []BetData{{1, []int64{},bn.N(1000000000)},{TWODIGITCOMBINATION, []int64{1,2},bn.N(1000000000)},{THREEDIGITCOMBINATION, []int64{1,2,3},bn.N(1000000000)},{FOURDIGITCOMBINATION, []int64{1,2,4,5},bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	betDataJsonBytes1, _ := jsoniter.Marshal(betData1)
	betDataJsonBytes2, _ := jsoniter.Marshal(betData2)
	betDataJsonBytes3, _ := jsoniter.Marshal(betData3)
	betDataJsonBytes4, _ := jsoniter.Marshal(betData4)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes2), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes3), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[4]).transfer(bn.N(4000000000)).PlaceBet(string(betDataJsonBytes4), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	//结算
	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 1), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 3), types.CodeOK)
}

//TestMyRoulette_RefundBets is a method of MySuite
func (mysuit *MySuite) TestMyRoulette_RefundBets(c *check.C) () {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)


	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)
	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), test.obj.sdk.Message().Contract().Account(), bn.N(1E11))
	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 2)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	commitLastBlock, pubKey, _, commit, signData, _ := PlaceBetHelper(100)

	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetPublicKey(pubKey[:]), types.CodeOK)


	betData := []BetData{{1, []int64{},bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	betData = []BetData{{1, []int64{},bn.N(1000000000)}}
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	// set bet time out
	count := 0
	for {
		utest.NextBlock(1)
		count++
		if count > 250 {
			break
		}
	}
	utest.AssertError(test.run().setSender(contractOwner).RefundBets(commit, 1), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).RefundBets(commit, 1), types.CodeOK)


}

//hempHeight 想对于下注高度和生效高度之间的差值
//acct 合约的owner
func PlaceBetHelper(tempHeight int64) (commitLastBlock int64, pubKey [32]byte, reveal, commit []byte, signData [64]byte, err types.Error) {
	acct, err := Load("./.keystore/local_owner.wal", []byte(password), nil)
	if err.ErrorCode != types.CodeOK {
		return
	}

	localBlockHeight := utest.UTP.ISmartContract.Block().Height()

	pubKey = acct.PubKey.(crypto.PubKeyEd25519)

	commitLastBlock = localBlockHeight + tempHeight
	decode := crypto.CRandBytes(32)
	revealStr := hex.EncodeToString(algorithm.SHA3256(decode))
	reveal, _ = hex.DecodeString(revealStr)

	commit = algorithm.SHA3256(reveal)

	signByte := append(bn.N(commitLastBlock).Bytes(), commit...)
	signData = acct.PrivKey.Sign(signByte).(crypto.SignatureEd25519)

	return
}

