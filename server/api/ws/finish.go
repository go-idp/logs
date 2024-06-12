package ws

import (
	"github.com/go-idp/logs/server/service"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event/cs"
)

type DataFinish struct {
	ID string `json:"id"`
}

func Finish() func(conn conn.Conn, payload cs.EventPayload, callback func(error, cs.EventPayload)) {
	return func(conn conn.Conn, payload cs.EventPayload, callback func(error, cs.EventPayload)) {
		var data DataFinish
		if err := payload.Bind(&data); err != nil {
			callback(err, nil)
			return
		}

		err := service.Get().Finish(conn.Context(), data.ID)
		if err != nil {
			callback(err, nil)
		} else {
			callback(nil, nil)
		}
	}
}
