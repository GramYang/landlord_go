package hubstatus

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/timer"
	"landlord_go/basefx/model"
	"landlord_go/proto"
	"landlord_go/svc/hub/api"
	"landlord_go/svc/hub/model"
	"time"
)

var (
	recvLoop *timer.Loop
)

const (
	statusUpdateTimeout = time.Second * 3
)

func StartRecvStatus(channelNames []string, svcStatusHandler *func(ev cellnet.Event)) {

	for _, channelName := range channelNames {
		hubapi.Subscribe(channelName)
	}

	*svcStatusHandler = func(ev cellnet.Event) {

		msg := ev.Message().(*proto.SvcStatusACK)

		model.UpdateStatus(&model.Status{
			UserCount: msg.UserCount,
			SvcID:     msg.SvcID,
		})
	}

	// 保证可以重入
	if recvLoop == nil {
		recvLoop = timer.NewLoop(fxmodel.Queue, statusUpdateTimeout, func(loop *timer.Loop) {

			now := time.Now()
			var timeoutSvcID []string

			model.VisitStatus(func(status *model.Status) bool {
				if now.Sub(status.LastUpdate) > statusUpdateTimeout {
					timeoutSvcID = append(timeoutSvcID, status.SvcID)
				}

				return true
			})

			for _, svcid := range timeoutSvcID {
				//log.Debugln("remove svc status: ", svcid)
				model.RemoveStatus(svcid)
			}

		}, nil)

		recvLoop.Notify().Start()
	}

}
