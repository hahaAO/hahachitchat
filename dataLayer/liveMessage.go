package dataLayer

import (
	"code/Hahachitchat/definition"
	"github.com/gorilla/websocket"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 10 * time.Second
	pingPeriod = 3 * time.Second
)

type registrant struct {
	uId    uint64
	wsConn *websocket.Conn
}
type notification struct {
	TargetUid   uint64                 `json:"target_uid"`
	MessageType definition.MessageType `json:"message_type"`
}

type notificationHub struct {
	registerChan     chan registrant
	unregisterChan   chan registrant
	pingChan         chan registrant
	notificationChan chan notification

	register map[uint64][]*websocket.Conn
}

var hub notificationHub

func RunNotificationHub() {
	// 初始化在线消息中心
	hub.registerChan = make(chan registrant, 50)
	hub.unregisterChan = make(chan registrant, 50)
	hub.pingChan = make(chan registrant, 50)
	hub.notificationChan = make(chan notification, 1000)

	hub.register = make(map[uint64][]*websocket.Conn)

	// Run 在线消息中心
	for {
		select {
		case register := <-hub.registerChan: // 注册
			hub.register[register.uId] = append(hub.register[register.uId], register.wsConn)

			register.wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			register.wsConn.SetPongHandler(func(appData string) error { register.wsConn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

			stop := make(chan struct{})
			go func() { // 发ping
				tick := time.Tick(pingPeriod)
				for range tick {
					select {
					case <-stop:
						break
					default:
						Ping(register.uId, register.wsConn)
					}
				}
			}()

			go func() { // 收pong
				for {
					_, _, err := register.wsConn.ReadMessage()
					if err != nil {
						close(stop)
						break
					}
				}
				UnregisterNotificationHub(register.uId, register.wsConn)
			}()
		case register := <-hub.unregisterChan: // 断连
			for i, conn := range hub.register[register.uId] {
				if conn == register.wsConn {
					hub.register[register.uId] = append(hub.register[register.uId][:i], hub.register[register.uId][i+1:]...)
					register.wsConn.Close()
				}
			}
		case register := <-hub.pingChan: // 发送ping
			register.wsConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := register.wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				Serverlog.Println("[sendPing]", err)
				UnregisterNotificationHub(register.uId, register.wsConn)
			}
		case noti := <-hub.notificationChan: // 发送提示
			for _, wsConn := range hub.register[noti.TargetUid] {
				if err := wsConn.WriteJSON(notification{
					TargetUid:   noti.TargetUid,
					MessageType: noti.MessageType,
				}); err != nil { // 发送失败断开连接
					Serverlog.Println("[sendNotification]", err)
					UnregisterNotificationHub(noti.TargetUid, wsConn)
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

func UnregisterNotificationHub(uId uint64, ws *websocket.Conn) {
	if ws != nil {
		hub.unregisterChan <- registrant{
			uId:    uId,
			wsConn: ws,
		}
	}
}

func Ping(uId uint64, ws *websocket.Conn) {
	if ws != nil {
		hub.pingChan <- registrant{
			uId:    uId,
			wsConn: ws,
		}
	}
}

func Notify(targetUid uint64, messageType definition.MessageType) {
	hub.notificationChan <- notification{
		TargetUid:   targetUid,
		MessageType: messageType,
	}
}
