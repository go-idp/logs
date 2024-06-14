package ws

import (
	"github.com/go-idp/logs/server/service"
	"github.com/go-zoox/websocket/conn"
	cs "github.com/go-zoox/websocket/extension/event/entity"
)

type DataPublish struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func Publish() func(conn conn.Conn, payload cs.EventPayload, callback func(error, cs.EventPayload)) {
	return func(conn conn.Conn, payload cs.EventPayload, callback func(error, cs.EventPayload)) {
		var data DataPublish
		if err := payload.Bind(&data); err != nil {
			callback(err, nil)
			return
		}

		err := service.Get().Publish(conn.Context(), data.ID, data.Message)
		if err != nil {
			callback(err, nil)
		} else {
			callback(nil, nil)
		}
	}
}
