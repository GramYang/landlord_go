package landlord

type CancelReadyRequest struct {
	IsCancelReady bool
}

type CardsOutRequest struct {
	IsPass bool
	FromSeatNum int32
	ToSeatNum int32
	CardsOut []int32
}

type ChatMsgRequest struct {
	ChatFlag int32
	UserName string
	Msg string
	TableNum int32
}

type EndGameRequest struct {
	WinnerSeatNum int32
}

type EndGrabLandlordRequest struct {
	MeSeatNum int32
}

type EnterTableRequest struct {
	UserName string
	TableNum int32
}

type ExitHallRequest struct {
	UserName string
}

type ExitSeatRequest struct {
	YourSeatNum int32
}

type GiveUpLandlordRequest struct {
	SeatNum int32
}

type InitHallRequest struct {}

type LandlordMultipleWagerRequest struct {
	LandlordSeatNum int32
	MultipleNum int32
}

type LoginRequest struct {
	UserName string
	Password string
}

type MultipleWagerRequest struct {
	Agreed int32
}

type ReadyRequest struct {
	IsReady bool
}

type UserInfoRequest struct {
	UserName string
}

type GameResultRequest struct {
	Result bool
	UserName string
	Password string
	Money int32
}