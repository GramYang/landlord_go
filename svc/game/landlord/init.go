package landlord

import (
	"github.com/davyxu/golog"
	"reflect"
	"sync"
)

type hallTables []*HallTable
func(h hallTables) Len() int {return len(h)}
func(h hallTables) Swap(i,j int) {h[i],h[j]=h[j],h[i]}
func(h hallTables) Less(i, j int) bool {return h[i].TableNum<h[j].TableNum}

var (
	log = golog.New("landlord")
	Code2bean = make(map[int32]reflect.Type)
	playerMap sync.Map //座位号--*Player
	tableMap sync.Map //牌桌号--*Tabler
	userName2Player sync.Map //用户名--*Player
	clientID2Player sync.Map //cid-*Player
	hallList hallTables//游戏大厅信息
	CARDS = []int32{1,2,3,4,5,6,7,8,9,10,11,12,13,21,22,23,24,25,26,27,28,29,30,31,32,33,41,42,43,
		44,45,46,47,48,49,50,51,52,53,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75}
)

func init() {
	Code2bean[10] = reflect.TypeOf((*CancelReadyRequest)(nil)).Elem()
	Code2bean[11] = reflect.TypeOf((*CardsOutRequest)(nil)).Elem()
	Code2bean[12] = reflect.TypeOf((*ChatMsgRequest)(nil)).Elem()
	Code2bean[13] = reflect.TypeOf((*EndGameRequest)(nil)).Elem()
	Code2bean[14] = reflect.TypeOf((*EndGrabLandlordRequest)(nil)).Elem()
	Code2bean[15] = reflect.TypeOf((*EnterTableRequest)(nil)).Elem()
	Code2bean[16] = reflect.TypeOf((*ExitHallRequest)(nil)).Elem()
	Code2bean[17] = reflect.TypeOf((*ExitSeatRequest)(nil)).Elem()
	Code2bean[18] = reflect.TypeOf((*GiveUpLandlordRequest)(nil)).Elem()
	Code2bean[19] = reflect.TypeOf((*InitHallRequest)(nil)).Elem()
	Code2bean[20] = reflect.TypeOf((*LandlordMultipleWagerRequest)(nil)).Elem()
	Code2bean[21] = reflect.TypeOf((*LoginRequest)(nil)).Elem()
	Code2bean[22] = reflect.TypeOf((*MultipleWagerRequest)(nil)).Elem()
	Code2bean[23] = reflect.TypeOf((*ReadyRequest)(nil)).Elem()
	Code2bean[24] = reflect.TypeOf((*UserInfoRequest)(nil)).Elem()
	Code2bean[25] = reflect.TypeOf((*GameResultRequest)(nil)).Elem()

	Code2bean[100] = reflect.TypeOf((*CancelReadyResponse)(nil)).Elem()
	Code2bean[101] = reflect.TypeOf((*CardsOutResponse)(nil)).Elem()
	Code2bean[102] = reflect.TypeOf((*ChatMsgResponse)(nil)).Elem()
	Code2bean[103] = reflect.TypeOf((*EndGameResponse)(nil)).Elem()
	Code2bean[104] = reflect.TypeOf((*EndGrabLandlordResponse)(nil)).Elem()
	Code2bean[105] = reflect.TypeOf((*EnterTableResponse)(nil)).Elem()
	Code2bean[106] = reflect.TypeOf((*ExitSeatResponse)(nil)).Elem()
	Code2bean[107] = reflect.TypeOf((*GiveUpLandlordResponse)(nil)).Elem()
	Code2bean[108] = reflect.TypeOf((*GrabLandlordResponse)(nil)).Elem()
	Code2bean[109] = reflect.TypeOf((*InitHallResponse)(nil)).Elem()
	Code2bean[110] = reflect.TypeOf((*LandlordMultipleWagerResponse)(nil)).Elem()
	Code2bean[111] = reflect.TypeOf((*LoginResponse)(nil)).Elem()
	Code2bean[112] = reflect.TypeOf((*MultipleWagerResponse)(nil)).Elem()
	Code2bean[113] = reflect.TypeOf((*ReadyResponse)(nil)).Elem()
	Code2bean[114] = reflect.TypeOf((*RefreshHallResponse)(nil)).Elem()
	Code2bean[115] = reflect.TypeOf((*UserInfoResponse)(nil)).Elem()
	Code2bean[116] = reflect.TypeOf((*GameResultResponse)(nil)).Elem()

	for i:=1; i<=300; i++ {
		playerMap.Store(int32(i), &Player{})
	}
	for j:=1; j<=100; j++ {
		tableMap.Store(int32(j), &Table{
			ThrowOutCards:make(map[int32]int32),
			PlayersCardsOut:make(map[int32]int32),
		})
	}
	hallList = make(hallTables, 100)
}