package hubapi

import (
	"github.com/davyxu/cellnet/relay"
	"landlord_go/basefx"
	"landlord_go/basefx/model"
	"landlord_go/proto"
	"landlord_go/svc/hub/model"
)

// 传入你的服务名, 连接到hub
func ConnectToHub(hubReady func()) {

	model.OnHubReady = hubReady
	basefx.CreateCommnicateConnector(fxmodel.ServiceParameter{
		SvcName:      "hub",
		NetProcName:  "tcp.hub",
		MaxConnCount: 1,
	})
}

func Subscribe(channel string) {

	if model.HubSession == nil {
		log.Errorf("hub session not ready, channel: %s", channel)
		return
	}

	model.HubSession.Send(&proto.SubscribeChannelREQ{
		Channel: channel,
	})
}

func Publish(channel string, msg interface{}) {

	if model.HubSession == nil {
		log.Errorf("hub session not ready, channel: %s", channel)
		return
	}

	_ = relay.Relay(model.HubSession, msg, channel)
}
