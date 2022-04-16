package ServiceV2

import (
	"code/Hahachitchat/dataLayer"
	"code/Hahachitchat/definition"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func InitClientLog() {
	definition.ClientLogChan = make(chan string, 500)
	definition.WsLogConnChan = make(chan struct{}, 1)
}

func LogWebSocketConnect(c *gin.Context) {
	select {
	case definition.WsLogConnChan <- struct{}{}:
	default: // 只允许一个连接
		c.JSON(http.StatusOK, definition.CommonResponse{
			State:        definition.ServerError,
			StateMessage: "服务器繁忙",
		})
		return
	}
	defer func() {
		<-definition.WsLogConnChan // 释放连接
	}()

	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	defer ws.Close()
	if err != nil {
		dataLayer.Serverlog.Println("[WebSocketConnect] err: ", err)
		return
	}

	for {
		log := <-definition.ClientLogChan
		err := ws.WriteMessage(websocket.TextMessage, []byte(log))
		if err != nil {
			dataLayer.Serverlog.Println("[LogWebSocketConnect] err: ", err)
			break
		}
	}
}
