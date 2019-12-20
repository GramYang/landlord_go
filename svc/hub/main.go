package main

import (
	_ "github.com/davyxu/cellnet/codec/gogopb"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/golog"
	"landlord_go/basefx"
	"landlord_go/basefx/model"
	_ "landlord_go/svc/hub/subscribe"
)

var log = golog.New("main")

func main() {

	basefx.Init("hub")

	basefx.CreateCommnicateAcceptor(fxmodel.ServiceParameter{
		SvcName:     "hub",
		NetProcName: "tcp.svc",
		ListenAddr:  ":6791",
	})

	basefx.StartLoop(nil)

	basefx.Exit()
}
