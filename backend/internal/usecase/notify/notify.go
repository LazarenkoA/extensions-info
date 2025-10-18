package notify

import (
	"context"
	"log"
)

type IWS interface {
	Write(msg string) error
}

type ISub interface {
	LogSubscribe(ctx context.Context, f func(msg string))
}

func Listener(ctx context.Context, sub ISub, ws IWS) {
	sub.LogSubscribe(ctx, func(msg string) {
		if err := ws.Write(msg); err != nil {
			log.Printf("ws write error: %v", err)
		}
	})
}
