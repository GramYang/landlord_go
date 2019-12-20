package json

import (
	"encoding/json"
	"github.com/davyxu/cellnet"
	"landlord_go/proto"
	"landlord_go/svc/agent/api"
	"landlord_go/svc/game/landlord"
	"reflect"
)

func init() {
	proto.Handle_Game_JsonREQ = agentapi.HandleBackendMessage(func(ev cellnet.Event, cid proto.ClientID) {
		switch msg := ev.Message().(type) {
		case *proto.JsonREQ:
			log.Infof("receive json message: %v content: %+v\n", landlord.Code2bean[msg.JsonType], string(msg.Content))
			request := reflect.New(landlord.Code2bean[msg.JsonType]).Interface()
			switch msg.JsonType {
			case 21:
				req := request.(*landlord.LoginRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.Login(req, &cid)
			case 19:
				req := request.(*landlord.InitHallRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.InitHall(req, &cid)
			case 15:
				req := request.(*landlord.EnterTableRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.EnterTable(req, &cid)
			case 12:
				req := request.(*landlord.ChatMsgRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.ChatMsg(req, &cid)
			case 23:
				req := request.(*landlord.ReadyRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.Ready(req, &cid)
			case 10:
				req := request.(*landlord.CancelReadyRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.CancelReady(req, &cid)
			case 18:
				req := request.(*landlord.GiveUpLandlordRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.GiveUpLandlord(req, &cid)
			case 14:
				req := request.(*landlord.EndGrabLandlordRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.EndGrabLandlord(req, &cid)
			case 20:
				req := request.(*landlord.LandlordMultipleWagerRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.LandlordMultipleWager(req, &cid)
			case 22:
				req := request.(*landlord.MultipleWagerRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.MultipleWager(req, &cid)
			case 11:
				req := request.(*landlord.CardsOutRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.CardsOut(req, &cid)
			case 13:
				req := request.(*landlord.EndGameRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.EndGame(req, &cid)
			case 17:
				req := request.(*landlord.ExitSeatRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.ExitSeat(req, &cid)
			case 16:
				req := request.(*landlord.ExitHallRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.ExitHall(req, &cid)
			case 24:
				req := request.(*landlord.UserInfoRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.UserInfo(req, &cid)
			case 25:
				req := request.(*landlord.GameResultRequest)
				_ = json.Unmarshal(msg.Content, req)
				landlord.GameResult(req, &cid)
			}
		case *proto.ClientClosedACK:
			log.Infof("client with id: %s has exit", msg.ID)
			landlord.ExitOrException(&cid)
		}
	})
}