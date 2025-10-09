package ws

import (
	"context"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type WSServer struct {
	upgrader *websocket.Upgrader
	conn     *websocket.Conn
	ctx      context.Context
	close    context.CancelFunc
	fail     atomic.Int32
	lastMsg  time.Time
	mx       sync.RWMutex
}

const (
	failCount = 2
)

func NewWSServer(pctx context.Context) *WSServer {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ctx, cancel := context.WithCancel(pctx)
	srv := &WSServer{
		upgrader: upgrader,
		ctx:      ctx,
		close:    cancel,
	}

	return srv
}

func (ws *WSServer) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) error {
	var err error

	ws.conn, err = ws.upgrader.Upgrade(w, r, responseHeader)
	if err == nil {
		go ws.readMsg()
	}

	return err
}

func (ws *WSServer) Close() {
	ws.conn.Close()
	ws.close()
}

func (ws *WSServer) readMsg() {
	msgChan := make(chan struct{})
	defer close(msgChan)

	// контроль сообщений
	go func(msg <-chan struct{}) {
		for {
			select {
			case <-ws.ctx.Done():
				log.Println("WS: context done")
				return
			case <-time.Tick(time.Second * 10):
				// если не получаем сообщения 10 секунд, закрываем конект (фронт переоткроет). Пинги должны приходить раз в 5 сек.

				ws.mx.RLock()
				if time.Since(ws.lastMsg).Seconds() >= 10 {
					if count := ws.fail.Add(1); count >= failCount {
						ws.Close()
					}
				}
				ws.mx.RUnlock()
			case <-msg:
				ws.mx.Lock()
				ws.lastMsg = time.Now()
				ws.mx.Unlock()
			}
		}
	}(msgChan)

	for {
		select {
		case <-ws.ctx.Done():
			log.Println("WS: context done")
			return
		default:
		}

		msgType, _, err := ws.conn.ReadMessage()
		if err != nil {
			log.Println(errors.Wrap(err, "ws read error"))
		}

		if msgType == -1 { // разрыв сокета
			log.Println("websocket closed")
			break
		}

		msgChan <- struct{}{}
		//log.Println(string(msg))
		//s.subscriber.SDKNotify(msg)
	}

}

func (ws *WSServer) WriteMsg(data string) error {
	if ws.conn == nil {
		return errors.New("ws connection not initialized")
	}

	err := ws.conn.WriteMessage(websocket.TextMessage, []byte(data))
	return err
}

func (ws *WSServer) WriteByteMsg(data []byte) error {
	if ws.conn == nil {
		return errors.New("ws connection not initialized")
	}

	err := ws.conn.WriteMessage(websocket.BinaryMessage, data)
	return err
}
