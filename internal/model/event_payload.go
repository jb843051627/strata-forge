package model

import (
	"encoding/json"
	"fmt"
)

type StateEvent struct {
	From string `json:"from"`
	To   string `json:"to"`
	By   string `json:"by"`
	Note string `json:"note"`
}

func EncodeStateEvent(event StateEvent) (string, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("encode state event: %w", err)
	}
	return string(data), nil
}

func DecodeStateEvent(payload string) (StateEvent, error) {
	var event StateEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		return StateEvent{}, fmt.Errorf("decode state event: %w", err)
	}
	return event, nil
}
