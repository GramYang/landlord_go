package landlord

import (
	"landlord_go/proto"
	"landlord_go/svc/agent/api"
	"math"
	"math/rand"
	"time"
)

//根据牌桌号获取座位号
func GetSeatNum(tableNum, tablePlayerCount int32) int32 {
	return (tableNum - 1) * 3 + tablePlayerCount + 1
}

//获取随机牌，牌洗了三遍
func GetRandomCards() []int32 {
	source := make([]int32, len(CARDS))
	copy(source, CARDS)
	rand.Seed(time.Now().UnixNano())
	for i:=0; i<3; i++ {
		rand.Shuffle(len(source), func(a, b int) {
			source[a], source[b] = source[b], source[a]
		})
	}
	return source
}

//随机选地主
func GetRandomLandlord(tableNum int32) int32{
	return (tableNum - 1) * 3 + 1 + int32(rand.Float32() * 2.0)
}

func GetLeftRivalSeatNum(yourSeatNum int32, players []*Player) int32 {
	list := make([]int32, 3)
	for _, v := range players {
		list = append(list, v.SeatNum)
	}
	maxSeatNum := max(list)
	if maxSeatNum - yourSeatNum == 0{
		return maxSeatNum - 1
	} else if maxSeatNum - yourSeatNum == 1 {
		return maxSeatNum - 2
	} else {
		return maxSeatNum
	}
}

func GetRightRivalSeatNum(yourSeatNum int32, players []*Player) int32 {
	list := make([]int32, 3)
	for _, v := range players {
		list = append(list, v.SeatNum)
	}
	maxSeatNum := max(list)
	if maxSeatNum - yourSeatNum == 0{
		return maxSeatNum - 2
	} else if maxSeatNum - yourSeatNum == 1 {
		return maxSeatNum
	} else {
		return maxSeatNum - 1
	}
}

//返回table中的座位号-用户名map
func RefreshSeatNum2UserName(table *Table) map[int32]string {
	tablePlayers := make(map[int32]string)
	for _, v := range table.Players {
		if v != nil {
			tablePlayers[v.SeatNum] = v.UserName
		}
	}
	return tablePlayers
}

func max(a []int32) int32 {
	var m int32
	m = math.MinInt32
	for _, v := range a {
		if v > m {
			m = v
		}
	}
	return m
}

func min(a []int32) int32 {
	var m int32
	m = math.MaxInt32
	for _, v := range a {
		if v < m {
			m = v
		}
	}
	return m
}

//群发的封装
func WrapMultiSend(players []*Player, ack *proto.JsonACK, cid *proto.ClientID) {
	var cids []*proto.ClientID
	for _, v := range players {
		if cid != nil && v.Cid.ID == cid.ID {
			continue
		}
		if v != nil {
			cids = append(cids, v.Cid)
		}
	}
	agentapi.MultipleSend(cids, ack)
}