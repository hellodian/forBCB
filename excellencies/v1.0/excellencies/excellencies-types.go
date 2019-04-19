package excellencies

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/jsoniter"
	"blockchain/smcsdk/sdk/types"
)

func (sg *Excellencies) LockedInBetsInit(tokenNameList []types.Address) {
	for _, value := range tokenNameList {
		sg._setLockedInBets(value, bn.N(0))
	}
}

func (sg *Excellencies) createLimitMaps(tokenNameList []types.Address) map[string]Limit {
	limitMaps := make(map[string]Limit, len(tokenNameList))
	for _, value := range tokenNameList {
		limit := Limit{
			MaxProfit: 2E12,
			MaxLimit:  2E10,
			MinLimit:  1E8,
		}
		limitMaps[value] = limit
	}
	return limitMaps
}

type Settings struct {
	LimitMaps           map[string]Limit `json:"limitMaps"`           // Maximum winning amount、Maximum bet、Minimum bet limit (cong)
	FeeRatio            int64            `json:"feeRatio"`            // Percentage of handling fee after winning the lottery (thousand-point ratio)
	FeeMiniNum          int64            `json:"feeMiniNum"`          // Minimum handling charge
	SendToCltRatio      int64            `json:"sendToCltRatio"`      // Part of the handling fee sent to CLT (thousandths)
	BetExpirationBlocks int64            `json:"betExpirationBlocks"` // Timeout block interval
	TokenNames          []string         `json:"tokenNames"`          // List of supported token names
	PoolFeeRatio        int64            `json:"poolFeeRatio"`		 //Pool amount ratio
	CarveUpPoolRatio    int64            `json:"CarveUpPoolRatio"`
}

type Limit struct {
	MaxProfit int64 `json:"maxProfit"` // Maximum winning amount (cong)
	MaxLimit  int64 `json:"maxLimit"`  // Maximum bet limit (cong)
	MinLimit  int64 `json:"minLimit"`  // Minimum bet limit unit (cong)
}

type RecFeeInfo struct {
	RecFeeRatio int64         `json:"recFeeRatio"` // Commission allocation ratio
	RecFeeAddr  types.Address `json:"recFeeAddr"`  // List of addresses to receive commissions
}

type RoundInfo struct {
	Commit           []byte               `json:"commit"`           // Current game random number hash value
	TotalBuyAmount   map[string]bn.Number `json:"totalBuyAmount"`   // Current total bet amount map key：tokenName(The name of the currency)
	TotalBetCount    int64                `json:"totalBetCount"`    // Current total number of bets
	State            int64                `json:"state"`            // Current wheel state 0 Not the lottery 1Has the lottery 2 refunded 3In the lottery
	ProcessCount     int64                `json:"processCount"`     // Current status processing bet quantity (settlement, refund subscript)
	FirstBlockHeight int64                `json:"firstBlockHeight"` // Block height when the current wheel initializes to determine whether to timeout
	Settings         Settings             `json:"settings"`         // Configuration information for the current wheel
	BetInfos         []BetInfo            `json:"betInfos"`         // Key1 player address, key2 currency name
	Game             Game                 `json:"game"`             // Game Result, Maker and Player Type
	PoolFlag         bool                 `json:"poolFlag"`          // 大小三公 获得额外的奖励是否结算
	BetAmountByModel map[string]map[string]bn.Number   `json:"betAmountByModel"`   //不同牌型投注总额
	OriginPokers     []*Poker			  `json:"originPokers"`     //原始牌面，没经过排序的牌
}

func NewRoundInfo(commit []byte, settings Settings) *RoundInfo {
	return &RoundInfo{
		Commit:           commit,
		TotalBuyAmount:   make(map[string]bn.Number),
		TotalBetCount:    0,
		State:            0,
		ProcessCount:     0,
		FirstBlockHeight: 0,
		Settings:         settings,
		BetInfos:         make([]BetInfo, 0),
		Game:             Game{},
		PoolFlag:false,
		BetAmountByModel:make(map[string]map[string]bn.Number),
		OriginPokers:make([]*Poker,0),

	}
}

type BetInfo struct {
	TokenName string        `json:"tokenName"` // Players bet on currency names
	Gambler   types.Address `json:"gambler"`   // Player betting address
	Amount    bn.Number     `json:"amount"`    // Players bet the total amount
	BetData   []BetData     `json:"betData"`   // Player betting details
	WinAmount bn.Number     `json:"winAmount"` // Players this bet the largest bonus
	Settled   bool          `json:"settled"`   // Whether the current bet has been settled
}

func NewBetInfo(tokenName string, gambler types.Address) *BetInfo {
	return &BetInfo{
		TokenName: tokenName,
		Gambler:   gambler,
		Amount:    bn.N(0),
		BetData:   make([]BetData, 0),
		WinAmount: bn.N(0),
		Settled:   false,
	}
}

func (bi *BetInfo) UpdateBetInfo(bet BetData) {
	data := bi.BetData
	bi.BetData = append(data, bet)
	number := bet.BetAmount
	bi.Amount = bi.Amount.Add(number)
}

type SlipperInfo struct {
	SettlementTime  bn.Number     `json:"settlementTime"`	//Clearing time
	PlayersAddress  types.Address `json:"playersAddress"`	//Players address
	TokenName		string        `json:"tokenName"`        //Players bet on currency names
	SettlementMoney bn.Number     `json:"settlementMoney"`	//The settlement amount
}


type BetData struct {
	BetMode   string    `json:"betMode"`   // A means banker B, C, D, E means idle family
	BetAmount bn.Number `json:"betAmount"` // Betting amount
}

func BuildBetData(jsonData []byte) []BetData {
	data := make([]BetData, 0)
	unmarshal := jsoniter.Unmarshal(jsonData, &data)
	sdk.RequireNotError(unmarshal, types.ErrInvalidParameter)
	return data
}


func CheckList(t string,s []string) (flag bool) {
	for _,v :=range s{
		if v==t {
			flag=true
			return
		}
	}
	return
}

