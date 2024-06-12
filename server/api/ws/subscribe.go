package ws

import (
	"encoding/json"

	"github.com/go-idp/logs/server/service"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/event/cs"
)

type DataSubscribe struct {
	ID string `json:"id"`
}

func Subscribe() func(conn conn.Conn, payload cs.EventPayload, callback func(error, cs.EventPayload)) {
	return func(conn conn.Conn, payload cs.EventPayload, callback func(error, cs.EventPayload)) {
		var data DataSubscribe
		if err := payload.Bind(&data); err != nil {
			callback(err, nil)
			return
		}

		err := service.Get().Subscribe(conn.Context(), data.ID, func(err error, message string) {
			if err != nil {
				callback(err, nil)
				return
			}

			var ep cs.EventPayload
			if err := json.Unmarshal([]byte(message), &ep); err != nil {
				callback(err, nil)
				return
			}

			callback(err, ep)
		})
		if err != nil {
			callback(err, nil)
		} else {
			callback(nil, nil)
		}
	}
}
