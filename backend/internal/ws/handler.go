package ws

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

type WS struct {
	wsServer *WSServer
	mx       *sync.Mutex
}

func New() *WS {
	return &WS{
		mx: new(sync.Mutex),
	}
}

func (ws *WS) Register(route *gin.RouterGroup) {
	route.GET("/ws", ws.openWS)
}

func (ws *WS) openWS(ctx *gin.Context) {
	ws.mx.Lock()
	defer ws.mx.Unlock()

	// старый коннект закрываем, новый открываем
	if ws.wsServer != nil {
		ws.wsServer.Close()
	}

	ws.wsServer = NewWSServer(context.Background())
	err := ws.wsServer.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(fmt.Errorf("server error ws: %w", err))
		return
	}
}

func (ws *WS) Write(msg string) error {
	ws.mx.Lock()
	defer ws.mx.Unlock()

	if ws.wsServer == nil {
		return nil
	}

	return ws.wsServer.WriteMsg(msg)
}
