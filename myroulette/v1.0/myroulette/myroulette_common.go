package myroulette

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

//check bet data betmodel and betvalue is matching
func (mr *MyRoulette) checkBetData(betdata BetData) bool {
	//是否在合理的点位
	for _, v := range betdata.BetValue {
		if !(v >= DICENUMBERMIN && v <= DICENUMBERMAX) {
			return false
		}
	}
	//是否value的元素个数是否
	flag := true
	switch betdata.BetMode {
	//这都是有对应的模型value 不用传具体的元素 所以元素个数为0
	case REDTYPE, BLACKTYPE, SINGLETYPE, DOUBLETYPE, BIGTYPE, SMALLTYPE, FIRSETAREA, SECONDAREA, THIRDAREA, FIRSTINLINE, SECONDINLINE, THIRDINLINE, FIVEDIGITCOMBINATION:
		if len(betdata.BetValue) != 0 {
			flag = false
		}
	case ASINGLENUMBER:
		if len(betdata.BetValue) != 1 {
			flag = false
		}
	case TWODIGITCOMBINATION:
		if (len(betdata.BetValue) != 2) || !mr.BetTwoCheck(betdata.BetValue) {
			flag = false
		}
	case THREEDIGITCOMBINATION:
		if (len(betdata.BetValue) != 3) || !mr.BetThreeCheck(betdata.BetValue) {
			flag = false
		}
	case FOURDIGITCOMBINATION:
		if (len(betdata.BetValue) != 4) || !mr.BetFourCheck(betdata.BetValue) {
			flag = false
		}

	}
	return flag
}

//Verify that the two points conform to the specification
func (mr *MyRoulette) BetTwoCheck(resultsArray []int64) bool {
	oneValue := resultsArray[0]
	twoValue := resultsArray[1]

	if (oneValue+1) == twoValue || oneValue+3 == twoValue {
		return true
	}

	return false
}

//Verify that the three points conform to the specification
func (mr *MyRoulette) BetThreeCheck(resultsArray []int64) bool {
	oneValue := resultsArray[0]
	twoValue := resultsArray[1]
	threeValue := resultsArray[2]
	//(0,37,2)match===>(0,00,2)
	//verify is not exceptional case
	if twoValue == 37 {
		if oneValue+twoValue+threeValue == 39 {
			return true
		} else {
			return false
		}
	}
	if (oneValue + 1) != twoValue {
		return false
	}
	if oneValue+2 != threeValue {
		return false
	}

	return true
}

//Verify that the four points conform to the specification
func (mr *MyRoulette) BetFourCheck(resultsArray []int64) bool {
	oneValue := resultsArray[0]
	twoValue := resultsArray[1]
	threeValue := resultsArray[2]
	fourValue := resultsArray[3]
	if (oneValue + 1) != twoValue {
		return false
	}
	if oneValue+3 != threeValue {
		return false
	}
	if twoValue+3 != fourValue {
		return false
	}
	return true
}

//Initializes the map according to tokenNameList
func (mr *MyRoulette) CreateMapByTokenName(tokenNameList []string) (maps map[string]bn.Number) {
	maps = make(map[string]bn.Number, len(tokenNameList))
	for _, value := range tokenNameList {
		maps[value] = bn.N(0)
	}
	return
}

func (mr *MyRoulette) GetTransferData(settings *Setting, tokenName *string, transferAmount *bn.Number) {
	for _, name := range settings.TokenNames {
		transferReceipt := mr.sdk.Message().GetTransferToMe(name)
		if transferReceipt != nil {
			*tokenName = name
			*transferAmount = transferReceipt.Value
			break
		}
	}

}

func (mr *MyRoulette) checkSettings(newSettings *Setting) {

	sdk.Require(len(newSettings.TokenNames) > 0,
		types.ErrInvalidParameter, "tokenNames cannot be empty")

	for _, tokenName := range newSettings.TokenNames {
		token := mr.sdk.Helper().TokenHelper().TokenOfName(tokenName)
		sdk.Require(token != nil,
			types.ErrInvalidParameter, fmt.Sprintf("tokenName=%s is not exist", tokenName))
	}

	sdk.Require(newSettings.MaxLimit > 0,
		types.ErrInvalidParameter, "MaxBet must be bigger than zero")

	sdk.Require(newSettings.MaxProfit >= 0,
		types.ErrInvalidParameter, "MaxProfit can not be negative")

	sdk.Require(newSettings.MinLimit > 0 && newSettings.MinLimit < newSettings.MaxLimit,
		types.ErrInvalidParameter, "MinBet must be bigger than zero and smaller than MaxBet")

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
}

func (mr *MyRoulette) checkRecFeeInfo(infos []RecFeeInfo) {
	sdk.Require(len(infos) > 0,
		types.ErrInvalidParameter, "The length of RecvFeeInfos must be larger than zero")

	allRatio := int64(0)
	for _, info := range infos {
		sdk.Require(info.RecFeeRatio > 0,
			types.ErrInvalidParameter, "ratio must be larger than zero")
		sdk.RequireAddress(mr.sdk, info.RecFeeAddr)
		sdk.Require(info.RecFeeAddr != mr.sdk.Message().Contract().Account(),
			types.ErrInvalidParameter, "address cannot be contract account address")

		allRatio += info.RecFeeRatio
	}

	//The allocation ratio set must add up to 1000
	sdk.Require(allRatio <= 1000,
		types.ErrInvalidParameter, "The sum of ratio must be less or equal 1000")
}

//Transfer to fee's receiving address
func (mr *MyRoulette) transferToRecFeeAddr(tokenNameList []types.Address, recFeeMap map[string]bn.Number) {

	account := mr.sdk.Helper().AccountHelper().AccountOf(mr.sdk.Message().Contract().Account())
	for _, tokenName := range tokenNameList {
		recFee := recFeeMap[tokenName]
		if recFee.CmpI(0) <= 0 {
			continue
		}
		infos := mr._recFeeInfo()
		for _, info := range infos {
			account.TransferByName(tokenName, info.RecFeeAddr, recFee.MulI(info.RecFeeRatio).DivI(PERMILLE))
		}
	}

}
