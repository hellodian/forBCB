package excellencies

import (
	"blockchain/smcsdk/sdk"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/smcsdk/sdk/types"
	"fmt"
	"sort"
	"blockchain/smcsdk/sdk/crypto/sha3"
)

type Poker struct {
	Name string
	Flag uint8
	Type  uint8  //牌的颜色 3 2 1 0 分别表示黑 红 梅 方
}

func NewPoker(name string, flag uint8) *Poker {
	return &Poker{
		Name: name,
		Flag: flag,
	}
}

type PokerSet struct {
	No     string
	Pokers []*Poker
	IsWin  bool
	TypeSan  uint8  //THREESAME MIXTHREE THREENOSAME  代表三张牌一样  混三公  点数
	ThreeSum uint8 // three porkers in sum
}

func (ps *PokerSet) Len() int {
	return len(ps.Pokers)
}

func (ps *PokerSet) Less(i, j int) bool {
	return ps.Pokers[i].Flag > ps.Pokers[j].Flag
}

func (ps *PokerSet) Swap(i, j int) {
	ps.Pokers[i], ps.Pokers[j] = ps.Pokers[j], ps.Pokers[i]
}

func NewPokerSet(no string) *PokerSet {
	return &PokerSet{
		No:     no,
		Pokers: make([]*Poker, 0),
	}
}

func (ps *PokerSet) addPai(p *Poker) {
	ps.Pokers = append(ps.Pokers, p)
}

type PokerPool struct {
	Pool []*Poker
}

func NewPokerPool(hash []byte,poker []*Poker,e *Excellencies) *PokerPool {
	pool := &PokerPool{
		Pool: make([]*Poker, 0),
	}
	pool.convert(hash,poker,e)
	return pool
}

func (pp *PokerPool) Discard() *Poker {
	pokers := pp.Pool
	i := len(pokers)
	pp.Pool = pokers[:i-1]
	return pokers[i-1]
}

func (pp *PokerPool) convert(hash []byte,poker []*Poker,e *Excellencies) {
	bytes := []byte(hash)
	for i := 0; i < 15; i++ {
		w:=bn.NBytes(sha3.Sum256(bytes[2:], e.sdk.Block().BlockHash(), e.sdk.Block().RandomNumber()))
		bytes=bytes[2:]
		index:=w.ModI(int64(POKERSIZE-i)).V.Int64()
		pp.Pool = append(pp.Pool, poker[index])
		poker=append(poker[0:index],poker[index+1:]...)

	}
}



type Game struct {
	Gamer map[string]*PokerSet
	Pp    *PokerPool
}

func NewGame(pp *PokerPool) *Game {
	game := &Game{
		Gamer: make(map[string]*PokerSet, 0),
		Pp:    pp,
	}
	game.init()
	return game
}

func (g *Game) init() {
	sets := g.Gamer
	sets["A"] = NewPokerSet("A")
	sets["B"] = NewPokerSet("B")
	sets["C"] = NewPokerSet("C")
	sets["D"] = NewPokerSet("D")
	sets["E"] = NewPokerSet("E")
}

func (g *Game) start() {
	sta := [5]string{"B", "C", "D", "E", "A"}
	pool := g.Pp
	len := len(pool.Pool)
	index := 0
	for i := 0; i < len; i++ {
		s := sta[index]
		set := g.Gamer[s]
		set.addPai(g.Pp.Discard())
		index++
		if index == 5 {
			index = 0
		}
	}

}

func (g *Game) GetGamerCards(flag string) *PokerSet {
	set := g.Gamer[flag]
	sort.Sort(set)
	return set
}

//func (e *Excellencies) StartGame() *Game {
func (e *Excellencies) StartGame(reveal []byte) (*Game,[]*Poker) {
	w:=bn.NBytes(sha3.Sum256(reveal, e.sdk.Block().BlockHash(), e.sdk.Block().RandomNumber()))
	//重新洗牌
	poker:=e.ShufflePoker(w)
	hash := e.sdk.Block().BlockHash()
	shash := fmt.Sprint(hash)
	pool := NewPokerPool([]byte(shash),poker,e)
	result:=pool.Pool
	game := NewGame(pool)
	game.start()
	return game,result
}

func (e *Excellencies) ShufflePoker(w bn.Number) []*Poker{
	poker:=e._poker()
	newPoker:=make([]*Poker,0)

	for i:=0;i<POKERSIZE;i++{
		if i==POKERSIZE{
			newPoker=append(newPoker,poker[0])
		}
		//算随机数
		index:=w.ModI(int64(POKERSIZE-i)).V.Int64()
		value:=poker[index]
		newPoker=append(newPoker,value)
		poker=append(poker[0:index],poker[index+1:]...)
	}
	e._setPoker(newPoker)
	return newPoker


}
func (e *Excellencies) MaybeWinAmountAndFeeByList(amount bn.Number, betList []BetData) (totalMaybeWinAmount, fee bn.Number) {
	fee = bn.N(0)
	totalMaybeWinAmount = bn.N(0)
	settings := e._settings()
	// fee
	fee = amount.MulI(settings.FeeRatio).DivI(PERMILLE)
	if fee.CmpI(settings.FeeMiniNum) < 0 {
		fee = bn.N(settings.FeeMiniNum)
	}

	sdk.Require(fee.Cmp(amount) <= 0,
		types.ErrInvalidParameter, "Bet doesn't even cover fee")

	for _, bet := range betList {
		totalMaybeWinAmount = totalMaybeWinAmount.Add(bet.BetAmount.MulI(ODDS))
	}
	return
}

func (e *Excellencies) GetGameResult(reveal []byte,roundinfo *RoundInfo ) *Game {

	game,pool:= e.StartGame(reveal)
	roundinfo.OriginPokers=pool
	banker := game.GetGamerCards(BANKER)
	player1 := game.GetGamerCards(PLAYER_1)
	player2 := game.GetGamerCards(PLAYER_2)
	player3 := game.GetGamerCards(PLAYER_3)
	player4 := game.GetGamerCards(PLAYER_4)
	banker.CalcPokerTypeAndSum()
	banker.JudgeBankerAndPlayerWin(player1)
	banker.JudgeBankerAndPlayerWin(player2)
	banker.JudgeBankerAndPlayerWin(player3)
	banker.JudgeBankerAndPlayerWin(player4)
	return game
}
func (ps *PokerSet) CalcPokerTypeAndSum() {
	ps.CalcPokerType()
	ps.CountSumPoints()
}
func (ps *PokerSet) CalcPokerType() {

	typeRueslt:=ps.JudgePokersIdentical()
	if typeRueslt {
		ps.TypeSan=THREESAME
		return
	}
	typeRueslt=ps.JudgePokersMixedThree()
	if typeRueslt {
		ps.TypeSan=MIXEDTHREE
		return
	}
	ps.TypeSan=THREENOSAME
	return

}
//return true banker win
func (ps *PokerSet) JudgeBankerAndPlayerWin(pr *PokerSet) (result bool) {

	pr.CalcPokerTypeAndSum()
	//先比较类型
	if ps.TypeSan>pr.TypeSan{
		pr.IsWin=false
		return
	}
	if ps.TypeSan<pr.TypeSan{
		pr.IsWin=true
		return
	}
	//类型相同
	if ps.TypeSan==pr.TypeSan{
		if ps.ThreeSum>pr.ThreeSum{
			pr.IsWin=false
			return
		}

		if ps.ThreeSum<pr.ThreeSum{
			pr.IsWin=true
			return
		}

		if ps.ThreeSum==pr.ThreeSum{
			//比较公牌数量
			if ps.CountPaiSum()>pr.CountPaiSum(){
				pr.IsWin=false
				return
			}
			if ps.CountPaiSum()<pr.CountPaiSum(){
				pr.IsWin=true
				return
			}
			if ps.CountPaiSum()==pr.CountPaiSum(){//公牌数相等
				//比较最大那张牌
				if ps.Pokers[0].Flag>ps.Pokers[0].Flag{
					pr.IsWin=false
					return
				}
				if ps.Pokers[0].Flag<ps.Pokers[0].Flag{
					pr.IsWin=true
					return
				}
				if ps.Pokers[0].Flag==ps.Pokers[0].Flag{
					//比较花色
					if ps.Pokers[0].Type>pr.Pokers[0].Type{
						pr.IsWin=false

					}else{
						pr.IsWin=true
					}
					return
				}


			}
			//比较花色
			if ps.Pokers[0].Type>pr.Pokers[0].Type{
				pr.IsWin=false
				return
			}else {
				pr.IsWin=true
				return
			}
		}

	}




	return

}

func (ps *PokerSet) JudgePokersIdentical() (result bool) {

	result = ps.Pokers[0].Flag == ps.Pokers[1].Flag && ps.Pokers[1].Flag == ps.Pokers[2].Flag
	return
}

func (ps *PokerSet) JudgePokersMixedThree() (result bool) {

	result = ps.Pokers[2].Flag > 10
	return
}

func (ps *PokerSet) JudgeMixedThree() (result bool) {

	return ps.Pokers[0].Flag>uint8(10) && ps.Pokers[1].Flag>uint8(10) && ps.Pokers[2].Flag>uint8(10)


}


//算三张牌的总点数
func (ps *PokerSet) CountSumPoints() {

	pokers := ps.Pokers
	sum:=uint8(0)

	for i := 0; i < len(pokers); i++ {
		if pokers[i].Flag<10{
			sum += pokers[i].Flag
		}
	}
	ps.ThreeSum=sum%10
	return
}

//算公牌数量
func (ps *PokerSet) CountPaiSum() uint8 {

	pokers := ps.Pokers
	sum:=uint8(0)

	for i := 0; i < len(pokers); i++ {
		if pokers[i].Flag>10{
			sum ++
		}
	}

	return sum
}


func (g *Game) GetBetWinAmount(data []BetData) (winAmount, winFeeAmount bn.Number) {

	winAmount = bn.N(0)
	winFeeAmount = bn.N(0)
	for _, bet := range data {

		if g.Gamer[bet.BetMode].IsWin {
			winAmount = winAmount.Add(bet.BetAmount.MulI(ODDS))
			winFeeAmount = winFeeAmount.Add(bet.BetAmount)
		}

	}
	return
}
