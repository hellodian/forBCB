package myroulette

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/crypto/ed25519"
	"blockchain/smcsdk/sdk/crypto/sha3"
	_ "blockchain/smcsdk/sdk/crypto/sha3"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/types"
	"encoding/hex"
)

//@:public:receipt
type receipt interface {
	emitSetPublicKey(newPublicKey types.PubKey)
	emitSetSetting(tokenNames []string, minLimit, maxLimit, maxProfit, feeRatio, feeMiniNum, sendToCltRatio, betExpirationBlocks int64)
	emitSetRecFeeInfo(info []RecFeeInfo)
	emitWithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number)
	emitPlaceBet(tokenName string, gambler types.Address, totalMaybeWinAmount bn.Number, betDataList []BetData, commitLastBlock int64, commit, signData []byte, refAddress types.Address)
	emitSettleBet(tokenName []string, reveal, commit []byte, gambler []types.Address, winNumber int64, totalWinAmount map[string]bn.Number, finished bool)
	emitRefundBet(commit []byte, tokenName []string, gambler []types.Address, refundedAmount map[string]bn.Number, finished bool)
}

//InitChain - InitChain Constructor of this Roulette
//@:constructor
func (mr *MyRoulette) InitChain() {

	//init dataN
	setting := Setting{}
	setting.MaxProfit = 2E12
	setting.MaxLimit = 2E10
	setting.MinLimit = 1E8
	setting.FeeRatio = 50
	setting.FeeMiniNum = 300000
	setting.SendToCltRatio = 100
	setting.BetExpirationBlocks = 250
	setting.TokenNames = []string{mr.sdk.Helper().GenesisHelper().Token().Name()}

	mr._setSetting(&setting)
	mr.LockedInBetsInit(setting.TokenNames)
}

// SetPublicKey - Set up the public key
//@:public:method:gas[500]
func (mr *MyRoulette) SetPublicKey(newPublicKey types.PubKey) {

	sdk.RequireOwner(mr.sdk)
	sdk.Require(len(newPublicKey) == 32,
		types.ErrInvalidParameter, "length of newPublicKey must be 32 bytes")

	//Save to database
	mr._setPublicKey(newPublicKey)

	// fire event
	mr.emitSetPublicKey(newPublicKey)
}

// SetSettings - Change game settings
//@:public:method:gas[500]
func (mr *MyRoulette) SetSettings(newSettingsStr string) {

	sdk.RequireOwner(mr.sdk)

	//Check that the Settings are valid
	newSetting := new(Setting)
	err := jsoniter.Unmarshal([]byte(newSettingsStr), newSetting)
	sdk.RequireNotError(err, types.ErrInvalidParameter)
	mr.checkSettings(newSetting)

	//Settings can only be set after all settlement is completed and the refund is completed
	setting := mr._setting()
	for _, tokenName := range setting.TokenNames {
		lockedAmount := mr._lockedInBets(tokenName)
		sdk.Require(lockedAmount.CmpI(0) == 0,
			types.ErrUserDefined, "only lockedAmount is zero that can do SetSettings()")
	}

	mr._setSetting(newSetting)

	// fire event
	mr.emitSetSetting(
		newSetting.TokenNames,
		newSetting.MinLimit,
		newSetting.MaxLimit,
		newSetting.MaxProfit,
		newSetting.FeeRatio,
		newSetting.FeeMiniNum,
		newSetting.SendToCltRatio,
		newSetting.BetExpirationBlocks,
	)
}

// SetRecFeeInfo - Set ratio of fee and receiver's account address
//@:public:method:gas[500]
func (mr *MyRoulette) SetRecFeeInfo(recFeeInfoStr string) {

	sdk.RequireOwner(mr.sdk)

	info := make([]RecFeeInfo, 0)
	err := jsoniter.Unmarshal([]byte(recFeeInfoStr), &info)
	sdk.RequireNotError(err, types.ErrInvalidParameter)
	//Check that the parameters are valid
	mr.checkRecFeeInfo(info)

	mr._setRecFeeInfo(info)
	// fire event
	mr.emitSetRecFeeInfo(info)
}

// WithdrawFunds - Funds withdrawal
//@:public:method:gas[500]
func (mr *MyRoulette) WithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) {

	sdk.RequireOwner(mr.sdk)
	sdk.Require(withdrawAmount.CmpI(0) > 0,
		types.ErrInvalidParameter, "withdrawAmount must be larger than zero")

	account := mr.sdk.Helper().AccountHelper().AccountOf(mr.sdk.Message().Contract().Account())
	lockedAmount := mr._lockedInBets(tokenName)
	unlockedAmount := account.BalanceOfName(tokenName).Sub(lockedAmount)
	sdk.Require(unlockedAmount.Cmp(withdrawAmount) >= 0,
		types.ErrInvalidParameter, "Not enough funds")

	// transfer to beneficiary
	account.TransferByName(tokenName, beneficiary, withdrawAmount)

	// fire event
	mr.emitWithdrawFunds(tokenName, beneficiary, withdrawAmount)
}

// PlaceBet - place bet
//@:public:method:gas[500]
func (mr *MyRoulette) PlaceBet(betInfoJson string, commitLastBlock int64, commit, signData []byte, refAddress types.Address) {

	sdk.Require(mr.sdk.Message().Sender().Address() != mr.sdk.Message().Contract().Owner(),
		types.ErrNoAuthorization, "The contract owner cannot bet")

	//1. Verify whether the signature and current round betting are legal
	data := append(bn.N(commitLastBlock).Bytes(), commit...)
	sdk.Require(ed25519.VerifySign(mr._publicKey(), data, signData),
		types.ErrInvalidParameter, "Incorrect signature")

	gambler := mr.sdk.Message().Sender().Address()
	hexCommit := hex.EncodeToString(commit)
	//Is late
	sdk.Require(mr.sdk.Block().Height() <= commitLastBlock,
		types.ErrInvalidParameter, "Commit has expired")

	//Verify that the user has made a note
	sdk.Require(!mr._chkBetInfo(hexCommit, gambler),
		types.ErrInvalidParameter, "Commit should be new")

	setting := mr._setting()
	var roundInfo *RoundInfo
	if !mr._chkRoundInfo(hexCommit) {
		roundInfo = &RoundInfo{
			Commit:           commit,
			State:            NOAWARD,
			FirstBlockHeight: mr.sdk.Block().Height(),
			Setting:          setting,
			TotalBuyAmount:   mr.CreateMapByTokenName(setting.TokenNames),
		}
		mr._setRoundInfo(hexCommit, roundInfo)
	}
	roundInfo = mr._roundInfo(hexCommit)
	//Whether the current wheel state allows betting
	sdk.Require(NOAWARD == roundInfo.State,
		types.ErrInvalidParameter, "No betting on the current wheel")
	sdk.Require(roundInfo.FirstBlockHeight+roundInfo.Setting.BetExpirationBlocks > mr.sdk.Block().Height(),
		types.ErrInvalidParameter, "This round is time out")

	tokenName := ""
	transferAmount := bn.N(0)
	mr.GetTransferData(setting, &tokenName, &transferAmount)

	//Verify receipt of bet transfer
	sdk.Require(tokenName != "" && transferAmount.CmpI(0) > 0,
		types.ErrUserDefined, "Must transfer tokens to me before place a bet")

	betDataList := make([]BetData, 0)
	err := jsoniter.Unmarshal([]byte(betInfoJson), &betDataList)
	sdk.RequireNotError(err, types.ErrInvalidParameter)

	totalAmount := bn.N(0)
	for _, betData := range betDataList {
		amount := betData.BetAmount

		//Verify that the betting scheme is legal
		sdk.Require(betData.BetMode >= REDTYPE && betData.BetMode <= ELEVENSIXDIGIT,
			types.ErrInvalidParameter, "The betting range is illegal")
		//Verify that the bet mode and bet value is legal
		sdk.Require(mr.checkBetData(betData),
			types.ErrInvalidParameter, "The betting model and value not matching")
		//Verify that the bet amount is legal
		sdk.Require(amount.CmpI(setting.MinLimit) >= 0 && amount.CmpI(setting.MaxLimit) <= 0,
			types.ErrInvalidParameter, "Amount should be within range")

		totalAmount = totalAmount.Add(amount)
	}

	//Check whether the total amount of bet is equal to the transfer amount
	sdk.Require(totalAmount.Cmp(transferAmount) == 0,
		types.ErrUserDefined, "transfer amount not equal place bet amount")

	//Verify that the reimbursable amount is sufficient
	//Calculate the amount of money you can win
	totalMaybeWinAmount, feeAmount := mr.GetMayWinData(setting, totalAmount, betDataList)
	//Is the amount likely to be won less than or equal to the maximum bonus amount
	sdk.Require(totalMaybeWinAmount.CmpI(setting.MaxProfit) <= 0,
		types.ErrInvalidParameter, "MaxProfit limit violation")

	contractAcct := mr.sdk.Helper().AccountHelper().AccountOf(mr.sdk.Message().Contract().Account())
	//Lock in the amount that may need to be paid
	totalLockedAmount := mr._lockedInBets(tokenName).Add(totalMaybeWinAmount).Add(feeAmount)
	//Contract account balance
	totalUnlockedAmount := contractAcct.BalanceOfName(tokenName)
	//Is the contract account balance greater than or equal to the balance that may need to be paid
	sdk.Require(totalUnlockedAmount.Cmp(totalMaybeWinAmount) >= 0,
		types.ErrInvalidParameter, "Cannot afford to lose this bet")
	mr._setLockedInBets(tokenName, totalLockedAmount)

	//Store bet information
	betInfo := &BetInfo{}
	betInfo.TokenName = tokenName
	betInfo.Amount = totalAmount
	betInfo.BetData = betDataList
	betInfo.WinAmount = totalMaybeWinAmount
	betInfo.Settled = false
	betInfo.Gambler = gambler

	mr._setBetInfo(hexCommit, gambler, betInfo)

	roundInfo = mr._roundInfo(hexCommit)
	//Round information add betting information k2
	roundInfo.BetInfoSerialNumber = append(roundInfo.BetInfoSerialNumber, gambler)
	roundInfo.TotalBuyAmount[tokenName] = roundInfo.TotalBuyAmount[tokenName].Add(totalAmount)
	roundInfo.TotalBetCount += 1
	//Save the current wheel information
	mr._setRoundInfo(hexCommit, roundInfo)

	mr.emitPlaceBet(tokenName, gambler, totalMaybeWinAmount, betDataList, commitLastBlock, commit, signData, refAddress)
}

// SettleBet - The lottery and settlement
//@:public:method:gas[500]
func (mr *MyRoulette) SettleBet(reveal []byte, settleCount int64) {

	sdk.Require(len(reveal) > 0,
		types.ErrInvalidParameter, "Commit does not exist")

	sdk.RequireOwner(mr.sdk)
	hexCommit := hex.EncodeToString(sha3.Sum256(reveal))
	sdk.Require(mr._chkRoundInfo(hexCommit),
		types.ErrInvalidParameter, "Commit should be not exist")
	roundInfo := mr._roundInfo(hexCommit)

	//Current wheel configuration
	settings := roundInfo.Setting
	//The bet height of the round to be settled should be less than the settlement height
	sdk.Require(roundInfo.FirstBlockHeight < mr.sdk.Block().Height(),
		types.ErrInvalidParameter, "SettleBet block can not be in the same block as placeBet, or before.")

	sdk.Require(NOAWARD == roundInfo.State || OPENINGAPRIZE == roundInfo.State,
		types.ErrInvalidParameter, "This state does not operate for settlement")

	//For the first settlement, the round information cannot expire
	if NOAWARD == roundInfo.State {
		sdk.Require(roundInfo.FirstBlockHeight+settings.BetExpirationBlocks > mr.sdk.Block().Height(),
			types.ErrInvalidParameter, "This round is time out")
		roundInfo.WinningResult = mr.OpenNumber()
		//No one bets on the current round
		if roundInfo.TotalBetCount <= 0 {
			roundInfo.State = AWARDED
			mr._setRoundInfo(hexCommit, roundInfo)
			return
		}
		roundInfo.State = OPENINGAPRIZE
		mr._setRoundInfo(hexCommit, roundInfo)
	}
	//Determine whether all bets have been settled
	if roundInfo.TotalBetCount == roundInfo.ProcessCount {
		roundInfo.State = AWARDED
		mr._setRoundInfo(hexCommit, roundInfo)
		sdk.Require(false,
			types.ErrInvalidParameter, "This round is complete")
	}

	tokenNameList := settings.TokenNames
	//Money that could be won
	totalPossibleWinAmount := mr.CreateMapByTokenName(tokenNameList)
	//The actual winning money
	totalWinAmount := mr.CreateMapByTokenName(tokenNameList)
	//Total handling charge
	totalFeeAmount := mr.CreateMapByTokenName(tokenNameList)
	//Key of betting information
	betInfoSerialNumberList := roundInfo.BetInfoSerialNumber
	//The lottery information
	winningResult := roundInfo.WinningResult
	//Contract account
	contractAcct := mr.sdk.Helper().AccountHelper().AccountOf(mr.sdk.Message().Contract().Account())
	var winCount int64 = 0

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
	for i := startIndex; i < endIndex; i++ {
		betInfoKey := betInfoSerialNumberList[i]
		//for _, betInfoKey := range betInfoSerialNumberList {
		betInfo := mr._betInfo(hexCommit, betInfoKey)
		if betInfo.Settled {
			continue
		}
		//The currency name of the bet such as BCB
		tokenName := betInfo.TokenName

		//resultType int64, mr *MyRoulette, resultList []int64
		betDataList := betInfo.BetData
		winAmount := winningResult.GetBetWinAmount(mr, betDataList)

		//If the winning amount is greater than 0, transfer the prize money to the player's address
		if winAmount.CmpI(0) > 0 {
			totalWinAmount[tokenName] = totalWinAmount[tokenName].Add(winAmount)
			contractAcct.TransferByName(tokenName, betInfo.Gambler, winAmount)
			winCount++
		}

		//Update the current betting status database to be settled
		betInfo.Settled = true
		mr._setBetInfo(hexCommit, betInfoKey, betInfo)
		totalMaybeWinAmount, feeAmount := mr.GetMayWinData(settings, betInfo.Amount, betDataList)
		totalPossibleWinAmount[tokenName] = totalPossibleWinAmount[tokenName].Add(totalMaybeWinAmount).Add(feeAmount)

		//fee
		totalFeeAmount[tokenName] = totalFeeAmount[tokenName].Add(feeAmount)
	}

	//Unlock lock amount
	for _, tokenName := range tokenNameList {
		lockedInBet := mr._lockedInBets(tokenName)
		mr._setLockedInBets(tokenName, lockedInBet.Sub(totalPossibleWinAmount[tokenName]))
	}

	//participation in profit
	if settings.SendToCltRatio > 0 {
		for _, tokenName := range tokenNameList {
			amount := totalFeeAmount[tokenName].MulI(roundInfo.Setting.SendToCltRatio).DivI(PERMILLE)
			contractAcct.TransferByName(tokenName, mr.sdk.Helper().BlockChainHelper().CalcAccountFromName("clt",""), amount)
			totalFeeAmount[tokenName] = totalFeeAmount[tokenName].Sub(amount)
		}
	}

	//Transfer to other handling address
	mr.transferToRecFeeAddr(tokenNameList, totalFeeAmount)
	roundInfo.ProcessCount = endIndex
	mr._setRoundInfo(hexCommit, roundInfo)
	//Send the receipt
	mr.emitSettleBet(tokenNameList, reveal, roundInfo.Commit, roundInfo.BetInfoSerialNumber, roundInfo.WinningResult.Value, totalWinAmount, roundInfo.State == AWARDED)
}

// RefundBets - Refund will be made if the prize is not paid after the time limit
//@:public:method:gas[500]
func (mr *MyRoulette) RefundBets(commit []byte, refundCount int64) {

	sdk.Require(len(commit) > 0, types.ErrInvalidParameter, "Commit should be not exist")

	sdk.RequireOwner(mr.sdk)

	hexCommit := hex.EncodeToString(commit)

	sdk.Require(mr._chkRoundInfo(hexCommit),
		types.ErrInvalidParameter, "Commit should be not exist")

	//Determine whether the bet can be refunded
	roundInfo := mr._roundInfo(hexCommit)
	//Current wheel configuration
	settings := roundInfo.Setting
	//The bet height of the round to be settled should be less than the settlement height
	sdk.Require(mr.sdk.Block().Height() > roundInfo.FirstBlockHeight+settings.BetExpirationBlocks,
		types.ErrInvalidParameter, "SettleBet block can not be in the same block as placeBet, or before.")
	//Whether the current round status can be refunded
	sdk.Require(NOAWARD == roundInfo.State,
		types.ErrInvalidParameter, "This status does not operate for a refund")
	//Whether the number of bets processed is less than the total number of bets
	sdk.Require(roundInfo.TotalBetCount > roundInfo.ProcessCount,
		types.ErrInvalidParameter, "There are currently no refundable bets")

	betInfoSerialNumberList := roundInfo.BetInfoSerialNumber
	//Whether betting information exists
	sdk.Require(len(betInfoSerialNumberList) > 0,
		types.ErrInvalidParameter, "There are currently no refundable bets")

	//Whether the dice result has been drawn
	sdk.Require(roundInfo.WinningResult == nil,
		types.ErrInvalidParameter, "The current round has been lottery, can not operate refund")

	//Contract account
	contractAcct := mr.sdk.Helper().AccountHelper().AccountOf(mr.sdk.Message().Contract().Account())

	tokenNameList := settings.TokenNames

	//Money that could be won
	totalPossibleWinAmount := mr.CreateMapByTokenName(tokenNameList)
	refundedAmount := mr.CreateMapByTokenName(tokenNameList)

	//Initial index
	startIndex := roundInfo.ProcessCount
	if startIndex < 0 {
		startIndex = 0
	}
	endIndex := startIndex + refundCount

	if endIndex >= roundInfo.TotalBetCount {
		endIndex = roundInfo.TotalBetCount
		roundInfo.State = REFUNDED
	}
	for i := startIndex; i < endIndex; i++ {
		betInfoKey := betInfoSerialNumberList[i]

		betInfo := mr._betInfo(hexCommit, betInfoKey)
		if betInfo.Settled {
			continue
		}
		//The currency name of the bet such as BCB
		tokenName := betInfo.TokenName
		totalAmount := betInfo.Amount
		//If the bet amount is greater than 0, the bet will be transferred to the player's address
		if betInfo.Amount.CmpI(0) > 0 {
			contractAcct.TransferByName(tokenName, betInfo.Gambler, totalAmount)
			refundedAmount[tokenName] = refundedAmount[tokenName].Add(totalAmount)
		}
		//Update the current betting status database to be settled
		betInfo.Settled = true
		mr._setBetInfo(hexCommit, betInfoKey, betInfo)

		//Calculate the amount of money you can win
		totalMaybeWinAmount, feeAmount := mr.GetMayWinData(settings, totalAmount, betInfo.BetData)
		totalPossibleWinAmount[tokenName] = totalPossibleWinAmount[tokenName].Sub(totalMaybeWinAmount).Sub(feeAmount)

	}

	//Unlock lock amount
	for _, tokenName := range tokenNameList {
		lockedInBet := mr._lockedInBets(tokenName)
		mr._setLockedInBets(tokenName, lockedInBet.Sub(totalPossibleWinAmount[tokenName]))
	}

	//Set the database state state to refunded, and ProcessCount updates it
	roundInfo.ProcessCount = endIndex
	mr._setRoundInfo(hexCommit, roundInfo)

	//Send the receipt
	mr.emitRefundBet(commit, tokenNameList, betInfoSerialNumberList, refundedAmount, roundInfo.State == REFUNDED)
}
