package main

import (
	"encoding/json"
	"net/url"

	"github.com/payfazz/go-errors"
)

type messageRepresentation struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func (m *messageRepresentation) parseFromJSON(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (m *messageRepresentation) asJSON() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m *messageRepresentation) parseFromURLEncoded(data []byte) error {
	values, err := url.ParseQuery(string(data))
	if err != nil {
		return errors.Wrap(err)
	}
	m.Title = values.Get("title")
	m.Message = values.Get("message")
	return nil
}

func (m *messageRepresentation) asURLEncoded() []byte {
	values := make(url.Values)
	values.Set("title", m.Title)
	values.Set("message", m.Message)
	return []byte(values.Encode())
}
