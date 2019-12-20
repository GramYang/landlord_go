package agentapi

import (
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet"
	"landlord_go/proto"
	"landlord_go/svc/agent/backend"
)

// 传入用户处理网关消息回调,返回消息源回调
func HandleBackendMessage(userHandler func(ev cellnet.Event, cid proto.ClientID)) func(ev cellnet.Event) {

	return func(incomingEv cellnet.Event) {

		switch ev := incomingEv.(type) {
		case *backend.RecvMsgEvent:

			var cid proto.ClientID
			cid.ID = ev.ClientID

			if agentCtx := service.SessionToContext(ev.Session()); agentCtx != nil {
				cid.SvcID = agentCtx.SvcID
			}

			userHandler(incomingEv, cid)
		}
	}
}
