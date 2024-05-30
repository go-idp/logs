package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Encode encodes an SSE event to a string
func Encode(event Event) string {
	data, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("event: %s\ndata: %s\n\n", event.Type, string(data))
}

// Send sends an SSE event to a client
func Send(w http.ResponseWriter, event Event) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Write([]byte(Encode(event)))
}
