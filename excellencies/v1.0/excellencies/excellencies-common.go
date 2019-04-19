package excellencies

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

func (sg *Excellencies) checkSettings(newSettings *Settings) {

	sdk.Require(len(newSettings.TokenNames) > 0,
		types.ErrInvalidParameter, "tokenNames cannot be empty")

	for _, tokenName := range newSettings.TokenNames {
		token := sg.sdk.Helper().TokenHelper().TokenOfName(tokenName)
		sdk.Require(token != nil,
			types.ErrInvalidParameter, fmt.Sprintf("tokenName=%s is not exist", tokenName))

		sdk.Require(newSettings.LimitMaps[tokenName].MaxLimit > 0,
			types.ErrInvalidParameter, "MaxBet must be bigger than zero")

		sdk.Require(newSettings.LimitMaps[tokenName].MaxProfit >= 0,
			types.ErrInvalidParameter, "MaxProfit can not be negative")

		sdk.Require(newSettings.LimitMaps[tokenName].MinLimit > 0 && newSettings.LimitMaps[tokenName].MinLimit < newSettings.LimitMaps[tokenName].MaxLimit,
			types.ErrInvalidParameter, "MinBet must be bigger than zero and smaller than MaxBet")
	}

	sdk.Require(newSettings.SendToCltRatio >= 0 && newSettings.SendToCltRatio < PERMILLE,
		types.ErrInvalidParameter,
		fmt.Sprintf("SendToCltRatio must be bigger than zero and smaller than %d", PERMILLE))

	sdk.Require(newSettings.FeeRatio > 0 && newSettings.FeeRatio < PERMILLE,
		types.ErrInvalidParameter,
		fmt.Sprintf("FeeRatio must be bigger than zero and  smaller than %d", PERMILLE))

	sdk.Require(newSettings.FeeMiniNum > 0,
		types.ErrInvalidParameter, "FeeMinimum must be bigger than zero")

	sdk.Require(newSettings.BetExpirationBlocks > 0,
		types.ErrInvalidParameter, "BetExpirationBlocks must be bigger than zero")
	sdk.Require(newSettings.CarveUpPoolRatio > 0 && newSettings.CarveUpPoolRatio < PERMILLE,
		types.ErrInvalidParameter,
		fmt.Sprintf("CarveUpPoolRatio must be bigger than zero and  smaller than %d", PERMILLE))
}

func (sg *Excellencies) checkRecFeeInfo(infos []RecFeeInfo) {
	sdk.Require(len(infos) > 0,
		types.ErrInvalidParameter, "The length of RecvFeeInfos must be larger than zero")

	allRatio := int64(0)
	for _, info := range infos {
		sdk.Require(info.RecFeeRatio > 0,
			types.ErrInvalidParameter, "ratio must be larger than zero")
		sdk.RequireAddress(sg.sdk, info.RecFeeAddr)
		sdk.Require(info.RecFeeAddr != sg.sdk.Message().Contract().Account(),
			types.ErrInvalidParameter, "address cannot be contract account address")

		allRatio += info.RecFeeRatio
	}

	//The allocation ratio set must add up to 1000
	sdk.Require(allRatio <= 1000,
		types.ErrInvalidParameter, "The sum of ratio must be less or equal 1000")
}

//Settlement public method
func (e *Excellencies) DealSettle(roundInfo *RoundInfo, hexCommit string, startIndex int64, endIndex int64) {

	settings := roundInfo.Settings
	tokenNameList := settings.TokenNames
	//The actual winning money
	totalWinAmount := e.CreateMapByTokenName(tokenNameList)
	//Total handling charge
	totalFeeAmount := e.CreateMapByTokenName(tokenNameList)
	//Money that could be won
	totalPossibleWinAmount := e.CreateMapByTokenName(tokenNameList)

	//奖池金额
	totalPoolAmount := e.CreateMapByTokenName(tokenNameList)

	//Contract account
	contractAcct := e.sdk.Helper().AccountHelper().AccountOf(e.sdk.Message().Contract().Account())

	//game result
	game := roundInfo.Game
	var winCount int64 = 0

	for i := startIndex; i < endIndex; i++ {

		if roundInfo.BetInfos[i].Settled {
			continue
		}
		//The currency name of the bet such as BCB
		tokenName := roundInfo.BetInfos[i].TokenName
		//The actual winning amount, the actual winning amount
		winAmount, winFeeAmount := game.GetBetWinAmount(roundInfo.BetInfos[i].BetData)

		totalMaybeWinAmount, feeAmount := e.MaybeWinAmountAndFeeByList(roundInfo.BetInfos[i].Amount, roundInfo.BetInfos[i].BetData)
		//If the winning amount is greater than 0, transfer the prize money to the player's address
		if winAmount.CmpI(0) > 0 {

			//减去中奖的总投注金额对应的手续费
			winAmount = winAmount.Sub(winFeeAmount.MulI(settings.FeeRatio).DivI(PERMILLE))
			//奖池费用
			poolFee := winAmount.MulI(settings.PoolFeeRatio).DivI(PERMILLE)
			//中奖金额减去奖池金额
			actualWinAmount := winAmount.Sub(poolFee)

			totalPoolAmount[tokenName] = totalPoolAmount[tokenName].Add(poolFee)

			totalWinAmount[tokenName] = totalWinAmount[tokenName].Add(winAmount)

			contractAcct.TransferByName(tokenName, roundInfo.BetInfos[i].Gambler, actualWinAmount)

			roundInfo.BetInfos[i].WinAmount = actualWinAmount
			winCount++

		}
		//Update the current betting status database to be settled
		roundInfo.BetInfos[i].Settled = true
		//Calculate the amount of money you can win
		totalPossibleWinAmount[tokenName] = totalPossibleWinAmount[tokenName].Add(totalMaybeWinAmount)
		//fee
		totalFeeAmount[tokenName] = totalFeeAmount[tokenName].Add(feeAmount)

	}

	poolTokens:=make(map[string]bn.Number)
	//Unlock lock amount
	for _, tokenName := range tokenNameList {
		lockedInBet := e._lockedInBets(tokenName)
		e._setLockedInBets(tokenName, lockedInBet.Sub(totalPossibleWinAmount[tokenName]))
		amount:=e._poolAmount(tokenName).Add(totalPoolAmount[tokenName])
		//奖池
		if totalPoolAmount[tokenName].CmpI(0) > 0 {
			e._setPoolAmount(tokenName, amount)
		}
		poolTokens[tokenName]=amount
	}

	//participation in profit
	if settings.SendToCltRatio > 0 {
		for _, tokenName := range tokenNameList {
			amount := totalFeeAmount[tokenName].MulI(roundInfo.Settings.SendToCltRatio).DivI(PERMILLE)
			contractAcct.TransferByName(tokenName, e.sdk.Helper().BlockChainHelper().CalcAccountFromName("clt",""), amount)
			totalFeeAmount[tokenName] = totalFeeAmount[tokenName].Sub(amount)
		}
	}

	//Transfer to other handling address
	e.transferToRecFeeAddr(tokenNameList, totalFeeAmount)

	roundInfo.ProcessCount = endIndex
	e._setRoundInfo(hexCommit, roundInfo)

	//Send the receipt
	e.emitSettleBet(tokenNameList, roundInfo.Commit, roundInfo.Game, winCount, totalWinAmount, roundInfo.State == AWARDED,poolTokens,roundInfo.OriginPokers)

}

//Initializes the map according to tokenNameList
func (e *Excellencies) CreateMapByTokenName(tokenNameList []string) (maps map[string]bn.Number) {
	maps = make(map[string]bn.Number, len(tokenNameList))
	for _, value := range tokenNameList {
		maps[value] = bn.N(0)
	}
	return
}

//Transfer to fee's receiving address
func (e *Excellencies) transferToRecFeeAddr(tokenNameList []types.Address, recFeeMap map[string]bn.Number) {

	account := e.sdk.Helper().AccountHelper().AccountOf(e.sdk.Message().Contract().Account())
	for _, tokenName := range tokenNameList {
		recFee := recFeeMap[tokenName]
		if recFee.CmpI(0) <= 0 {
			continue
		}
		infos := e._recFeeInfo()
		for _, info := range infos {
			account.TransferByName(tokenName, info.RecFeeAddr, recFee.MulI(info.RecFeeRatio).DivI(PERMILLE))
		}
	}
}


//Transfer to players
func (e *Excellencies) TransferToPlayers(tokenName string, gamblerAddr string, accountMoney bn.Number) {
	account := e.sdk.Helper().AccountHelper().AccountOf(e.sdk.Message().Contract().Account())
	account.TransferByName(tokenName, gamblerAddr, accountMoney)
}

