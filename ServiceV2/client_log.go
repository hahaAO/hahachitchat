package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func InitClientLog() {
	definition.ClientLogChan = make(chan string, 500)
	definition.WsLogConnChan = make(chan struct{}, 1)
}

func LogWebSocketConnect(c *gin.Context) {
	fmt.Println("来了连接")
	select {
	case definition.WsLogConnChan <- struct{}{}:
	default: // 只允许一个连接
		fmt.Println("连接被占用")
		c.JSON(http.StatusOK, definition.CommonResponse{
			State:        definition.ServerError,
			StateMessage: "服务器繁忙",
		})
		return
	}

	defer func() {
		fmt.Println("连接释放1")
		<-definition.WsLogConnChan // 释放连接
	}()

	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	defer func() {
		ws.Close()
		fmt.Println("连接释放2")
	}()

	if err != nil {
		dataLayer.Serverlog.Println("[WebSocketConnect] err: ", err)
		return
	}

	for {
		select {
		case log := <-definition.ClientLogChan:
			ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := ws.WriteMessage(websocket.TextMessage, []byte(log))
			if err != nil {
				dataLayer.Serverlog.Println("[LogWebSocketConnect] err: ", err)
				return
			}
		default:
			ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
			err := ws.WriteMessage(websocket.TextMessage, []byte("PING"))
			if err != nil {
				fmt.Println("发ping失败")
				dataLayer.Serverlog.Println("[LogWebSocketConnect] err: ", err)
				return
			}
			fmt.Println("发ping成功")
			ws.SetReadDeadline(time.Now().Add(5 * time.Second))
			_, message, err := ws.ReadMessage()
			if err != nil || string(message) != "PONG" {
				fmt.Println("收pong失败")
				dataLayer.Serverlog.Println("[LogWebSocketConnect] err: ", err)
				return
			}
			fmt.Println("收pong成功")
			time.Sleep(2 * time.Second)
		}

	}
}
