package main

import (
	"github.com/davyxu/cellnet"
	_ "github.com/davyxu/cellnet/codec/gogopb"
	_ "github.com/davyxu/cellnet/peer/gorillaws"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/proc/gorillaws"
	"github.com/davyxu/golog"
	"landlord_go/basefx"
	fxmodel "landlord_go/basefx/model"
	"landlord_go/proto"
	hubapi "landlord_go/svc/hub/api"
	hubstatus "landlord_go/svc/hub/status"
	_ "landlord_go/svc/login/login"
)

var log = golog.New("main")

func main() {

	basefx.Init("login")

	// 与客户端通信的处理器
	proc.RegisterProcessor("ws.client", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {

		bundle.SetTransmitter(new(gorillaws.WSMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(gorillaws.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})

	switch *fxmodel.FlagCommunicateType {
	case "tcp":
		basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
			SvcName:     "login",
			NetPeerType: "tcp.Acceptor",
			NetProcName: "tcp.ltv",
			ListenAddr:  ":6789",
		})
	case "ws":
		basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
			SvcName:     "login",
			NetPeerType: "gorillaws.Acceptor",
			NetProcName: "ws.client",
			ListenAddr:  ":6790",
		})
	}

	hubapi.ConnectToHub(func() {

		// 开始接收game状态
		hubstatus.StartRecvStatus([]string{"game_status", "agent_status"}, &proto.Handle_Login_SvcStatusACK)
	})

	basefx.StartLoop(nil)

	basefx.Exit()
}
