package excellencies

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/crypto/ed25519"
	"blockchain/smcsdk/sdk/crypto/sha3"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/types"
	"encoding/hex"
)

//excellencies This is struct of contract
//@:contract:excellencies
//@:version:1.0
//@:organization:orgNUjCm1i8RcoW2kVTbDw4vKW6jzfMxewJHjkhuiduhjuikjuyhnnjkuhujk111
//@:author:ef94556a937618c72ffaf173b1533c533d77aa3ea2a63f053bb904feefe5a92f
type Excellencies struct {
	sdk sdk.ISmartContract

	//@:public:store:cache
	secretSigner types.PubKey // Check to sign the public key

	//@:public:store:cache
	lockedInBets map[string]bn.Number // Lock amount (unit cong) key: currency name

	//@:public:store:cache
	settings *Settings

	//@:public:store:cache
	recFeeInfo []RecFeeInfo

	//@:public:store
	roundInfo map[string]*RoundInfo //key1轮标识

	//@:public:store:cache
	poolAmount map[string]bn.Number

	//@:public:store:cache
	slipper map[string]SlipperInfo

	//@:public:store
	poker  []*Poker

}

//@:public:receipt
type receipt interface {
	emitSetSecretSigner(newSecretSigner types.PubKey)
	emitSetOwner(newAddress types.Address)
	emitSetSettings(tokenNames []string, limit map[string]Limit, feeRatio, feeMiniNum, sendToCltRatio, betExpirationBlocks,CarveUpPoolRatio int64)
	emitPlaceBet(tokenName string, gambler types.Address, totalMaybeWinAmount bn.Number, betDataList []BetData, commitLastBlock int64, commit, signData []byte, refAddress types.Address)
	emitSetRecFeeInfo(info []RecFeeInfo)
	emitWithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number)
	emitSettleBet(tokenName []string, commit []byte, game Game, winNumber int64, totalWinAmount map[string]bn.Number, finished bool,poolAmount map[string]bn.Number,OriginPokers []*Poker)
	emitSlipperHighestTransfer(tokenName string, playersAddress types.Address, winningAmount bn.Number)
	emitCarveUpPool(countNum int64,amountPoolList map[string]bn.Number)
}

//InitChain Constructor of this excellencies
//@:constructor
func (sg *Excellencies) InitChain() {

	settings := Settings{}

	settings.FeeRatio = 50
	settings.FeeMiniNum = 300000
	settings.SendToCltRatio = 100
	settings.BetExpirationBlocks = 250
	settings.TokenNames = []string{sg.sdk.Helper().GenesisHelper().Token().Name()}
	settings.PoolFeeRatio = 20
	settings.CarveUpPoolRatio = 50

	limitMaps := sg.createLimitMaps(settings.TokenNames)
	settings.LimitMaps = limitMaps
	sg._setSettings(&settings)
	sg.LockedInBetsInit(settings.TokenNames)
	sg.PoolAmountInit(settings.TokenNames)
	//set poker
	pokers:=[]*Poker{
		&Poker{"2",2,0},&Poker{"2",2,1},&Poker{"2",2,2},&Poker{"2",2,3},
		&Poker{"3",3,0},&Poker{"3",3,1},&Poker{"3",3,2},&Poker{"3",3,3},
		&Poker{"4",4,0},&Poker{"4",4,1},&Poker{"4",4,2},&Poker{"4",4,3},
		&Poker{"5",5,0},&Poker{"5",5,1},&Poker{"5",5,2},&Poker{"5",5,3},
		&Poker{"6",6,0},&Poker{"6",6,1},&Poker{"6",6,2},&Poker{"6",6,3},
		&Poker{"7",7,0},&Poker{"7",7,1},&Poker{"7",7,2},&Poker{"7",7,3},
		&Poker{"8",8,0},&Poker{"8",8,1},&Poker{"8",8,2},&Poker{"8",8,3},
		&Poker{"9",9,0},&Poker{"9",9,1},&Poker{"9",9,2},&Poker{"9",9,3},
		&Poker{"10",10,0},&Poker{"10",10,1},&Poker{"10",10,2},&Poker{"10",10,3},
		&Poker{"A",1,0},&Poker{"A",1,1},&Poker{"A",1,2},&Poker{"A",1,3},
		&Poker{"J",11,0},&Poker{"J",11,1},&Poker{"J",11,2},&Poker{"J",11,3},
		&Poker{"Q",12,0},&Poker{"Q",12,1},&Poker{"Q",12,2},&Poker{"Q",12,3},
		&Poker{"K",13,0},&Poker{"K",13,1},&Poker{"K",13,2},&Poker{"K",13,3},
		}
	sg._setPoker(pokers)

}

func (e *Excellencies) PoolAmountInit(tokenNameList []types.Address) {
	for _, value := range tokenNameList {
		e._setPoolAmount(value, bn.N(0))
		sinfo:= SlipperInfo{
			SettlementTime:e.sdk.Block().Now(),
			PlayersAddress:types.Address(""),
			TokenName:value,
			SettlementMoney:bn.N(0),
		}
		e._setSlipper(value,sinfo)
	}
}

//SetSecretSigner - Set up the public key
//@:public:method:gas[500]
func (sg *Excellencies) SetSecretSigner(newSecretSigner types.PubKey) {

	sdk.RequireOwner(sg.sdk)
	sdk.Require(len(newSecretSigner) == 32,
		types.ErrInvalidParameter, "length of newSecretSigner must be 32 bytes")

	//Save to database
	sg._setSecretSigner(newSecretSigner)

	// fire event
	sg.emitSetSecretSigner(newSecretSigner)
}

//SetOwner -Set contract owner
//@:public:method:gas[500]
func (sg *Excellencies) SetOwner(newOwnerAddr types.Address) {
	//only contract owner just can set new owner
	sdk.RequireOwner(sg.sdk)

	sdk.Require(newOwnerAddr != sg.sdk.Message().Contract().Account(),
		types.ErrInvalidParameter, "NewOwner address cannot be contract account address")

	sdk.Require(newOwnerAddr != sg.sdk.Message().Contract().Address(),
		types.ErrInvalidParameter, "NewOwner address cannot be contract address")

	sdk.Require(sg.sdk.Helper().BlockChainHelper().CheckAddress(newOwnerAddr) == nil,
		types.ErrInvalidParameter, "NewOwner must be valid address")

	sg.sdk.Message().Contract().SetOwner(newOwnerAddr)

	//fire event
	sg.emitSetOwner(sg.sdk.Message().Contract().Owner())
}

//SetSettings - Change game settings
//@:public:method:gas[500]
func (sg *Excellencies) SetSettings(newSettinsStr string) {
	sdk.RequireOwner(sg.sdk)

	//check that the Settings are valid
	newSettings := new(Settings)
	err := jsoniter.Unmarshal([]byte(newSettinsStr), newSettings)
	sdk.RequireNotError(err, types.ErrInvalidParameter)
	sg.checkSettings(newSettings)

	//Settings can only be set after all settlement is completed and the refund is completed
	settings := sg._settings()
	for _, tokenName := range settings.TokenNames {
		lockedAmount := sg._lockedInBets(tokenName)
		sdk.Require(lockedAmount.CmpI(0) == 0,
			types.ErrUserDefined, "only lockedAmount is zero that can do SetSettings()")
	}
	sg._setSettings(newSettings)
	// fire event
	sg.emitSetSettings(
		newSettings.TokenNames,
		newSettings.LimitMaps,
		newSettings.FeeRatio,
		newSettings.FeeMiniNum,
		newSettings.SendToCltRatio,
		newSettings.BetExpirationBlocks,
		newSettings.CarveUpPoolRatio,
	)
}

// SetRecFeeInfo - Set ratio of fee and receiver's account address
//@:public:method:gas[500]
func (sg *Excellencies) SetRecFeeInfo(recFeeInfoStr string) {

	sdk.RequireOwner(sg.sdk)

	info := make([]RecFeeInfo, 0)
	err := jsoniter.Unmarshal([]byte(recFeeInfoStr), &info)
	sdk.RequireNotError(err, types.ErrInvalidParameter)
	//Check that the parameters are valid
	sg.checkRecFeeInfo(info)

	sg._setRecFeeInfo(info)
	// fire event
	sg.emitSetRecFeeInfo(info)
}

// WithdrawFunds - Funds withdrawal
//@:public:method:gas[500]
func (sg *Excellencies) WithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) {

	sdk.RequireOwner(sg.sdk)
	sdk.Require(withdrawAmount.CmpI(0) > 0,
		types.ErrInvalidParameter, "withdrawAmount must be larger than zero")

	account := sg.sdk.Helper().AccountHelper().AccountOf(sg.sdk.Message().Contract().Account())
	lockedAmount := sg._lockedInBets(tokenName)
	unlockedAmount := account.BalanceOfName(tokenName).Sub(lockedAmount)
	sdk.Require(unlockedAmount.Cmp(withdrawAmount) >= 0,
		types.ErrInvalidParameter, "Not enough funds")

	// transfer to beneficiary
	account.TransferByName(tokenName, beneficiary, withdrawAmount)

	// fire event
	sg.emitWithdrawFunds(tokenName, beneficiary, withdrawAmount)
}

//PlaceBet - place bet
//@:public:method:gas[500]
func (sg *Excellencies) PlaceBet(betJson string, commitLastBlock int64, commit, signData []byte, refAddress types.Address) {

	sdk.Require(sg.sdk.Message().Sender().Address() != sg.sdk.Message().Contract().Owner(),
		types.ErrNoAuthorization, "The contract owner cannot bet")

	//1. Verify whether the signature and current round betting are legal
	data := append(bn.N(commitLastBlock).Bytes(), commit...)
	sdk.Require(ed25519.VerifySign(sg._secretSigner(), data, signData),
		types.ErrInvalidParameter, "Incorrect signature")

	gambler := sg.sdk.Message().Sender().Address()
	hexCommit := hex.EncodeToString(commit)
	//Is late
	sdk.Require(sg.sdk.Block().Height() <= commitLastBlock,
		types.ErrInvalidParameter, "Commit has expired")

	settings := sg._settings()
	if !sg._chkRoundInfo(hexCommit) {
		info := NewRoundInfo([]byte(hexCommit), *settings)
		info.FirstBlockHeight=sg.sdk.Block().Height()
		sg._setRoundInfo(hexCommit, info)
	}
	roundInfo := sg._roundInfo(hexCommit)
	sdk.Require(NOAWARD == roundInfo.State,
		types.ErrInvalidParameter, "No betting on the current wheel")
	sdk.Require(roundInfo.FirstBlockHeight+roundInfo.Settings.BetExpirationBlocks > sg.sdk.Block().Height(),
		types.ErrInvalidParameter, "This round is time out")

	tokenName := ""
	transferAmount := bn.N(0)
	//Obtain transfer records
	for _, name := range settings.TokenNames {
		transferReceipt := sg.sdk.Message().GetTransferToMe(name)
		if transferReceipt != nil {
			tokenName = name
			transferAmount = transferReceipt.Value
			break
		}
	}

	//Verify receipt of bet transfer
	sdk.Require(tokenName != "" && transferAmount.CmpI(0) > 0,
		types.ErrUserDefined, "Must transfer tokens to me before place a bet")

	betDataList := BuildBetData([]byte(betJson))
	totalAmount := bn.N(0)
	//amountByModel:=make(map[string]map[string]bn.Number)
	for _, betData := range betDataList {
		amount := betData.BetAmount
		//Verify that the betting scheme is legal
		sdk.Require(betData.BetMode == PLAYER_1 || betData.BetMode == PLAYER_2 || betData.BetMode == PLAYER_3 || betData.BetMode == PLAYER_4,
			types.ErrInvalidParameter, "The betting range is illegal")
		//Verify that the bet amount is legal
		sdk.Require(amount.CmpI(settings.LimitMaps[tokenName].MinLimit) >= 0 && amount.CmpI(settings.LimitMaps[tokenName].MaxLimit) <= 0,
			types.ErrInvalidParameter, "Amount should be within range")
		if roundInfo.BetAmountByModel==nil{
			roundInfo.BetAmountByModel=make(map[string]map[string]bn.Number)
		}
		if _,ok:=roundInfo.BetAmountByModel[betData.BetMode];!ok{
			roundInfo.BetAmountByModel[betData.BetMode]=make(map[string]bn.Number)
			//amountByModel[betData.BetMode]=tmp
			if _,ok=roundInfo.BetAmountByModel[betData.BetMode][tokenName];!ok{
				roundInfo.BetAmountByModel[betData.BetMode][tokenName]=bn.N(0)
			}
		}
		roundInfo.BetAmountByModel[betData.BetMode][tokenName]=roundInfo.BetAmountByModel[betData.BetMode][tokenName].Add(amount)
		totalAmount = totalAmount.Add(amount)
	}

	//Check whether the total amount of bet is equal to the transfer amount
	sdk.Require(totalAmount.Cmp(transferAmount) == 0,
		types.ErrUserDefined, "transfer amount not equal place bet amount")

	//Verify that the reimbursable amount is sufficient
	//Calculate the amount of money you can win
	totalMaybeWinAmount, feeAmount := sg.MaybeWinAmountAndFeeByList(totalAmount, betDataList)
	//Is the amount likely to be won less than or equal to the maximum bonus amount
	sdk.Require(totalMaybeWinAmount.CmpI(settings.LimitMaps[tokenName].MaxProfit) <= 0,
		types.ErrInvalidParameter, "MaxProfit limit violation")
	contractAcct := sg.sdk.Helper().AccountHelper().AccountOf(sg.sdk.Message().Contract().Account())
	//Lock in the amount that may need to be paid
	totalLockedAmount := sg._lockedInBets(tokenName).Add(totalMaybeWinAmount).Add(feeAmount)
	//Contract account balance
	totalUnlockedAmount := contractAcct.BalanceOfName(tokenName)
	//Is the contract account balance greater than or equal to the balance that may need to be paid
	sdk.Require(totalUnlockedAmount.Cmp(totalMaybeWinAmount) >= 0,
		types.ErrInvalidParameter, "Cannot afford to lose this bet")
	sg._setLockedInBets(tokenName, totalLockedAmount)
	betInfo := NewBetInfo(tokenName, gambler)
	betInfo.BetData = betDataList
	betInfo.Amount = totalAmount
	if _, ok := roundInfo.TotalBuyAmount[tokenName]; !ok {
		roundInfo.TotalBuyAmount[tokenName] = bn.N(0)
	}
	roundInfo.TotalBuyAmount[tokenName] = roundInfo.TotalBuyAmount[tokenName].Add(totalAmount)
	roundInfo.TotalBetCount += 1
	roundInfo.BetInfos = append(roundInfo.BetInfos, *betInfo)
	sg._setRoundInfo(hexCommit, roundInfo)

	sg.emitPlaceBet(tokenName, gambler, totalMaybeWinAmount, betDataList, commitLastBlock, commit, signData, refAddress)
}

// SettleBet - The lottery and settlement
//@:public:method:gas[500]
func (e *Excellencies) SettleBet(reveal []byte, settleCount int64) {
	sdk.Require(len(reveal) > 0,
		types.ErrInvalidParameter, "Commit should be not exist")

	sdk.RequireOwner(e.sdk)
	hexCommit := hex.EncodeToString(sha3.Sum256(reveal))
	sdk.Require(e._chkRoundInfo(hexCommit), types.ErrInvalidParameter, "Commit should be not exist")

	roundInfo := e._roundInfo(hexCommit)
	//Current wheel configuration
	settings := roundInfo.Settings
	//The bet height of the round to be settled should be less than the settlement height
	sdk.Require(roundInfo.FirstBlockHeight < e.sdk.Block().Height(),
		types.ErrInvalidParameter, "SettleBet block can not be in the same block as placeBet, or before.")

	sdk.Require(NOAWARD == roundInfo.State || OPENINGAPRIZE == roundInfo.State,
		types.ErrInvalidParameter, "This state does not operate for settlement")

	//For the first settlement, the round information cannot expire
	if NOAWARD == roundInfo.State {
		sdk.Require(roundInfo.FirstBlockHeight+settings.BetExpirationBlocks > e.sdk.Block().Height(),
			types.ErrInvalidParameter, "This round is time out")
		roundInfo.Game = *e.GetGameResult(reveal,roundInfo)
		//No one bets on the current round
		if roundInfo.TotalBetCount <= 0 {
			roundInfo.State = AWARDED
			e._setRoundInfo(hexCommit, roundInfo)
			return
		}
		roundInfo.State = OPENINGAPRIZE
		e._setRoundInfo(hexCommit, roundInfo)
	}
	//Determine whether all bets have been settled
	if roundInfo.TotalBetCount == roundInfo.ProcessCount {
		roundInfo.State = AWARDED
		e._setRoundInfo(hexCommit, roundInfo)
		sdk.Require(false,
			types.ErrInvalidParameter, "This round is complete")
	}
	//Initial index
	startIndex := roundInfo.ProcessCount
	if startIndex < 0 {
		startIndex = 0
	}
	endIndex := startIndex + settleCount

	if endIndex >= roundInfo.TotalBetCount {
		endIndex = roundInfo.TotalBetCount
		//Set the database state state to lottery
		roundInfo.State = AWARDED
	}
	e.DealSettle(roundInfo, hexCommit, startIndex, endIndex)
}



//CarveUpPool - carve up pool
//@:public:method:gas[500]
func (e *Excellencies) CarveUpPool (commit []byte)  {
	sdk.RequireOwner(e.sdk)
	hexCommit := hex.EncodeToString(sha3.Sum256(commit))
	sdk.Require(e._chkRoundInfo(hexCommit), types.ErrInvalidParameter, "Commit should be not exist")

	roundInfo := e._roundInfo(hexCommit)
	setting:=roundInfo.Settings
	//The bet height of the round to be settled should be less than the settlement height
	sdk.Require(roundInfo.FirstBlockHeight < e.sdk.Block().Height(),
		types.ErrInvalidParameter, "SettleBet block can not be in the same block as placeBet, or before.")

	sdk.Require(false == roundInfo.PoolFlag,
		types.ErrInvalidParameter, "The operators is illegal")
	//Pool settlement
	amountPoolList:=make(map[string]bn.Number) //奖池
	for _,v:=range  e._settings().TokenNames {
		//拿出百分之5来派奖
		amountPoolList[v]=e._poolAmount(v).MulI(setting.CarveUpPoolRatio).DivI(PERMILLE)
	}


	modelAmountList:=roundInfo.BetAmountByModel //各个模型的投注总额
	tmp:=roundInfo.Game.Gamer
	winModel:=make([]string,0)//哪些model是符合派奖的模型
	if tmp[BETMODEL_B].TypeSan>THREENOSAME{ //表示是要派奖的类型
		winModel=append(winModel,BETMODEL_B)
	}
	if tmp[BETMODEL_C].TypeSan>THREENOSAME{ //表示是要派奖的类型
		winModel=append(winModel,BETMODEL_C)
	}
	if tmp[BETMODEL_D].TypeSan>THREENOSAME{ //表示是要派奖的类型
		winModel=append(winModel,BETMODEL_D)
	}
	if tmp[BETMODEL_E].TypeSan>THREENOSAME{ //表示是要派奖的类型
		winModel=append(winModel,BETMODEL_E)
	}
	sdk.Require(len(winModel)>=1,
		types.ErrInvalidParameter, "this operator is illegal")


	total:=make(map[string]bn.Number) //每个币种对应总的
	for _,v:=range winModel{ //v 代表模型
		for _,value:=range setting.TokenNames{//value 表示每个币种
		tmp:=modelAmountList[v][value]
			if _,ok:=total[value];!ok{
				total[value]=bn.N(0)
			}
			total[value]=total[value].Add(tmp)
		}

	}

	everyToken:=make(map[string]bn.Number)
	for _,value:=range setting.TokenNames{//value 表示每个币种
		everyToken[value]=amountPoolList[value].MulI(PERMILLE).Div(total[value])
	}

	contractAcct := e.sdk.Helper().AccountHelper().AccountOf(e.sdk.Message().Contract().Account())
	countNum:=int64(0)
	for _,v:=range roundInfo.BetInfos{ //v 每个人的投注
		//查出这个模型下的所有投注
		for _,value:=range v.BetData{
			if CheckList(value.BetMode,winModel){
				//表示投注的是符合派奖的牌型  给他派奖
				contractAcct.TransferByName(v.TokenName, v.Gambler,value.BetAmount.Mul(everyToken[v.TokenName]).DivI(PERMILLE))
				countNum++
			}
		}

	}
	roundInfo.PoolFlag=true
	e._setRoundInfo(hexCommit,roundInfo)
	e.emitCarveUpPool(countNum,amountPoolList)

}


//WithdrawWin - Player settlement
//@:public:method:gas[500]
func (e *Excellencies) WithdrawWin (commit []byte)  {
	hexCommit := hex.EncodeToString(commit)
	sdk.Require(e._chkRoundInfo(hexCommit), types.ErrInvalidParameter, "Commit should be not exist")

	roundInfo := e._roundInfo(hexCommit)
	//The bet height of the round to be settled should be less than the settlement height
	sdk.Require(roundInfo.FirstBlockHeight < e.sdk.Block().Height(),
		types.ErrInvalidParameter, "SettleBet block can not be in the same block as placeBet, or before.")

	sdk.Require(AWARDED != roundInfo.State,
		types.ErrInvalidParameter, "The operators have yet to draw a prize")

	totalBetCount := roundInfo.TotalBetCount
	sdk.Require(totalBetCount != roundInfo.ProcessCount,
		types.ErrInvalidParameter, "This round is complete")

	startIndex := int64(0)
	e.DealSettle(roundInfo, hexCommit, startIndex, totalBetCount)
}

// SlipperHighestTransfer - Slipper Highest Transfer
//@:public:method:gas[500]
func (e *Excellencies) SlipperHighestTransfer(tokenName string, playersAddress types.Address){
	sdk.RequireOwner(e.sdk)
	sdk.Require(e.sdk.Message().Sender().Address() == e.sdk.Message().Contract().Owner(),
		types.ErrNoAuthorization, "The contract owner can do")

	poolAmount := e._poolAmount(tokenName)
	winningAmount := poolAmount.DivI(2)
	sdk.Require(winningAmount.CmpI(0) > 0,
		types.ErrInvalidParameter, "winningAmount must be larger than zero")

	//account := bj.sdk.Helper().AccountHelper().AccountOf(bj.sdk.Message().Contract().Account())
	nowTime := e.sdk.Block().Now()
	//获得上一次返奖信息
	slipperInfo := e._slipper(tokenName)
	oldTime := slipperInfo.SettlementTime
	timeDifference := nowTime.Sub(oldTime)
	sdk.Require(timeDifference.IsGE(bn.N(ONEDAYSECONDS)),
		types.ErrInvalidParameter, "The call time is less than 24 hours")
	newslipperInfo:=SlipperInfo{
		SettlementTime:e.sdk.Block().Now(),
		PlayersAddress:playersAddress,
		TokenName:tokenName,
		SettlementMoney:winningAmount,
	}
	e.TransferToPlayers(tokenName, playersAddress, winningAmount)
	//account.TransferByToken(tokenName, playersAddress, winningAmount)
	e._setPoolAmount(tokenName, winningAmount)
	e._setSlipper(tokenName,newslipperInfo)
	e.emitSlipperHighestTransfer(tokenName, playersAddress, winningAmount)
}