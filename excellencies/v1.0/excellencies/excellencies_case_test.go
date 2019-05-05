package excellencies

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

//Test This is a function
func Test(t *testing.T) { check.TestingT(t) }

//MySuite This is a struct
type MySuite struct{}

var _ = check.Suite(&MySuite{})

//TestExcellencies_SetSecretSigner This is a method of MySuite
func (mysuit *MySuite) TestExcellencies_SetSecretSigner(c *check.C) {
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
		utest.AssertError(test.run().setSender(item.account).SetSecretSigner(item.pubKey), item.code)
	}

}

//TestExcellencies_SetOwner is a method of MySuite
func (mysuit *MySuite) TestExcellencies_SetOwner(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	account := utest.NewAccount(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1000000000))
	//ss :=account.Address()
	//print(ss)
	var tests = []struct {
		account sdk.IAccount
		addr    types.Address
		desc    string
		code    uint32
	}{
		{contractOwner, types.Address("local5wNLBx3qhuQoSuA5ZGsyw8Fj6dJGwdFww"), "--正常流程--", types.CodeOK},
		{contractOwner, types.Address("local5wNLBx3qhuQoSuA5ZGsyw8Fj6dJGwdFwwrrr"), "--异常流程--公钥长度不正确--", types.ErrInvalidParameter},
		{account, types.Address("local5wNLBx3qhuQoSuA5ZGsyw8Fj6dJGwdFwr"), "--异常流程--非owner调用--", types.ErrNoAuthorization},
	}

	//for _, item := range tests {
	//	utest.AssertError(test.run().setSender(item.account).SetOwner(item.addr), item.code)
	//}
	//utest.AssertError(test.run().setSender(tests[0].account).SetOwner(tests[0].addr), tests[0].code)
	//utest.AssertError(test.run().setSender(tests[1].account).SetOwner(tests[1].addr), tests[1].code)
	utest.AssertError(test.run().setSender(tests[2].account).SetOwner(tests[2].addr), tests[2].code)
}

//TestExcellencies_PlaceBet This is a method of MySuite
func (mysuit *MySuite) TestExcellencies_PlaceBet(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	//TODO
	contract := utest.UTP.Message().Contract()
	genesisOwner := utest.UTP.Helper().GenesisHelper().Token().Owner()
	utest.Assert(test.run().setSender(utest.UTP.Helper().AccountHelper().AccountOf(genesisOwner)) != nil)

	utest.Transfer(nil, test.obj.sdk.Helper().GenesisHelper().Token().Name(), contract.Account(), bn.N(1E11))
	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 5)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	commitLastBlock, pubKey, _, commit, signData, _ := PlaceBetHelper(100)
	//utest.AssertError(err, types.CodeOK)

	test.run().setSender(contractOwner).InitChain()
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)

	betData := []BetData{{"C", bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	betData1 := []BetData{{"B", bn.N(1000000000)}}
	betDataJsonBytes1, _ := jsoniter.Marshal(betData1)
	// PlaceBet(betInfoJson string, commitLastBlock int64,betIndex string, commit, signData []byte, refAddress types.Address)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
}

func (mysuit *MySuite) TestExcellencies_Game(c *check.C) {
	//utest.Init(orgID)
	//contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	//test := NewTestObject(contractOwner)
	////TODO
	//hash := test.obj.sdk.Block().BlockHash()
	//sprint := fmt.Sprint(hash)
	//fmt.Println(sprint)
	//_, _, reval, _, _, _ := PlaceBetHelper(100)
	////utest.AssertError(err, types.CodeOK)
	//game := test.obj.StartGame(reval)
	//cards1 := game.GetGamerCards("A")
	//cards2 := game.GetGamerCards("B")
	//cards3 := game.GetGamerCards("C")
	//cards4 := game.GetGamerCards("D")
	//cards5 := game.GetGamerCards("E")
	//fmt.Println(*cards1.Pokers[0], *cards1.Pokers[1], *cards1.Pokers[2], cards1.No)
	//fmt.Println(*cards2.Pokers[0], *cards2.Pokers[1], *cards2.Pokers[2], cards2.No)
	//fmt.Println(*cards3.Pokers[0], *cards3.Pokers[1], *cards3.Pokers[2], cards3.No)
	//fmt.Println(*cards4.Pokers[0], *cards4.Pokers[1], *cards4.Pokers[2], cards4.No)
	//fmt.Println(*cards5.Pokers[0], *cards5.Pokers[1], *cards5.Pokers[2], cards5.No)
}

//TestSicBo_SetSettings This is a method of MySuite
func (mysuit *MySuite) TestSicBo_SetSettings(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

	accounts := utest.NewAccounts(test.obj.sdk.Helper().GenesisHelper().Token().Name(), bn.N(1E13), 1)
	if accounts == nil {
		panic("初始化newOwner失败")
	}

	settings := Settings{}

	settings.FeeRatio = 50
	settings.FeeMiniNum = 300000
	settings.SendToCltRatio = 100
	settings.BetExpirationBlocks = 250
	settings.TokenNames = []string{test.obj.sdk.Helper().GenesisHelper().Token().Name()}
	limitMaps := make(map[string]Limit, len(settings.TokenNames))
	for _, value := range settings.TokenNames {
		limit := Limit{
			MaxProfit: 2E12,
			MaxLimit:  2E10,
			MinLimit:  1E8,
		}
		limitMaps[value] = limit
	}
	settings.LimitMaps = limitMaps
	resBytes1, _ := jsoniter.Marshal(settings)
	fmt.Println(string(resBytes1))
	for _, value := range settings.TokenNames {
		limit := Limit{
			MaxProfit: 2E12,
			MaxLimit:  2E9,
			MinLimit:  2E10,
		}
		limitMaps[value] = limit
	}
	settings.LimitMaps = limitMaps
	resBytes2, _ := jsoniter.Marshal(settings)

	for _, value := range settings.TokenNames {
		limit := Limit{
			MaxProfit: 2E12,
			MaxLimit:  2E10,
			MinLimit:  2E8,
		}
		limitMaps[value] = limit
	}
	settings.LimitMaps = limitMaps
	settings.TokenNames = []string{}
	resBytes3, _ := jsoniter.Marshal(settings)

	settings.TokenNames = []string{test.obj.sdk.Helper().GenesisHelper().Token().Name()}
	for _, value := range settings.TokenNames {
		limit := Limit{
			MaxProfit: 2E12,
			MaxLimit:  0,
			MinLimit:  2E10,
		}
		limitMaps[value] = limit
	}
	settings.LimitMaps = limitMaps
	resBytes4, _ := jsoniter.Marshal(settings)

	for _, value := range settings.TokenNames {
		limit := Limit{
			MaxProfit: 2E12,
			MaxLimit:  2E10,
			MinLimit:  -1,
		}
		limitMaps[value] = limit
	}
	settings.LimitMaps = limitMaps
	resBytes5, _ := jsoniter.Marshal(settings)

	for _, value := range settings.TokenNames {
		limit := Limit{
			MaxProfit: math.MinInt64,
			MaxLimit:  2E8,
			MinLimit:  2E10,
		}
		limitMaps[value] = limit
	}
	settings.LimitMaps = limitMaps
	resBytes6, _ := jsoniter.Marshal(settings)

	for _, value := range settings.TokenNames {
		limit := Limit{
			MaxProfit: 2E12,
			MaxLimit:  2E9,
			MinLimit:  2E10,
		}
		limitMaps[value] = limit
	}
	settings.LimitMaps = limitMaps
	settings.FeeMiniNum = -1
	resBytes7, _ := jsoniter.Marshal(settings)

	settings.FeeMiniNum = 300000
	settings.FeeRatio = -1
	resBytes8, _ := jsoniter.Marshal(settings)

	settings.FeeRatio = 1001
	resBytes9, _ := jsoniter.Marshal(settings)

	settings.FeeRatio = 50
	settings.SendToCltRatio = -1
	resBytes10, _ := jsoniter.Marshal(settings)

	settings.SendToCltRatio = 1001
	resBytes11, _ := jsoniter.Marshal(settings)

	settings.SendToCltRatio = 100
	settings.BetExpirationBlocks = -1
	fmt.Println()
	fmt.Println("setting setting setting setting")
	fmt.Println(settings)
	fmt.Println()
	fmt.Println()
	fmt.Println()

	resBytes12, _ := jsoniter.Marshal(settings)

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

//TestExcellencies_SetSettings is a method of MySuite
func (mysuit *MySuite) TestExcellencies_GetGameResult(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	test.setSender(contractOwner).InitChain()

	hash := test.obj.sdk.Block().BlockHash()
	sprint := fmt.Sprint(hash)
	fmt.Println(sprint)
	//game := test.obj.GetGameResult()
	//cards1 := game.GetGamerCards("A")
	//cards2 := game.GetGamerCards("B")
	//cards3 := game.GetGamerCards("C")
	//cards4 := game.GetGamerCards("D")
	//cards5 := game.GetGamerCards("E")
	//fmt.Println(cards1.No, *cards1.Pokers[0], *cards1.Pokers[1], *cards1.Pokers[2])
	//fmt.Println(cards2.No, *cards2.Pokers[0], *cards2.Pokers[1], *cards2.Pokers[2], cards2.IsWin)
	//fmt.Println(cards3.No, *cards3.Pokers[0], *cards3.Pokers[1], *cards3.Pokers[2], cards3.IsWin)
	//fmt.Println(cards4.No, *cards4.Pokers[0], *cards4.Pokers[1], *cards4.Pokers[2], cards4.IsWin)
	//fmt.Println(cards5.No, *cards5.Pokers[0], *cards5.Pokers[1], *cards5.Pokers[2], cards5.IsWin)
}

//TestExcellencies_SetRecFeeInfo is a method of MySuite
func (mysuit *MySuite) TestExcellencies_SetRecFeeInfo(c *check.C) {
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

//TestExcellencies_SetSettings is a method of MySuite
func (mysuit *MySuite) TestExcellencies_SetSettings(c *check.C) {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	test.setSender(contractOwner).InitChain()

}

//TestExcellencies_WithdrawFunds is a method of MySuite
func (mysuit *MySuite) TestExcellencies_WithdrawFunds(c *check.C) {
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
		//{contractOwner, genesisToken.Name(), contractAccount, bn.N(1E10), "--异常流程--", types.ErrInvalidParameter},
		//{contractOwner, "xt", contractOwner.Address(), bn.N(1E10), "--异常流程--", types.ErrInvalidParameter},
		{accounts[0], genesisToken.Name(), contractOwner.Address(), bn.N(1E10), "--异常流程--", types.ErrNoAuthorization},
	}

	for _, item := range tests {
		utest.AssertError(test.run().setSender(item.account).WithdrawFunds(item.tokenName, item.beneficiary, item.withdrawAmount), item.code)
	}
}

//TestExcellencies_SettleBet is a method of MySuite
func (mysuit *MySuite) TestExcellencies_SettleBet(c *check.C) {

	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

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
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)

	betData := []BetData{{"C", bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	betData1 := []BetData{{"B", bn.N(1000000000)}}
	betDataJsonBytes1, _ := jsoniter.Marshal(betData1)

	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 1), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 4), types.CodeOK)
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

//TestExcellencies_WithdrawWin is a method of MySuite
func (mysuit *MySuite) TestExcellencies_WithdrawWin(c *check.C) () {
	//utest.Init(orgID)
	//contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	//test := NewTestObject(contractOwner)
	//test.setSender(contractOwner).InitChain()
	////TODO

	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

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
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)

	betData := []BetData{{"C", bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	betData1 := []BetData{{"B", bn.N(1000000000)}}
	betDataJsonBytes1, _ := jsoniter.Marshal(betData1)

	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)

	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 1), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).WithdrawWin(commit), types.CodeOK)
}

//TestExcellencies_SlipperHighestTransfer is a method of MySuite
func (mysuit *MySuite) TestExcellencies_SlipperHighestTransfer(c *check.C) () {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

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
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)

	betData := []BetData{{"C", bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	betData1 := []BetData{{"B", bn.N(1000000000)}}
	betDataJsonBytes1, _ := jsoniter.Marshal(betData1)

	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 1), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 4), types.CodeOK)


	//utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	//utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 1), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SlipperHighestTransfer(test.obj.sdk.Helper().GenesisHelper().Token().Name(),accounts[0].Address()), types.ErrInvalidParameter)
	//utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).SlipperHighestTransfer(test.obj.sdk.Helper().GenesisHelper().Token().Name(),accounts[0].Address()), types.ErrInvalidParameter)

}

//TestExcellencies_GetPoolAmount is a method of MySuite
func (mysuit *MySuite) TestExcellencies_GetPoolAmount(c *check.C) () {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)
	test.setSender(contractOwner).InitChain()
	//TODO
}

//TestExcellencies_CarveUpPool is a method of MySuite
func (mysuit *MySuite) TestExcellencies_CarveUpPool(c *check.C) () {
	utest.Init(orgID)
	contractOwner := utest.DeployContract(c, contractName, orgID, contractMethods, contractInterfaces)
	test := NewTestObject(contractOwner)

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
	utest.AssertError(test.run().setSender(contractOwner).SetSecretSigner(pubKey[:]), types.CodeOK)

	betData := []BetData{{"C", bn.N(1000000000)}}
	betDataJsonBytes, _ := jsoniter.Marshal(betData)
	betData1 := []BetData{{"B", bn.N(1000000000)}}
	betData2 := []BetData{{"D", bn.N(1000000000)}}
	betData3 := []BetData{{"E", bn.N(1000000000)}}
	betDataJsonBytes1, _ := jsoniter.Marshal(betData1)
	betDataJsonBytes2, _ := jsoniter.Marshal(betData2)
	betDataJsonBytes3, _ := jsoniter.Marshal(betData3)

	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes3), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes2), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes2), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes3), commitLastBlock, commit, signData[:], ""), types.CodeOK)

	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes1), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes3), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[0]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes2), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[1]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[2]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes2), commitLastBlock, commit, signData[:], ""), types.CodeOK)
	utest.AssertError(test.run().setSender(accounts[3]).transfer(bn.N(1000000000)).PlaceBet(string(betDataJsonBytes3), commitLastBlock, commit, signData[:], ""), types.CodeOK)

	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 1), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 4), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 1), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).SettleBet(reveal, 4), types.CodeOK)
	utest.AssertError(test.run().setSender(contractOwner).CarveUpPool(reveal), types.CodeOK)

}