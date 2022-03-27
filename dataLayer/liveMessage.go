package dataLayer

import (
	"code/Hahachitchat/definition"
	"github.com/gorilla/websocket"
)

type registrant struct {
	uId    uint64
	wsConn *websocket.Conn
}
type notification struct {
	targetUid   uint64                 `json:"target_uid"`
	messageType definition.MessageType `json:"message_type"`
}

type notificationHub struct {
	registerChan     chan registrant
	notificationChan chan notification

	register map[uint64][]*websocket.Conn
}

var hub notificationHub

func RunNotificationHub() {
	// 初始化在线消息中心
	hub.registerChan = make(chan registrant, 50)
	hub.notificationChan = make(chan notification, 1000)
	hub.register = make(map[uint64][]*websocket.Conn)

	// Run 在线消息中心
	for {
		select {
		case register := <-hub.registerChan:
			hub.register[register.uId] = append(hub.register[register.uId], register.wsConn)
		case noti := <-hub.notificationChan:
			for i, wsConn := range hub.register[noti.targetUid] {
				if err := wsConn.WriteJSON(notification{
					targetUid:   noti.targetUid,
					messageType: noti.messageType,
				}); err != nil {
					Serverlog.Println("[sendNotification]", err)
					// 自动断掉连接，退出 hub
					wsConn.Close()
					hub.register[noti.targetUid] = append(hub.register[noti.targetUid][:i], hub.register[noti.targetUid][i+1:]...)
				}
			}
		}
	}
}

func RegisterNotificationHub(uId uint64, ws *websocket.Conn) {
	if ws != nil {
		hub.registerChan <- registrant{
			uId:    uId,
			wsConn: ws,
		}
	}
}

func Notify(targetUid uint64, messageType definition.MessageType) {
	hub.notificationChan <- notification{
		targetUid:   targetUid,
		messageType: messageType,
	}
}
