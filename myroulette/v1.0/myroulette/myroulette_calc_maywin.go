package myroulette

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"fmt"
)

//Calculate the amount of money a long bet is likely to win: total bet amount betList scheme
func (mr *MyRoulette) GetMayWinData(settings *Setting, amount bn.Number, betList []BetData) (totalMaybeWinAmount, fee bn.Number) {
	fee = bn.N(0)
	totalMaybeWinAmount = bn.N(0)
	//settings := sb._settings()
	// fee
	fee = amount.MulI(settings.FeeRatio).DivI(PERMILLE)
	if fee.CmpI(settings.FeeMiniNum) < 0 {
		fee = bn.N(settings.FeeMiniNum)
	}

	sdk.Require(fee.Cmp(amount) <= 0,
		types.ErrInvalidParameter, "Bet doesn't even cover fee")

	for _, bet := range betList {

		totalMaybeWinAmount = totalMaybeWinAmount.Add(MaybeWinAmountByOne(settings, bet))
		fmt.Println(totalMaybeWinAmount)
	}
	return
}

//Calculate the amount of money a bet is likely to win
func MaybeWinAmountByOne(settings *Setting, bet BetData) (amount bn.Number) {

	var oddsValue int64

	switch bet.BetMode {
	//red
	case REDTYPE:
		oddsValue = ONETIMES
		//black
	case BLACKTYPE:
		oddsValue = ONETIMES
		//single double
	case SINGLETYPE, DOUBLETYPE:
		oddsValue = ONETIMES
	case BIGTYPE, SMALLTYPE:
		oddsValue = ONETIMES
	case FIRSETAREA, SECONDAREA, THIRDAREA:
		oddsValue = TWOTIMES
	case FIRSTINLINE, SECONDINLINE, THIRDINLINE:
		oddsValue = TWOTIMES
	case ASINGLENUMBER:
		oddsValue = THIRTYTIMES
	case TWODIGITCOMBINATION:
		oddsValue = SEVENTEENTIMES
	case THREEDIGITCOMBINATION:
		oddsValue = ELEVENTIMES
	case FOURDIGITCOMBINATION:
		oddsValue = EIGHTTIMES
	case FIVEDIGITCOMBINATION:
		oddsValue = SIXTIMES
	case FIRSTSIXDIGIT, SECONDSIXDIGIT, THIRDSIXDIGIT, FOURSIXDIGIT, FIVESIXDIGIT, SIXSIXDIGIT,
		SEVENSIXDIGIT, EIGHTSIXDIGIT, NINESIXDIGIT, TENSIXDIGIT, ELEVENSIXDIGIT:
		oddsValue = SIXTIMES
	}

	if oddsValue <= 0 {
		amount = bn.N(0)
	} else {
		amount = bet.BetAmount.MulI(oddsValue + (1 - settings.FeeRatio))
	}

	return
}
