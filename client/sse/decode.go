package sse

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

// Decode decodes an SSE event from a string
func Decode(data string) (Event, error) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Split(bufio.ScanLines)

	var event Event
	for scanner.Scan() {
		line := scanner.Text()

		if line == "event: message" {
			continue
		}

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			err := json.Unmarshal([]byte(data), &event)
			if err != nil {
				return Event{}, err
			}
		}

		if line == "" {
			break
		}
	}

	return event, nil
}

// Read reads SSE events from a reader
func Read(r io.Reader) ([]Event, error) {
	var events []Event
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "event: message" {
			continue
		}

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimPrefix(line, "data:")
			event, err := Decode(data)
			if err != nil {
				return nil, err
			}
			events = append(events, event)
		}

		if line == "" {
			continue
		}
	}

	return events, nil
}
