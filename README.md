<!-- TOC -->

- [什么是landlord_go](#什么是landlord_go)
- [使用](#使用)
- [框架分析](#框架分析)
    - [basefx](#basefx)
        - [init.go](#initgo)
        - [createpeer.go](#createpeergo)
    - [proto](#proto)
    - [svc.agent](#svcagent)
        - [流程](#流程)
        - [model目录](#model目录)
            - [frontend.go](#frontendgo)
            - [user.go](#usergo)
            - [routetab.go](#routetabgo)
        - [api目录](#api目录)
            - [broadcast.go](#broadcastgo)
            - [adaptor.go](#adaptorgo)
        - [backend目录](#backend目录)
            - [event.go](#eventgo)
            - [hooker.go](#hookergo)
        - [frontend目录](#frontend目录)
            - [bind.go](#bindgo)
            - [hooker.go](#hookergo-1)
            - [listener.go](#listenergo)
            - [packet.go](#packetgo)
            - [transmitter.go](#transmittergo)
        - [heartbeat目录](#heartbeat目录)
            - [heartbeat_msg.go](#heartbeat_msggo)
        - [routerule目录](#routerule目录)
            - [update.go](#updatego)
        - [agent.backend](#agentbackend)
            - [service.SvcEventHooker](#servicesvceventhooker)
            - [broadcasterHooker](#broadcasterhooker)
        - [tcp.frontend](#tcpfrontend)
            - [directTCPTransmitter](#directtcptransmitter)
            - [FrontendEventHooker](#frontendeventhooker)
    - [svc.game](#svcgame)
        - [流程](#流程-1)
        - [json目录](#json目录)
            - [json_msg.go](#json_msggo)
        - [database目录](#database目录)
        - [landlord目录](#landlord目录)
        - [verify目录](#verify目录)
            - [verify_msg.go](#verify_msggo)
        - [svc.backend](#svcbackend)
            - [BackendMsgHooker](#backendmsghooker)
    - [svc.hub](#svchub)
        - [流程](#流程-2)
        - [status目录](#status目录)
            - [api.go](#apigo)
            - [recv.go](#recvgo)
            - [send.go](#sendgo)
        - [api目录](#api目录-1)
            - [hooker.go](#hookergo-2)
            - [pubsub.go](#pubsubgo)
        - [subscribe目录](#subscribe目录)
            - [subcribe_msg.go：](#subcribe_msggo)
        - [model目录](#model目录-1)
            - [status.go](#statusgo)
            - [channel.go](#channelgo)
    - [svc.login](#svclogin)
        - [流程](#流程-3)
        - [login目录](#login目录)
            - [login.go](#logingo)
        - [tcp.client](#tcpclient)
        - [tcp.hub](#tcphub)

<!-- /TOC -->

# 什么是landlord_go

根据[cellmesh_demo](https://github.com/davyxu/cellmesh_demo)写成的斗地主游戏后端，对接Android前端landlord_client

# 使用

先开启[cellmesh](https://github.com/davyxu/cellmesh)，然后依次开启landlord_go/svc/下的服务hub、login、game、agent

# 框架分析

> 这实际上是[cellmesh_demo](https://github.com/davyxu/cellmesh_demo)的框架分析，有一些修改和简化。其中涉及到[cellnet](https://github.com/davyxu/cellnet)和[cellmesh](https://github.com/davyxu/cellmesh)的知识点请自行补习。

## basefx

### init.go

* init()

将proto.PingACK和proto.SvcStatusACK存入msglog的黑名单中
创建并开启fxmodel.Queue
调用service.Init(procName)保存procName，解析命令行参数
调用service.ConnectDiscovery()连接memsd服务

### createpeer.go

* CreateCommnicateAcceptor()

引用fxmodel.Queue，根据param创建peer
根据proto.GetMessageHandler(param.SvcName)获取usercallback，应用到peer上
设置socket的读写buffer
将peer加入fxmodel，开启peer，然后在memsd上注册peer

* CreateCommnicateConnector()

proto.GetMessageHandler(service.GetProcName())
获取usercallback
opt的maxcount默认是0，即没有效果。rules默认是单个切片，值为”dev”
引用fxmodel.Queue准备作为peer的消息队列
调用mp := service.DiscoveryService在服务发现上获取服务信息后创建connector连接此服务

* GetRemoteServiceWANAddress

result := service.QueryService(svcName,service.Filter_MatchSvcID(svcid))
从memsd根据svcName查询一堆desc后，在从中挑出svcid和desc.ID相等的desc，返回其地址

## proto

landlord.proto为前端后端通信的proto。msgbind.go中对消息和回调函数进行了注册。

## svc.agent

### 流程

1. 连接服务发现，从服务发现下载路由规则。所谓的路由规则就是哪些消息走哪些服务。所谓的路由规则在服务发现上的地址为config_demo/route_rule，值为RouteTable的json字符串。(这里我修改了一下，直接写死到代码里面去了)

2. 开始心跳检查，即遍历agent持有的与所有client的连接session，检查session中user的u.LastPingTime是否超时。检查的评率和超时时间从服务发现上面取。为"config/agent/heartbeat_sec"+5。（这里我也写死了）

3. 赋值agentmodel.AgentSvcID，这一值将被使用在tcp.frontend协议，该协议用于agent作为acceptor来处理client的连接。

4. 创建connector来连接game，协议为game.backend。这里不能使用usercallback。

5. 创建acceptor来应对client连接，协议为tcp.frontend。这里的usercallback为空。

6. 创建connector来连接hub，协议tcp.hub，回调函数调用hubstatus.StartSendStatus()发送频道名为agent_status的状态。

7. import中调用proto包初始化协议，调用agent/backend包初始化协议

### model目录

#### frontend.go
frontend指代前端，也就是与agent保持长连接的client。其中保存AgentSvcID和FrontendSessionManager，FrontendSessionManager管理多个client与一个agent的session。

1. SessionToUser

func SessionToUser(clientSes cellnet.Session) \*User 
从ses中的peer中的contextset查找user键对应的值，类型为\*User 

#### user.go
其中主要是网关用户User结构体的定义和相关方法。User就是对前端client连接的抽象，后端通过访问User来访问前端client。结构体如下：
```
type Backend struct {
	SvcName string
	SvcID   string // 只保留绑定后台的svcid,即便后台更换session,也无需同步
}

type User struct {
	ClientSession cellnet.Session
	Targets       []*Backend
	LastPingTime  time.Time
	CID proto.ClientID
}
```

1. BroadcastToBackends

广播到这个用户绑定的所有后台
遍历user中的Targets，根据其信息对比连接到服务发现的session，然后向这些session发送msg

2. TransmitToBackend

将client的消息传递到后端的服务
也是在服务发现中获取svcID对应的session，然后发送一个TransmitACK，其中包含要传递的msgID、msgData和client的ID

3. SetBackend

绑定用户后台，就是把后台服务信息写到user里面去
先判断一下Targets里是否有重名的svcName，有的话就将其value替换成svcID 
把svcName和svcID 键值对写入Targets
把隔壁的AgentSvcID写入CID

4. GetBackend

遍历Targets返回svcName 对应的SvcID

#### routetab.go

保存：
```
ruleByMsgName = map[string]*table.RouteRule{}
ruleByMsgID = map[int]*table.RouteRule{}
```
以及一系列操作函数。
而RouteRule就是路由规则，定义：
```
type RouteRule struct {
	MsgName string
	SvcName string
	Mode    string // auth: 需要授权 pass: 可通过
	MsgID   int
}
```
这里的RouteRule虽然没有用到过mode，但是可以代表一个消息-服务的映射关系：
一个消息对应一个服务，一个消息可以发送到这个服务中去。

### api目录

#### broadcast.go
1. BroadcastAll

将msg加密成data和meta，然后调用service.VisitRemoteService遍历所有ses中ctx对应值的Name是agent的ses，向其中发送一个TransmitACK，包裹meta.ID和data。

注意，这里访问的是remotesrv.go里保存的session，选取的是所有订阅了agent这个服务的session（这里game也只持有了与agent的连接）。而remotesrv.go里的服务是在SvcEventHooker中被添加的，涉及到
```
tcp.svc
svc.backend
agent.backend
```
这里的service.VisitRemoteService只是把消息发给了agent，由All这个标志来指示群发给所有连接agent的client

2. Send

从cid中获取svcID，从服务发现中获取agentSes。向agentSes中发送
```
TransmitACK：
		agentSes.Send(&proto.TransmitACK{
			MsgID:    uint32(meta.ID),
			MsgData:  data,
			ClientID: cid.ID,
		})
```

#### adaptor.go
就一个函数
```
func HandleBackendMessage(userHandler func(ev cellnet.Event, cid proto.ClientID)) func(ev cellnet.Event) 
```
该函数只在game的两个目录中被调用，用于proc的callback集中处理
简单来说，就是统一提取了cellnet.Event中的clientID和svcID
传入的ev的类型为*backend.RecvMsgEvent时，用ev中的ClientID、ev中的session的contextset中ctx的值中的svcID来填充ClientID，然后调用userHandler。

### backend目录
#### event.go
```
type RecvMsgEvent struct {
	Ses      cellnet.Session
	Msg      interface{}
	ClientID int64 //比常规的RecvMsgEvent多一个这个
}
```
该结构体实现了cellnet.Event接口，在svc.backend协议中被使用，是game用来连接服务发现的协议，用来对接TransmitACK
其专门用来对接TransmitACK，以存放其中的ClientID

#### hooker.go

注册协议svc.backend和agent.backend，定义了BackendMsgHooker网关消息处理和broadcasterHooker，svc.backend中用到了BackendMsgHooker，agent.backend中用到了broadcasterHooker

这里可以看到，svc.backend是后台服务使用的协议（game使用了svc.backend），用于处理agent发送来的消息。
而agent.backend是被agent使用，用于处理后台服务发送来的消息。其他服务调用agent的api来向agent发送消息。


### frontend目录

#### bind.go
   
就一个函数
```
// 将客户端连接绑定到后台服务
func bindClientToBackend(backendSvcID string, clientSesID int64) (*model.User, error) 
```
就是把session的contextset中的ctx对应的值写到user对应的值里去。也就是把后端服务的session中ctx对应的信息，填写到前端client的user对应的User指针里，完成将客户端连接绑定到后端服务。
这里在user中存入svcname，svcid。只不过这两个信息都是在agent持有的后台服务的session中获取的。这里实际上是变相确认了agent是否持有前端要求绑定的后台服务！
这里注意：如果agentmodel.GetClientSession(clientSesID)不能根据id获取session，那就没办法了。。。要报错。


#### hooker.go
注册了tcp.frontend协议，定义FrontendEventHooker，tcp.frontend使用了FrontendEventHooker。FrontendEventHooker用于处理agent内部抛出的SessionClosed，将其广播给后台服务。
tcp.frontend用于创建tcpAcceptor来对接client的连接
其directTCPTransmitter会将client发来相关消息转发给后端服务或者是直接返回给client（这一点和一般的协议完全不同）
其FrontendEventHooker会在处理内部消息，处理的方式也是单点转发或者广播给后端服务
而agent在建立tcp.frontend的tcpAcceptor时根本不需要用户回调！！

1. ProcFrontendPacket

该函数在directTCPTransmitter中被使用，而directTCPTransmitter在tcp.frontend中被使用。
将client发来相关消息转发给后端服务或者是直接返回给client

#### listener.go
两个方法，agent用来开启面向client的tcpAcceptor。
1. Start

创建一个peer，绑定协议，但是用户回调函数是空的
将peer的sessionmanager赋给model.FrontendSessionManager
服务发现注册peer，fxmodel加入peer

2. Stop

关闭上面的peer
服务发现解注册该服务

#### packet.go
就一个函数
```
// 接收Length-Type-Value格式的封包流程
func RecvLTVPacketData(reader io.Reader, maxPacketSize int) (msgID int, msgData []byte, err error) 
```
就是普通的封包，我没搞懂为啥还要专门写一个函数出来

#### transmitter.go
定义directTCPTransmitter结构体，用于tcp.frontend协议

### heartbeat目录
#### heartbeat_msg.go
就一个函数StartCheck()

其中每隔几秒钟就检查FrontendSessionManager中session中contextset的user里的LastPingTime，超时了就关闭session。
这里的session是client与agent的连接。

### routerule目录
#### update.go
从服务发现上下载路由规则

### agent.backend
```
agent/backend/hooker.go
	proc.RegisterProcessor("agent.backend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(
			new(service.SvcEventHooker), // 服务互联处理
			new(broadcasterHooker),      // 网关消息处理
			new(tcp.MsgHooker)))         // tcp基础消息处理
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
```
在tcp.ltv的基础上添加了service.SvcEventHooker和broadcasterHooker。
该协议由agent来作为connector连接后端的服务，比如svc.backend。
该协议是没有usercallback的。

#### service.SvcEventHooker
用作维护connBySvcID
* OnInboundEvent

根据inputEvent.Message().(type)进行处理：
  1. ServiceIdentifyACK
   
调用AddRemoteService将ses存入connBySvcID，其类型为map[string]cellnet.Session{}

agent添加连接的后端服务，比如game

  2. SessionConnected

取出session中的contextset，从中取”sd”对应的值，取成功就调用
```
			// 用Connector的名称（一般是ProcName）让远程知道自己是什么服务，用于网关等需要反向发送消息的标识
			inputEvent.Session().Send(&ServiceIdentifyACK{
				SvcName: GetProcName(),
				SvcID:   GetLocalSvcID(),
			})
```
然后调用AddRemoteService(inputEvent.Session(), sd.ID, sd.Name)添加服务

向session发送一个ServiceIdentifyACK，这与上面的ServiceIdentifyACK对应

3. SessionClosed

调用RemoveRemoteService(inputEvent.Session())删除服务

* OnOutboundEvent

什么都不做

#### broadcasterHooker
用来转发后端服务（如game）传过来的消息给client。分为单发，群发，全部群发。
* OnInboundEvent
  
判断inputEvent.Message().(type)
1. TransmitACK
   
生成一个cellnet.RawPacket实例，单发、群发或者全部发。
单发是根据ClientID，调用

clientSes := model.GetClientSession(incomingMsg.ClientID)

群发是根据ClientIDList，遍历调用

clientSes := model.GetClientSession(cid)

全部发则是看All这个bool值，调动model.FrontendSessionManager.VisitSession遍历sessionmanager所有的session发送。

这里是将后端服务传递的信息转发到前端的client中去，通过后端消息附带的clientid和All标志来指定发给哪些client

* OnOutboundEvent
  
根据inputEvent.Message().(type)判断
1. TransmitACK

就log一下

### tcp.frontend
```
agent/frontend/hooker.go
	// 前端的processor
	proc.RegisterProcessor("tcp.frontend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {
		bundle.SetTransmitter(new(directTCPTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(
			new(tcp.MsgHooker),       //  TCP基础消息及日志
			new(FrontendEventHooker), // 内部消息处理
		))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
```
更换了directTCPTransmitter、添加了FrontendEventHooker
tcp.frontend被agent用作acceptor来对应client的连接

#### directTCPTransmitter
这里是应用路由规则，处理PingACK和VerifyREQ的地方

* OnRecvMessage

从session中获取socketOpt

conn, ok := ses.Raw().(net.Conn)
在ok的情况下，for循环中设置读超时，执行
```
// 接收来自客户端的封包
msgID, msgData, err = RecvLTVPacketData(reader, opt.MaxPacketSize())
msg, err = ProcFrontendPacket(msgID, msgData, ses)
```
与TCPMessageTransmitter相比，directTCPTransmitter的不同：
directTCPTransmitter中在ltv拆包后，先由ProcFrontendPacket来处理拆包后的数据，只有在
ProcFrontendPacket抛出err后才会传递给hooker来进一步处理拆包数据。

那么ProcFrontendPacket什么时候会报错呢？

正常情况下不会报错，也就不会执行到hooker。所以下面的FrontendEventHooker也只是处理SessionAccepted和SessionClosed。

1. ProcFrontendPacket

处理消息：

a.PingACK类型

根据ses获取User实例，如果不为空则调用ses.Send(&proto.PingACK{})回消息。如果为空则调用ses.Close()

b.VerifyREQ类型

u, err := bindClientToBackend(userMsg.GameSvcID, ses.ID())

根据GameSvcID获取backendSes，将backendSes变成contextset，取ctx对应的值*RemoteServiceContext也就是sd并返回。

根据ses.ID()调用GetClientSession获取clientSes，将clientSes变成contextset，去user对应的值User指针并返回。

又调用u = model.CreateUser(clientSes)，根据clientSes生成一个User指针，并覆盖上面的u，这里还把User指针存入contextset的user键中。

u.SetBackend(sd.Name, sd.SvcID)

该方法向u中写入值，主要是其self.CID和self.Targets

根据clientses创建user，其实就是用user包裹clientses，然后再将要绑定的后端服务的name和svcid存入user中。

本来呢，参数backendSvcID是可以直接推导出name和svcid的，但是这里特意根据backendSvcID取到backendses，再从里面获取name和svcid，也是变相的确认了一下agent是否持有该后端服务的连接。

在上面bindClientToBackend成功的情况下，调用
u.TransmitToBackend(userMsg.GameSvcID, msgID, msgData)
从本地根据GameSvcID获取ses，再调用ses的Send方法发送一个TransmitACK结构体，里面包含了msgID和msgData。

绑定client要求的后台服务，顺便转发信息

c.default

rule := model.GetRuleByMsgID(msgID)

返回的rule是一个RouteRule指针

u := model.SessionToUser(ses)

从ses中获取contextset中的user对应值，User指针

err = u.TransmitToBackend(u.GetBackend(rule.SvcName), msgID, msgData)

同上

ProcFrontendPacket中只有在消息未注册和后端服务session为空的时候才会抛出err！才会用到FrontendEventHooker

* OnSendMessage
  
在设置写超时的情况下调用err = util.SendLTVPacket(writer, ses.(cellnet.ContextSet), msg)将msg写入ses。注意：这里传入的ses.(cellnet.ContextSet)其实球用都没用，根本没用到


#### FrontendEventHooker
只处理SessionClosed

* OnInboundEvent
  
判断inputEvent.Message().(type)
1. SessionClosed

u := model.SessionToUser(inputEvent.Session())

从ses的contextset中获取user对应的User指针
```
u不为空的情况下调用
			u.BroadcastToBackends(&proto.ClientClosedACK{
				ID: proto.ClientID{
					ID:    inputEvent.Session().ID(),
					SvcID: model.AgentSvcID,
				},
			})
```
BroadcastToBackends中遍历User的Targets域，Targets是[]*Backend。依次从Backend的SvcID获取本地ses，调用ses的Send将msg发送出去。

在连接关闭时向该session（也就是一个user）的所有订阅的后端服务广播ClientClosedACK消息。


* OnOutboundEvent

什么都不做

## svc.game

### 流程
1. 连接服务发现。

2. 创建一个acceptor，协议svc.backend，面向agent的连接。

3. 连接hub，发送频道game_status的信息。按道理讲，game没有什么状态好发送的，毕竟并发连接是agent来扛的，game就是一个单session连接的计算服务器。

4. import中调用game/json包，这就是游戏后端的主要逻辑了。调用game/verify包用来处理VerifyREQ，返回一个VerifyACK


### json目录
#### json_msg.go

处理proto.JsonREQ消息，调用landlord目录的handler.go的不同方法进行处理

### database目录

持久化方法，使用的是mysql+database/sql。

### landlord目录

游戏的主要逻辑都在这里。具体不表，参考代码注释。

### verify目录
#### verify_msg.go
处理client发送的VerifyREQ，简单的log一下，并反馈一个VerifyACK。

VerifyREQ的作用主要是用来触发

u, err := bindClientToBackend(userMsg.GameSvcID, ses.ID())

根据后端服务的svcid获取后端服务的session（前提是后端服务已经连接了agent），然后生成user，填入后端服务的name和svcid，然后将user存入session。

### svc.backend
```
svc.backend在agent/backend/hooker.go
	proc.RegisterProcessor("svc.backend", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(
			new(service.SvcEventHooker), // 服务互联处理
			new(BackendMsgHooker),       // 网关消息处理
			new(tcp.MsgHooker)))         // tcp基础消息处理
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
```
在tcp.ltv的基础上增加了两个hooker

该协议被game用来作为acceptor来迎接来自agent的agent.backend协议的连接。

#### BackendMsgHooker
* OnInboundEvent

判断inputEvent.Message().(type) 
1. TransmitACK

userMsg, _, err := codec.DecodeMessage(int(incomingMsg.MsgID), incomingMsg.MsgData)

解析TransmitACK里的id和[]byte，返回userMsg
构建RecvMsgEvent实例，该实例实现了cellnet.Event
```
		ev := &RecvMsgEvent{
			Ses:      inputEvent.Session(),
			Msg:      userMsg,
			ClientID: incomingMsg.ClientID,
		}
```
专门对应TransmitACK，将TransmitACK拆包后用RecvMsgEvent包裹

2. default

什么都不做


* OnOutboundEvent

就log一下

## svc.hub
### 流程
1. 连接服务发现。

2. 创建一个acceptor，协议为tcp.svc。这个tcp.svc主要是用来维护connBySvcID的。

3. import中运行了hub/subscribe目录的init函数。其中


### status目录

#### api.go
SelectServiceByLowUserCount

1. model.VisitStatus中传入的回调函数：

name, _, group, err := service.ParseSvcID(status.SvcID)

就是反向解析，看看方法签名就知道了：

func ParseSvcID(svcid string) (svcName string, svcIndex int, svcGroup string, err error) 

如果name和svcName不一样就return true，这说明需要name和svcName相同。

在group为空或者等于svcGroup ，且service.GetRemoteService(status.SvcID)返回不为空，将statusBySvcID 的status添加进statusList。

2. 检查statusList的长度：

0：
说明没找到，返回空。

1：
返回statusList[0].SvcID

default：
好像是随机选一个

从hub的model的statusBySvcID 中选出复合要求的status，返回其中的svcID。要求：人数比较少

#### recv.go
就一个函数：

func StartRecvStatus(channelNames []string, svcStatusHandler *func(ev cellnet.Event))

channelNames为[]string{"game_status", "agent_status"}

svcStatusHandler为&proto.Handle_Login_SvcStatusACK

1. 遍历channelNames，调用hubapi.Subscribe向hub发送SubscribeChannelREQ。SubscribeChannelREQ的响应在callback中，但是callback中空函数的赋值在hub/subscribe/subscribe_msg.go中


2. 给svcStatusHandler赋值

处理SvcStatusACK

model.UpdateStatus更新statusBySrvID中的Status
svcStatusHandler用来处理发送给login的SvcStatusACK消息的。SvcStatusACK由hub发送给login，这样实现了hub向login推送消息更新相关信息。

然而，hub是不产生SvcStatusACK的， 这个SvcStatusACK是由其他服务发送给hub，然后由hub转发而来的。

3. 在一个loop中每3秒执行一个回调

回调中检查超过3秒没有更新的status，遍历调用model.RemoveStatus删除

#### send.go

只有一个函数

func StartSendStatus(channelName string, updateInterval time.Duration, statusGetter func() int) 

channelName为"agent_status"

updateInterval 是3秒

statusGetter返回的是agent连接的client的session数量

1. 在一个loop中每3秒执行一个回调

回调中构建一个SvcStatusACK，调用hubapi.Publish发送出去
这里才是产生SvcStatusACK的地方！！

### api目录
#### hooker.go
注册tcp.hub协议，定义subscriberHooker

#### pubsub.go
1. ConnectToHub

设置model.OnHubReady，然后再以tcp.hub的协议连接hub，在tcp.hub的hooker中在响应连接hub成功时会调用model.OnHubReady。
同时也会调用Subscribe()订阅本地的procname和svcID。譬如是login调用ConnectToHub，那么Subcribe订阅的就是login和loginSvcID

2. Subscribe
调用
```
	model.HubSession.Send(&proto.SubscribeChannelREQ{
		Channel: channel,
	})
```
向hub发送SubscribeChannelREQ

3. Publish

relay.Relay(model.HubSession, msg, channel)

向hub发送一个RelayACK，msg会被压缩成[]byte并包裹在RelayACK中，当然还会有对应的msgID

### subscribe目录
#### subcribe_msg.go：
* proto.Handle_Hub_SubscribeChannelREQ
  
model.AddSubscriber(msg.Channel, ev.Session())

将channel和session存入hub/model包的chanByName = map[string][]cellnet.Session{}

再向session返回一个SubscribeChannelACK。貌似对SubscribeChannelACK没有处理
这里是唯一的AddSubscriber被调用的地方

* proto.Handle_Hub_Default

只响应SessionClosed，调用model.RemoveSubscriber(ev.Session(), func(chanName string) 

* relay.SetBroadcaster(callback)

callback中先从even中抽取订阅的channelName，然后调用VisitSubscriber遍历chanByName，再调用relay.Relay(ses, event.Message(), channelName)给发送出去。

这里是实现向hub订阅消息获取反馈的关键：

就像是一面镜子，hub维护了

chanByName = map[string][]cellnet.Session{}

这样一个map，key为订阅的内容，value为订阅此内容的session切片。你向hub发送了相关内容的消息，这个消息会被转发到所有订阅此内容的session。
这是使用cellnet的relay实现的。

### model目录

#### status.go
所有连接hub的服务的状态，（他们的登录人数），3秒一更新
```
statusBySvcID = map[string]*Status{}
其中的Status的定义：
type Status struct {
	UserCount int32
	SvcID string
	LastUpdate time.Time
}
```

#### channel.go
保存所有订阅channelName的Session，channelName就是key

chanByName = map[string][]cellnet.Session{}

以及相关操作方法

```
tcp.svc
其定义在cellmesh/service/hooker.go
	proc.RegisterProcessor("tcp.svc", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(SvcEventHooker), new(tcp.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
```
可以看出与cellnet的tcp.ltv相比，就是hooker添加了一个SvcEventHooker在前面。

tcp.svc主要提供了一个服务端客户端双向添加连接服务的功能。

tcp.svc在hub中的acceptor被使用，而其他服务连接hub的协议则是tcp.hub

SvcEventHooker的两个方法：

* OnInboundEvent

处理inputEvent.Message()：
1. ServiceIdentifyACK

pre := GetRemoteService(msg.SvcID)

GetRemoteService是从本地的connBySvcID地图里取出cellnet.Session，connBySvcID应该是服务发现保存连接的地图。
调用GetRemoteService取cellnet.Session，如果为空就调用AddRemoteService添加进去。

针对ServiceIdentifyACK添加连接服务

2. SessionConnected

ctx := inputEvent.Session().Peer().(cellnet.ContextSet)
ctx.FetchContext("sd", &sd) 

这里的FetchContext的作用是根据key从ctx中取值，取到值就将值转换为\*sd的类型，并把其赋给\*sd。

ctx中的值是在multi.AddPeer时被存进去的。

从ctx中取值成功后，调用inputEvent.Session().Send发送一个
```
&ServiceIdentifyACK{
				SvcName: GetProcName(),
				SvcID:   GetLocalSvcID(),
			}
```
然后调用AddRemoteService(inputEvent.Session(), sd.ID, sd.Name)

3. SessionClosed

RemoveRemoteService(inputEvent.Session())

很明显

SessionConnected在连接后向对方发送一个ServiceIdentifyACK，包裹你自己的svcname和svcid。然后添加连接服务。
SessionClosed则删除连接服务。

* OnOutboundEvent

什么都不做。

## svc.login
### 流程
1. 连接服务发现。

2. 创建一个acceptor，协议用tcp.client。

3. 用tcp.hub协议创建connector连接hub，hub的地址从服务发现上面来找。连接hub后，开始接收game_status和agent_status信息，拿了信息什么都没有做。

4. import中调用login/login包来处理前端的登陆请求。

### login目录
#### login.go
init函数中给proto.Handle_Login_LoginREQ函数赋值：
1. agentSvcID := hubstatus.SelectServiceByLowUserCount("agent", "", false)

返回人数较少的agent的agentSvcID。这里的agentSvcID是hub的status中的。
如果agentSvcID为空，返回一个LoginACK返回ResultCode_AgentNotFound。

2. agentWAN := basefx.GetRemoteServiceWANAddress("agent", agentSvcID)

根据服务名称从服务发现中返回ServiceDesc，从中获取WANAddress的对应值，WANAddress就是外网地址。如果没有WANAddress就返回ServiceDesc的ip+port，就是内网地址。

3. 构建ServerInfo

其中的ip和port是根据上面的agentWAN 解析而来，将其赋给LoginACK

4. ack.GameSvcID = hubstatus.SelectServiceByLowUserCount("game", "", false)

这次返回的是game的svcID，赋给LoginACK

5. service.Reply(ev, &ack)返回

首先通过hubstatus返回人数最少的agentSvcID。注意：这个时候client还没有连接agent，只是先行连接login。

根据服务名和svcID从服务发现中获取其地址（可能是外网地址，也可能是内网地址），同样从hubstatus返回人数少的gamesvcID，然后返回。

### tcp.client
```
其定义在cellmesh/service/hooker.go
	proc.RegisterProcessor("tcp.client", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(tcp.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
```
和tcp.ltv一样的

### tcp.hub
```
其定义在hub/api/hooker.go中
	proc.RegisterProcessor("tcp.hub", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(proc.NewMultiHooker(new(subscriberHooker), new(tcp.MsgHooker)))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
```
tcp.hub相比于tcp.ltv，多了个subscriberHooker。

subscriberHooker的方法：
* OnInboundEvent

只处理inputEvent.Message()为SessionConnected：
```
model.HubSession = inputEvent.Session()
Subscribe(service.GetProcName())
Subscribe(service.GetLocalSvcID())
两个Subscribe函数会调用
	model.HubSession.Send(&proto.SubscribeChannelREQ{
		Channel: channel,
	})
```
调用cellnet.QueuedCall(fxmodel.Queue, model.OnHubReady)

SubscribeChannelREQ会在hub服务的usercallback中响应，调用
hubmodel.AddSubscriber(msg.Channel, ev.Session())
然后在connector的消息队列中调用 model.OnHubReady

* OnOutboundEvent
什么都不做
