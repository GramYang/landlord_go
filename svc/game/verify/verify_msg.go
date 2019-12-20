package verify

import (
	"fmt"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"landlord_go/proto"
	"landlord_go/svc/agent/api"
)

func init() {

	proto.Handle_Game_VerifyREQ = agentapi.HandleBackendMessage(func(ev cellnet.Event, cid proto.ClientID) {

		msg := ev.Message().(*proto.VerifyREQ)

		fmt.Printf("verfiy: %+v \n", msg.GameToken)

		service.Reply(ev, &proto.VerifyACK{Result:proto.ResultCode_NoError})
	})
}
