package main

import (
	"github.com/davyxu/cellmesh/service"
	_ "github.com/davyxu/cellnet/codec/gogopb"
	_ "github.com/davyxu/cellnet/codec/protoplus"
	_ "github.com/davyxu/cellnet/peer/gorillaws"
	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/golog"
	"landlord_go/basefx"
	"landlord_go/basefx/model"
	_ "landlord_go/proto" // 进入协议
	_ "landlord_go/svc/agent/backend"
	"landlord_go/svc/agent/frontend"
	"landlord_go/svc/agent/heartbeat"
	"landlord_go/svc/agent/model"
	"landlord_go/svc/agent/routerule"
	"landlord_go/svc/hub/api"
	"landlord_go/svc/hub/status"
	"time"
)

var log = golog.New("main")

func main() {

	basefx.Init("agent")

	//routerule.Download()

	routerule.GetRouteRule()

	heartbeat.StartCheck()

	model.AgentSvcID = service.GetLocalSvcID()

	// 要连接的服务列表
	basefx.CreateCommnicateConnector(fxmodel.ServiceParameter{
		SvcName:      "game",
		NetProcName:  "agent.backend",
	})

	switch *fxmodel.FlagCommunicateType {
	case "tcp":
		frontend.Start(model.FrontendParameter{
			SvcName:     "agent",
			ListenAddr:  ":6790",
			NetPeerType: "tcp.Acceptor",
			NetProcName: "tcp.frontend",
		})
	case "ws":
		frontend.Start(model.FrontendParameter{
			SvcName:     "agent",
			ListenAddr:  ":6791",
			NetPeerType: "gorillaws.Acceptor",
			NetProcName: "ws.frontend",
		})
	}

	hubapi.ConnectToHub(func() {

		// 发送网关连接数量
		hubstatus.StartSendStatus("agent_status", time.Second*3, func() int {
			return model.FrontendSessionManager.SessionCount()
		})
	})

	basefx.StartLoop(nil)

	frontend.Stop()

	basefx.Exit()
}
