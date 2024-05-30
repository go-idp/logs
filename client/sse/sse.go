package sse

// Event represents an SSE event
type Event struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Comment string      `json:"comment"`
}
