package main

import (
	_ "github.com/davyxu/cellnet/codec/gogopb"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/golog"
	"landlord_go/basefx"
	"landlord_go/basefx/model"
	_ "landlord_go/svc/game/json"
	_ "landlord_go/svc/game/verify"
	"landlord_go/svc/hub/api"
	"landlord_go/svc/hub/status"
	"time"
)

var log = golog.New("main")

func main() {

	basefx.Init("game")

	basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
		SvcName:     "game",
		NetProcName: "svc.backend",
		ListenAddr:  ":6792",
	})

	hubapi.ConnectToHub(func() {

		// 开始接收game状态
		hubstatus.StartSendStatus("game_status", time.Second*3, func() int {
			return 100
		})
	})

	basefx.StartLoop(nil)

	basefx.Exit()
}
