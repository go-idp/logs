package ws

import (
	"github.com/go-idp/logs/server/service"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event/cs"
)

type DataOpen struct {
	ID string `json:"id"`
}

func Open() func(conn conn.Conn, payload cs.EventPayload, callback func(error, cs.EventPayload)) {
	return func(conn conn.Conn, payload cs.EventPayload, callback func(error, cs.EventPayload)) {
		var data DataOpen
		if err := payload.Bind(&data); err != nil {
			callback(err, nil)
			return
		}

		err := service.Get().Open(conn.Context(), data.ID)
		if err != nil {
			callback(err, nil)
		} else {
			callback(nil, nil)
		}
	}
}
