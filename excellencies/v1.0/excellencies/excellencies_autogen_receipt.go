package excellencies

import (
	"blockchain/smcsdk/sdk/types"
	"blockchain/smcsdk/sdk/bn"
)

var _ receipt = (*Excellencies)(nil)

//emitSetSecretSigner This is a method of Excellencies
func (e *Excellencies) emitSetSecretSigner(newSecretSigner types.PubKey) {
	type setSecretSigner struct {
		NewSecretSigner types.PubKey `json:"newSecretSigner"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		setSecretSigner{
			NewSecretSigner: newSecretSigner,
		},
	)
}

//emitSetOwner This is a method of Excellencies
func (e *Excellencies) emitSetOwner(newAddress types.Address) {
	type setOwner struct {
		NewAddress types.Address `json:"newAddress"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		setOwner{
			NewAddress: newAddress,
		},
	)
}

//emitSetSettings This is a method of Excellencies
func (e *Excellencies) emitSetSettings(tokenNames []string, limit map[string]Limit, feeRatio, feeMiniNum, sendToCltRatio, betExpirationBlocks, CarveUpPoolRatio int64) {
	type setSettings struct {
		TokenNames          []string         `json:"tokenNames"`
		Limit               map[string]Limit `json:"limit"`
		FeeRatio            int64            `json:"feeRatio"`
		FeeMiniNum          int64            `json:"feeMiniNum"`
		SendToCltRatio      int64            `json:"sendToCltRatio"`
		BetExpirationBlocks int64            `json:"betExpirationBlocks"`
		CarveUpPoolRatio    int64            `json:"CarveUpPoolRatio"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		setSettings{
			TokenNames:          tokenNames,
			Limit:               limit,
			FeeRatio:            feeRatio,
			FeeMiniNum:          feeMiniNum,
			SendToCltRatio:      sendToCltRatio,
			BetExpirationBlocks: betExpirationBlocks,
			CarveUpPoolRatio:    CarveUpPoolRatio,
		},
	)
}

//emitPlaceBet This is a method of Excellencies
func (e *Excellencies) emitPlaceBet(tokenName string, gambler types.Address, totalMaybeWinAmount bn.Number, betDataList []BetData, commitLastBlock int64, commit, signData []byte, refAddress types.Address) {
	type placeBet struct {
		TokenName           string        `json:"tokenName"`
		Gambler             types.Address `json:"gambler"`
		TotalMaybeWinAmount bn.Number     `json:"totalMaybeWinAmount"`
		BetDataList         []BetData     `json:"betDataList"`
		CommitLastBlock     int64         `json:"commitLastBlock"`
		Commit              []byte        `json:"commit"`
		SignData            []byte        `json:"signData"`
		RefAddress          types.Address `json:"refAddress"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		placeBet{
			TokenName:           tokenName,
			Gambler:             gambler,
			TotalMaybeWinAmount: totalMaybeWinAmount,
			BetDataList:         betDataList,
			CommitLastBlock:     commitLastBlock,
			Commit:              commit,
			SignData:            signData,
			RefAddress:          refAddress,
		},
	)
}

//emitSetRecFeeInfo This is a method of Excellencies
func (e *Excellencies) emitSetRecFeeInfo(info []RecFeeInfo) {
	type setRecFeeInfo struct {
		Info []RecFeeInfo `json:"info"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		setRecFeeInfo{
			Info: info,
		},
	)
}

//emitWithdrawFunds This is a method of Excellencies
func (e *Excellencies) emitWithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) {
	type withdrawFunds struct {
		TokenName      string        `json:"tokenName"`
		Beneficiary    types.Address `json:"beneficiary"`
		WithdrawAmount bn.Number     `json:"withdrawAmount"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		withdrawFunds{
			TokenName:      tokenName,
			Beneficiary:    beneficiary,
			WithdrawAmount: withdrawAmount,
		},
	)
}

//emitSettleBet This is a method of Excellencies
func (e *Excellencies) emitSettleBet(tokenName []string, commit []byte, game Game, winNumber int64, totalWinAmount map[string]bn.Number, finished bool, poolAmount map[string]bn.Number, OriginPokers []*Poker) {
	type settleBet struct {
		TokenName      []string             `json:"tokenName"`
		Commit         []byte               `json:"commit"`
		Game           Game                 `json:"game"`
		WinNumber      int64                `json:"winNumber"`
		TotalWinAmount map[string]bn.Number `json:"totalWinAmount"`
		Finished       bool                 `json:"finished"`
		PoolAmount     map[string]bn.Number `json:"poolAmount"`
		OriginPokers   []*Poker             `json:"OriginPokers"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		settleBet{
			TokenName:      tokenName,
			Commit:         commit,
			Game:           game,
			WinNumber:      winNumber,
			TotalWinAmount: totalWinAmount,
			Finished:       finished,
			PoolAmount:     poolAmount,
			OriginPokers:   OriginPokers,
		},
	)
}

//emitSlipperHighestTransfer This is a method of Excellencies
func (e *Excellencies) emitSlipperHighestTransfer(tokenName string, playersAddress types.Address, winningAmount bn.Number) {
	type slipperHighestTransfer struct {
		TokenName      string        `json:"tokenName"`
		PlayersAddress types.Address `json:"playersAddress"`
		WinningAmount  bn.Number     `json:"winningAmount"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		slipperHighestTransfer{
			TokenName:      tokenName,
			PlayersAddress: playersAddress,
			WinningAmount:  winningAmount,
		},
	)
}

//emitCarveUpPool This is a method of Excellencies
func (e *Excellencies) emitCarveUpPool(countNum int64, amountPoolList map[string]bn.Number) {
	type carveUpPool struct {
		CountNum       int64                `json:"countNum"`
		AmountPoolList map[string]bn.Number `json:"amountPoolList"`
	}

	e.sdk.Helper().ReceiptHelper().Emit(
		carveUpPool{
			CountNum:       countNum,
			AmountPoolList: amountPoolList,
		},
	)
}
        
                         