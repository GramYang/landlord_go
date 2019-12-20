package proto

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"reflect"
)

// agent
var (
	Handle_Agent_CloseClientACK = func(ev cellnet.Event) { panic("'CloseClientACK' not handled") }
	Handle_Agent_Default        func(ev cellnet.Event)
)

// game
var (
	Handle_Game_JsonREQ   = func(ev cellnet.Event) { panic("'ChatREQ' not handled") }
	Handle_Game_VerifyREQ = func(ev cellnet.Event) { panic("'VerifyREQ' not handled") }
	Handle_Game_Default   func(ev cellnet.Event)
)

// hub
var (
	Handle_Hub_SubscribeChannelREQ = func(ev cellnet.Event) { panic("'SubscribeChannelREQ' not handled") }
	Handle_Hub_Default             func(ev cellnet.Event)
)

// login
var (
	Handle_Login_LoginREQ     = func(ev cellnet.Event) { panic("'LoginREQ' not handled") }
	Handle_Login_SvcStatusACK = func(ev cellnet.Event) { panic("'SvcStatusACK' not handled") }
	Handle_Login_Default      func(ev cellnet.Event)
)

// match
var (
	Handle_Match_SvcStatusACK = func(ev cellnet.Event) { panic("'SvcStatusACK' not handled") }
	Handle_Match_Default      func(ev cellnet.Event)
)

func GetMessageHandler(svcName string) cellnet.EventCallback {

	switch svcName {
	case "agent":
		return func(ev cellnet.Event) {
			switch ev.Message().(type) {
			case *CloseClientACK:
				Handle_Agent_CloseClientACK(ev)
			default:
				if Handle_Agent_Default != nil {
					Handle_Agent_Default(ev)
				}
			}
		}
	case "game":
		return func(ev cellnet.Event) {
			switch ev.Message().(type) {
			case *JsonREQ:
				Handle_Game_JsonREQ(ev)
			case *VerifyREQ:
				Handle_Game_VerifyREQ(ev)
			default:
				if Handle_Game_Default != nil {
					Handle_Game_Default(ev)
				}
			}
		}
	case "hub":
		return func(ev cellnet.Event) {
			switch ev.Message().(type) {
			case *SubscribeChannelREQ:
				Handle_Hub_SubscribeChannelREQ(ev)
			default:
				if Handle_Hub_Default != nil {
					Handle_Hub_Default(ev)
				}
			}
		}
	case "login":
		return func(ev cellnet.Event) {
			switch ev.Message().(type) {
			case *LoginREQ:
				Handle_Login_LoginREQ(ev)
			case *SvcStatusACK:
				Handle_Login_SvcStatusACK(ev)
			default:
				if Handle_Login_Default != nil {
					Handle_Login_Default(ev)
				}
			}
		}
	case "match":
		return func(ev cellnet.Event) {
			switch ev.Message().(type) {
			case *SvcStatusACK:
				Handle_Match_SvcStatusACK(ev)
			default:
				if Handle_Match_Default != nil {
					Handle_Match_Default(ev)
				}
			}
		}
	}

	return nil
}

func init() {

	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*PingACK)(nil)).Elem(),
		ID:    16241,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*LoginREQ)(nil)).Elem(),
		ID:    18888,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*LoginACK)(nil)).Elem(),
		ID:    18889,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*VerifyREQ)(nil)).Elem(),
		ID:    13457,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*VerifyACK)(nil)).Elem(),
		ID:    13458,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*JsonREQ)(nil)).Elem(),
		ID:    20000,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*JsonACK)(nil)).Elem(),
		ID:    20001,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*TestACK)(nil)).Elem(),
		ID:    9315,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*CloseClientACK)(nil)).Elem(),
		ID:    58040,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*ClientClosedACK)(nil)).Elem(),
		ID:    50844,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*TransmitACK)(nil)).Elem(),
		ID:    9941,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*SubscribeChannelREQ)(nil)).Elem(),
		ID:    27927,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*SubscribeChannelACK)(nil)).Elem(),
		ID:    55294,
	})
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("gogopb"),
		Type:  reflect.TypeOf((*SvcStatusACK)(nil)).Elem(),
		ID:    50227,
	})
}