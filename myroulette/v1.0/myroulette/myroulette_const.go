package myroulette

const (
	REDTYPE               = 1  //red    红
	BLACKTYPE             = 2  //black  黑
	SINGLETYPE            = 3  //single 单
	DOUBLETYPE            = 4  //double 双
	BIGTYPE               = 5  //big    大
	SMALLTYPE             = 6  //small  小
	FIRSETAREA            = 7  //define area   第1区
	SECONDAREA            = 8  //define area   第2区
	THIRDAREA             = 9  //define area   第3区
	FIRSTINLINE           = 10 //inlie    第一直列
	SECONDINLINE          = 11 //inlie    第二直列
	THIRDINLINE           = 12 //inlie    第三直列
	ASINGLENUMBER         = 13 //A single number 单个数字
	TWODIGITCOMBINATION   = 14 //Two digit combination  两个数字组合
	THREEDIGITCOMBINATION = 15 //Three-digit combination 三个数字组合
	FOURDIGITCOMBINATION  = 16 //Four-digit combination 四个数字组合
	FIVEDIGITCOMBINATION  = 17 //Five-digit combination 五个数字组合
	FIRSTSIXDIGIT         = 18 //Six-digit combination  六个数字组合第一区
	SECONDSIXDIGIT        = 19 //Six-digit combination  六个数字组合第二区
	THIRDSIXDIGIT         = 20 //Six-digit combination  六个数字组合第三区
	FOURSIXDIGIT          = 21 //Six-digit combination  六个数字组合第四区
	FIVESIXDIGIT          = 22 //Six-digit combination  六个数字组合第五区
	SIXSIXDIGIT           = 23 //Six-digit combination  六个数字组合第六区
	SEVENSIXDIGIT         = 24 //Six-digit combination  六个数字组合第七区
	EIGHTSIXDIGIT         = 25 //Six-digit combination  六个数字组合第八区
	NINESIXDIGIT          = 26 //Six-digit combination  六个数字组合第九区
	TENSIXDIGIT           = 27 //Six-digit combination  六个数字组合第十区
	ELEVENSIXDIGIT        = 28 //Six-digit combination  六个数字组合第十一区

	ONETIMES 				= 1
	TWOTIMES 				= 2
	THIRTYTIMES 			= 35
	SEVENTEENTIMES 			= 17
	ELEVENTIMES 			= 11
	EIGHTTIMES 				= 8
	SIXTIMES 				= 6

	RANDOMMODYLO			= 38
)


const (
	NOAWARD       = iota //Not the lottery
	AWARDED              //Has the lottery
	REFUNDED             //refunded
	OPENINGAPRIZE        //Is the lottery
)

const (
	PERMILLE = 1000
)
const (
	DICENUMBERMIN = 0
	DICENUMBERMAX = 37
)