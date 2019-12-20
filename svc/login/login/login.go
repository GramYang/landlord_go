package login

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"landlord_go/basefx"
	"landlord_go/proto"
	hubstatus "landlord_go/svc/hub/status"
	"landlord_go/svc/login/database"
	"strings"
)

const (
	loginSuccess = 200
	loginWrongUserName = 201
	loginWrongPassword = 202
)

func init() {
	proto.Handle_Login_LoginREQ = func(ev cellnet.Event) {

		msg := ev.Message().(*proto.LoginREQ)
		// TODO 第三方请求验证及信息拉取

		ack := &proto.LoginACK{}
		agentSvcID := hubstatus.SelectServiceByLowUserCount("agent", "", false)
		if agentSvcID == "" {
			ack.Result = proto.ResultCode_AgentNotFound

			ev.Session().Send(ack)
			return
		}

		agentWAN := basefx.GetRemoteServiceWANAddress("agent", agentSvcID)


		host, port, err := util.SpliteAddress(agentWAN)
		if err != nil {
			log.Errorf("invalid address: '%s' %s", agentWAN, err.Error())

			ack.Result = proto.ResultCode_AgentAddressError

			ev.Session().Send(ack)
			return
		}

		ack.Server = &proto.ServerInfo{
			IP:   host,
			Port: int32(port),
		}

		ack.GameSvcID = hubstatus.SelectServiceByLowUserCount("game", "", false)

		if ack.GameSvcID == "" {
			ack.Result = proto.ResultCode_GameNotFound

			ev.Session().Send(ack)
			return
		}

		ack.Result = proto.ResultCode_NoError
		var username string
		var password string
		token := string(msg.TokenReq)
		log.Debugf("receive LoginREQ, token: %s\n", token)
		if strings.Contains(token, ":") {
			tmp := strings.Split(token, ":")
			username = tmp[0]
			password = tmp[1]
		}
		userPassword := database.GetUserPassword(username)
		if userPassword == "" {
			ack.TokenAck = loginWrongUserName
			ev.Session().Send(ack)
			return
		}
		if password == userPassword {
			ack.TokenAck = loginSuccess
		} else {
			ack.TokenAck = loginWrongPassword
		}
		ev.Session().Send(ack)
	}
}
