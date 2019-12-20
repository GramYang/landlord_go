package landlord

type CancelReadyResponse struct {
	IsCancelReady bool
}

type CardsOutResponse struct {
	IsPass bool
	IsAllPass bool
	FromSeatNum int32
	ToSeatNum int32
	CardsOut []int32
	ThrowOutCards map[int32]int32
	PlayersCardsCount map[int32]int32
}

type ChatMsgResponse struct {
	ChatFlag int32
	UserName string
	Msg string
	TableNum int32
}

type EndGameResponse struct {
	WinnerSeatNum int32
}

type EndGrabLandlordResponse struct {
	FinalLandlordSeatNum int32
	ThreeCards []int32
}

type EnterTableResponse struct {
	IsSuccess bool
	TablePlayers map[int32]string
}

type ExitSeatResponse struct {
	UserName string
	SeatNum int32
	TablePlayers map[int32]string
}

type GiveUpLandlordResponse struct {
	NextLandlordSeatNum int32
}

type GrabLandlordResponse struct {
	LandlordSeatNum int32
	TablePlayers map[int32]string
	ThreeCards []int32
	Cards []int32
}

type InitHallResponse struct {
	HallTables []*HallTable
}

type LandlordMultipleWagerResponse struct {
	MultipleNum int32
}

type LoginResponse struct {
	UserName string
	IsSuccessFul bool
	ResponseMsg string
}

type MultipleWagerResponse struct {
	MultipleNum int32
}

type ReadyResponse struct {
	Ready bool
}

type RefreshHallResponse struct {
	HallTables []*HallTable
}

type UserInfoResponse struct {
	Name string
	Avatar string
	Win string
	Lose string
	Money string
}

type GameResultResponse struct {
	Status int32
}