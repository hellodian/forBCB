package myroulette

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
)

//Roulette This is struct of contract
//@:contract:Roulette
//@:version:1.0
//@:organization:orgNUjCm1i8RcoW2kVTbDw4vKW6jzfMxewJHjkhuiduhjuikjuyhnnjkuhujk111
//@:author:c5a13e26dedbca637ec4af527ae112361bf867ac886e9c4882d4be3bb6cce355
type MyRoulette struct {
	sdk sdk.ISmartContract

	//@:public:store:cache
	publicKey types.PubKey // Check to sign the public key

	//@:public:store:cache
	lockedInBets map[string]bn.Number // Lock amount (unit cong) key: currency name

	//@:public:store:cache
	setting *Setting

	//@:public:store:cache
	recFeeInfo []RecFeeInfo

	//@:public:store
	roundInfo map[string]*RoundInfo
	//@:public:store
	betInfo map[string]map[string]*BetInfo
}

func (mr *MyRoulette) LockedInBetsInit(tokenNameList []types.Address) {
	for _, value := range tokenNameList {
		mr._setLockedInBets(value, bn.N(0))
	}
}

type Setting struct {
	MaxProfit           int64    `json:"maxProfit"`           // Maximum winning amount (cong)
	MaxLimit            int64    `json:"maxLimit"`            // Maximum bet limit (cong)
	MinLimit            int64    `json:"minLimit"`            // Minimum bet limit unit (cong)
	FeeRatio            int64    `json:"feeRatio"`            // Percentage of handling fee after winning the lottery (thousand-point ratio)
	FeeMiniNum          int64    `json:"feeMiniNum"`          // Minimum handling charge
	SendToCltRatio      int64    `json:"sendToCltRatio"`      // Part of the handling fee sent to CLT (thousandths)
	BetExpirationBlocks int64    `json:"betExpirationBlocks"` // Timeout block interval
	TokenNames          []string `json:"tokenNames"`          // List of supported token names
}

type RecFeeInfo struct {
	RecFeeRatio int64         `json:"recFeeRatio"` // Commission allocation ratio
	RecFeeAddr  types.Address `json:"recFeeAddr"`  // List of addresses to receive commissions
}

type RoundInfo struct {
	Commit              []byte               `json:"commit"`              // Current game random number hash value
	TotalBuyAmount      map[string]bn.Number `json:"totalBuyAmount"`      // Current total bet amount map key：tokenName(The name of the currency)
	TotalBetCount       int64                `json:"totalBetCount"`       // Current total number of bets
	State               int64                `json:"state"`               // Current wheel state 0 Not the lottery 1Has the lottery 2 refunded 3In the lottery
	ProcessCount        int64                `json:"processCount"`        // Current status processing bet quantity (settlement, refund subscript)
	FirstBlockHeight    int64                `json:"firstBlockHeight"`    // Block height when the current wheel initializes to determine whether to timeout
	Setting             *Setting             `json:"settings"`            // Configuration information for the current wheel
	BetInfoSerialNumber []types.Address      `json:"betInfoSerialNumber"` // The current wheel betInfo is associated with the serial number
	WinningResult       *WinningResult       `json:"winningResult"`       // The current round of lottery results
}

type BetInfo struct {
	TokenName string        `json:"tokenName"` // Players bet on currency names
	Gambler   types.Address `json:"gambler"`   // Player betting address
	Amount    bn.Number     `json:"amount"`    // Players bet the total amount
	BetData   []BetData     `json:"betData"`   // Player betting details
	WinAmount bn.Number     `json:"winAmount"` // Players this bet the largest bonus
	Settled   bool          `json:"settled"`   // Whether the current bet has been settled
}

type BetData struct {

	BetMode   int64     `json:"betMode"`   // Betting plan
	BetValue  []int64   `json:"betValue"`  // Betting value
	BetAmount bn.Number `json:"betAmount"` // Betting amount
}
