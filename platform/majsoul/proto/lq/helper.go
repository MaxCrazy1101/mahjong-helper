package lq

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"reflect"
	"strings"
	"time"
)

var (
	lobbyClientMethodMap    = map[string]reflect.Type{}
	fastTestClientMethodMap = map[string]reflect.Type{}
)

func init() {
	t := reflect.TypeOf((*LobbyClient)(nil)).Elem()
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		lobbyClientMethodMap[method.Name] = method.Type
	}

	t = reflect.TypeOf((*FastTestClient)(nil)).Elem()
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fastTestClientMethodMap[method.Name] = method.Type
	}
}

func FindMethod(clientName string, methodName string) reflect.Type {
	methodName = strings.Title(methodName)
	if clientName == "Lobby" {
		return lobbyClientMethodMap[methodName]
	} else { // clientName == "FastTest"
		return fastTestClientMethodMap[methodName]
	}
}

// 下面补充一些功能

func (m *Friend) CLIString() string {
	return fmt.Sprintf("%9d   %s   %s   %s",
		m.Base.AccountId,
		time.Unix(int64(m.State.LoginTime), 0).Format("2006-01-02 15:04:05"),
		time.Unix(int64(m.State.LogoutTime), 0).Format("2006-01-02 15:04:05"),
		m.Base.Nickname,
	)
}

type FriendList []*Friend

func (l FriendList) String() string {
	out := "好友账号ID   好友上次登录时间        好友上次登出时间       好友昵称\n"
	for _, friend := range l {
		out += friend.CLIString() + "\n"
	}
	return out
}

func (m *ActionPrototype) ParseData() (proto.Message, error) {
	// 构造消息类型的全名，这里假设您的包名是 "lq"
	name := "lq." + m.Name

	// 查找消息类型
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(name))
	if err != nil {
		return nil, fmt.Errorf("ActionPrototype.ParseData 未找到类型 %s，请检查！", name)
	}

	// 创建消息类型的实例
	messagePtr := mt.New()
	if err := proto.Unmarshal(m.Data, messagePtr.Interface().(proto.Message)); err != nil {
		return nil, err
	}
	// 返回反序列化后的消息实例
	return messagePtr.Interface().(proto.Message), nil
}
