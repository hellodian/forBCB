package myroulette

import (
	"blockchain/smcsdk/sdk/bn"
)

type WinningResult struct {
	Value int64 `json:"value"`
}

//Get the lottery random number
func (mr *MyRoulette) OpenNumber() *WinningResult {
	bytes := bn.NBytes(mr.sdk.Block().RandomNumber())
	size := int64(RANDOMMODYLO)
	ran := bytes.ModI(size)
	value := ran.V.Int64()
	winningResult := WinningResult{}
	winningResult.Value = value
	return &winningResult
}

//Determine whether to win
func (wr *WinningResult) IsWinning(resultType int64, mr *MyRoulette, resultList []int64) (winningStatus bool, oddsValue int64) {
	winningStatus = false
	resultValue := wr.Value
	//Eat all
	if resultValue == 0 || resultValue == 37 {
		winningStatus = wr.EatAll(resultType)
		if winningStatus == false {
			oddsValue = 0
			return
		}
	}

	switch resultType {
	//red
	case REDTYPE:
		RED := []int64{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36}
		winningStatus = wr.CheckSlice(RED)
		oddsValue = ONETIMES
		//black
	case BLACKTYPE:
		BLACK := []int64{2, 4, 6, 8, 10, 11, 13, 15, 17, 20, 22, 24, 26, 28, 29, 31, 33, 35}
		winningStatus = wr.CheckSlice(BLACK)
		oddsValue = ONETIMES
		//single double
	case SINGLETYPE, DOUBLETYPE:
		winningStatus = wr.CheckSingleDouble(resultType)
		oddsValue = ONETIMES
	case BIGTYPE, SMALLTYPE:
		winningStatus = wr.CheckSize(resultType)
		oddsValue = ONETIMES
	case FIRSETAREA, SECONDAREA, THIRDAREA:
		winningStatus = wr.CheckArea(resultType)
		oddsValue = TWOTIMES
	case FIRSTINLINE, SECONDINLINE, THIRDINLINE:
		winningStatus = wr.CheckInline(resultType)
		oddsValue = TWOTIMES
	case ASINGLENUMBER:
		winningStatus = wr.CheckSinglePoint(resultList)
		oddsValue = THIRTYTIMES
	case TWODIGITCOMBINATION:
		winningStatus = wr.TwoDigitCombination(resultList, mr)
		oddsValue = SEVENTEENTIMES
	case THREEDIGITCOMBINATION:
		winningStatus = wr.ThreeDigitCombination(resultList, mr)
		oddsValue = ELEVENTIMES
	case FOURDIGITCOMBINATION:
		winningStatus = wr.FourDigitCombination(resultList, mr)
		oddsValue = EIGHTTIMES
	case FIVEDIGITCOMBINATION:
		five := []int64{0, 1, 2, 3, 37}
		winningStatus = wr.CheckSlice(five)
		oddsValue = SIXTIMES
	case FIRSTSIXDIGIT, SECONDSIXDIGIT, THIRDSIXDIGIT, FOURSIXDIGIT, FIVESIXDIGIT, SIXSIXDIGIT,
		SEVENSIXDIGIT, EIGHTSIXDIGIT, NINESIXDIGIT, TENSIXDIGIT, ELEVENSIXDIGIT:
		winningStatus = wr.SixDigitCombination(resultType)
		oddsValue = SIXTIMES
	}
	return
}

//Calculate multiple entry bonus
func (wr *WinningResult) GetBetWinAmount(mr *MyRoulette, betList []BetData) (totalWinAmount bn.Number) {
	totalWinAmount = bn.N(0)
	if len(betList) <= 0 {
		return
	}
	for _, bet := range betList {
		//Calculate note bonus
		betMode 		:= bet.BetMode
		resultList 		:= bet.BetValue
		winningStatus, oddsValue := wr.IsWinning(betMode, mr, resultList)
		if winningStatus == false {
			return
		}
		amount := bet.BetAmount.MulI(oddsValue)
		//Bonus greater than 0
		if bn.N(0).Cmp(amount) < 0 {
			totalWinAmount = totalWinAmount.Add(amount)
		}
	}
	return
}

//Determine if the winning result is within the slice
func (wr *WinningResult) CheckSlice(resultList []int64) bool {
	for _, value := range resultList {
		if wr.Value == value {
			return true
		}
	}
	return false
}

//Check single and double
func (wr *WinningResult) CheckSingleDouble(resultType int64) bool {
	remaining := wr.Value % 2
	if resultType == SINGLETYPE && remaining != 0 {
		return true
	} else if resultType == DOUBLETYPE && remaining == 0 {
		return true
	}
	return false
}

//Check size
func (wr *WinningResult) CheckSize(resultType int64) bool {
	remaining := wr.Value
	if remaining >= 1 && remaining <= 18 && resultType == SMALLTYPE {
		return true
	} else if remaining >= 19 && remaining <= 36 && resultType == BIGTYPE {
		return true
	}
	return false
}

//Check area
func (wr *WinningResult) CheckArea(resultType int64) bool {
	remaining := wr.Value
	if resultType == FIRSETAREA && remaining >= 1 && remaining <= 12 {
		return true
	} else if resultType == SECONDAREA && remaining >= 13 && remaining <= 24 {
		return true
	} else if resultType == THIRDAREA && remaining >= 25 && remaining <= 36 {
		return true
	}
	return false
}

//check inline
func (wr *WinningResult) CheckInline(resultType int64) bool {
	remaining := wr.Value % 3
	if resultType == FIRSTINLINE && remaining == 1 {
		return true
	} else if resultType == SECONDINLINE && remaining == 2 {
		return true
	} else if resultType == THIRDINLINE && remaining == 0 {
		return true
	}
	return false
}

//Check single point
func (wr *WinningResult) CheckSinglePoint(resultList []int64) bool {
	//00 == 37
	resultValue := wr.Value
	for _, value := range resultList {
		if resultValue == value && value >= 0 && value <= 37 {
			return true
		}
	}
	return false
}

//check two digit combination
func (wr *WinningResult) TwoDigitCombination(resultList []int64, mr *MyRoulette) bool {
	betTwoStatus := mr.BetTwoCheck(resultList)
	if betTwoStatus == true && wr.CheckSlice(resultList) == true {
		return true
	}
	return false
}

//check three digit combination
func (wr *WinningResult) ThreeDigitCombination(resultList []int64, mr *MyRoulette) bool {
	specialList := []int64{0, 2, 37}
	//Determine if the slices are equal
	inNumberList := 0
	for _, specialValue := range specialList {
		for _, getValue := range resultList {
			if getValue == specialValue {
				inNumberList += 1
			}
		}
	}

	if inNumberList ==3 {
		resultValue := wr.Value
		winningNumber := 0
		for _, value := range specialList {
			if resultValue == value {
				winningNumber += 1
			}
		}
		if winningNumber == 3 {
			return true
		}
	} else {
		betThreeStatus := mr.BetThreeCheck(resultList)
		if betThreeStatus == true && wr.CheckSlice(resultList) == true {
			return true
		}
	}
	return false
}

//check three digit combination
func (wr *WinningResult) FourDigitCombination(resultList []int64, mr *MyRoulette) bool {
	betFourStatus := mr.BetFourCheck(resultList)
	if betFourStatus == true && wr.CheckSlice(resultList) == true {
		return true
	}
	return false
}

//check six digit combination
func (wr *WinningResult) SixDigitCombination(resultType int64) bool {
	resultValue := wr.Value
	if resultType == FIRSTSIXDIGIT && resultValue >= 1 && resultValue <= 6 {
		return true
	} else if resultType == SECONDSIXDIGIT && resultValue >= 4 && resultValue <= 9 {
		return true
	} else if resultType == THIRDSIXDIGIT && resultValue >= 7 && resultValue <= 12 {
		return true
	} else if resultType == FOURSIXDIGIT && resultValue >= 10 && resultValue <= 15 {
		return true
	} else if resultType == FIVESIXDIGIT && resultValue >= 13 && resultValue <= 18 {
		return true
	} else if resultType == SIXSIXDIGIT && resultValue >= 16 && resultValue <= 21 {
		return true
	} else if resultType == SEVENSIXDIGIT && resultValue >= 19 && resultValue <= 24 {
		return true
	} else if resultType == EIGHTSIXDIGIT && resultValue >= 22 && resultValue <= 27 {
		return true
	} else if resultType == NINESIXDIGIT && resultValue >= 25 && resultValue <= 30 {
		return true
	} else if resultType == TENSIXDIGIT && resultValue >= 28 && resultValue <= 33 {
		return true
	} else if resultType == ELEVENSIXDIGIT && resultValue >= 31 && resultValue <= 36 {
		return true
	}
	return false
}

//Eat all
func (wr *WinningResult) EatAll(resultType int64) bool {
	switch resultType {
	case REDTYPE, BLACKTYPE, SINGLETYPE, DOUBLETYPE, BIGTYPE, SMALLTYPE, FIRSETAREA,
		SECONDAREA, THIRDAREA, FIRSTINLINE, SECONDINLINE, THIRDINLINE:
		return false
	}
	return true
}
