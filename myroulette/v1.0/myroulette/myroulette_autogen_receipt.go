package myroulette

import (
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

var _ receipt = (*MyRoulette)(nil)

//emitSetPublicKey This is a method of MyRoulette
func (mr *MyRoulette) emitSetPublicKey(newPublicKey types.PubKey) {
	type setPublicKey struct {
		NewPublicKey types.PubKey `json:"newPublicKey"`
	}

	mr.sdk.Helper().ReceiptHelper().Emit(
		setPublicKey{
			NewPublicKey: newPublicKey,
		},
	)
}

//emitSetSetting This is a method of MyRoulette
func (mr *MyRoulette) emitSetSetting(tokenNames []string, minLimit, maxLimit, maxProfit, feeRatio, feeMiniNum, sendToCltRatio, betExpirationBlocks int64) {
	type setSetting struct {
		TokenNames          []string `json:"tokenNames"`
		MinLimit            int64    `json:"minLimit"`
		MaxLimit            int64    `json:"maxLimit"`
		MaxProfit           int64    `json:"maxProfit"`
		FeeRatio            int64    `json:"feeRatio"`
		FeeMiniNum          int64    `json:"feeMiniNum"`
		SendToCltRatio      int64    `json:"sendToCltRatio"`
		BetExpirationBlocks int64    `json:"betExpirationBlocks"`
	}

	mr.sdk.Helper().ReceiptHelper().Emit(
		setSetting{
			TokenNames:          tokenNames,
			MinLimit:            minLimit,
			MaxLimit:            maxLimit,
			MaxProfit:           maxProfit,
			FeeRatio:            feeRatio,
			FeeMiniNum:          feeMiniNum,
			SendToCltRatio:      sendToCltRatio,
			BetExpirationBlocks: betExpirationBlocks,
		},
	)
}

//emitSetRecFeeInfo This is a method of MyRoulette
func (mr *MyRoulette) emitSetRecFeeInfo(info []RecFeeInfo) {
	type setRecFeeInfo struct {
		Info []RecFeeInfo `json:"info"`
	}

	mr.sdk.Helper().ReceiptHelper().Emit(
		setRecFeeInfo{
			Info: info,
		},
	)
}

//emitWithdrawFunds This is a method of MyRoulette
func (mr *MyRoulette) emitWithdrawFunds(tokenName string, beneficiary types.Address, withdrawAmount bn.Number) {
	type withdrawFunds struct {
		TokenName      string        `json:"tokenName"`
		Beneficiary    types.Address `json:"beneficiary"`
		WithdrawAmount bn.Number     `json:"withdrawAmount"`
	}

	mr.sdk.Helper().ReceiptHelper().Emit(
		withdrawFunds{
			TokenName:      tokenName,
			Beneficiary:    beneficiary,
			WithdrawAmount: withdrawAmount,
		},
	)
}

//emitPlaceBet This is a method of MyRoulette
func (mr *MyRoulette) emitPlaceBet(tokenName string, gambler types.Address, totalMaybeWinAmount bn.Number, betDataList []BetData, commitLastBlock int64, commit, signData []byte, refAddress types.Address) {
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

	mr.sdk.Helper().ReceiptHelper().Emit(
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

//emitSettleBet This is a method of MyRoulette
func (mr *MyRoulette) emitSettleBet(tokenName []string, reveal, commit []byte, gambler []types.Address, winNumber int64, totalWinAmount map[string]bn.Number, finished bool) {
	type settleBet struct {
		TokenName      []string             `json:"tokenName"`
		Reveal         []byte               `json:"reveal"`
		Commit         []byte               `json:"commit"`
		Gambler        []types.Address      `json:"gambler"`
		WinNumber      int64                `json:"winNumber"`
		TotalWinAmount map[string]bn.Number `json:"totalWinAmount"`
		Finished       bool                 `json:"finished"`
	}

	mr.sdk.Helper().ReceiptHelper().Emit(
		settleBet{
			TokenName:      tokenName,
			Reveal:         reveal,
			Commit:         commit,
			Gambler:        gambler,
			WinNumber:      winNumber,
			TotalWinAmount: totalWinAmount,
			Finished:       finished,
		},
	)
}

//emitRefundBet This is a method of MyRoulette
func (mr *MyRoulette) emitRefundBet(commit []byte, tokenName []string, gambler []types.Address, refundedAmount map[string]bn.Number, finished bool) {
	type refundBet struct {
		Commit         []byte               `json:"commit"`
		TokenName      []string             `json:"tokenName"`
		Gambler        []types.Address      `json:"gambler"`
		RefundedAmount map[string]bn.Number `json:"refundedAmount"`
		Finished       bool                 `json:"finished"`
	}

	mr.sdk.Helper().ReceiptHelper().Emit(
		refundBet{
			Commit:         commit,
			TokenName:      tokenName,
			Gambler:        gambler,
			RefundedAmount: refundedAmount,
			Finished:       finished,
		},
	)
}
