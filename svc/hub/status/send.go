package hubstatus

import (
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/cellnet/timer"
	"landlord_go/basefx/model"
	"landlord_go/proto"
	hubapi "landlord_go/svc/hub/api"
	"time"
)

func StartSendStatus(channelName string, updateInterval time.Duration, statusGetter func() int) {

	timer.NewLoop(fxmodel.Queue, updateInterval, func(loop *timer.Loop) {

		var ack proto.SvcStatusACK
		ack.SvcID = service.GetLocalSvcID()
		ack.UserCount = int32(statusGetter())

		hubapi.Publish(channelName, &ack)

	}, nil).Notify().Start()
}
