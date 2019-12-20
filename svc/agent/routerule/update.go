package routerule

import (
	"github.com/davyxu/cellmesh/discovery"
	"landlord_go/svc/agent/model"
	"landlord_go/table"
)

// 用Consul KV下载路由规则
func Download() error {

	log.Debugf("Download route rule from discovery...")

	var tab table.RouteTable

	err := discovery.Default.GetValue(model.ConfigPath, &tab)
	if err != nil {
		return err
	}

	model.ClearRule()

	for _, r := range tab.Rule {
		model.AddRouteRule(r)
	}

	log.Debugf("Total %d rules added", len(tab.Rule))

	return nil
}

//手动添加路由规则
func GetRouteRule() {
	log.Debugln("get route rule...")
	r1 := &table.RouteRule{MsgName:"JsonREQ",SvcName:"game",Mode:"auth",MsgID:20000}
	model.AddRouteRule(r1)
	r2 := &table.RouteRule{MsgName:"VerifyREQ",SvcName:"game",Mode:"pass",MsgID:13457}
	model.AddRouteRule(r2)
}