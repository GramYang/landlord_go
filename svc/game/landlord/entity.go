package landlord

import "landlord_go/proto"

type HallTable struct {
	TableNum int32
	UserNames []string
	IsPlay bool
	IsFull bool
}

type Player struct {
	Cid *proto.ClientID
	UserName string
	SeatNum int32
	Cards []int32
	TableNum int32
}

type Table struct {
	Players []*Player
	Landlord *Player
	Cards []int32
	TableNum int32
	PlayerCount int32
	ThrowOutCards map[int32]int32
	LastCardsOut []int32
	PlayersCardsOut map[int32]int32
	ThreeCards []int32
	ReadyCount int32
	PassLandlordCount int32
	WagerMultipleNum int32
	AnswerMultipleNum int32
	AgreedMultipleResult int32
	ContinuousPass int32
	IsWait bool
	IsGrab bool
	IsPlay bool
}