package echat

import (
	"encoding/json"
	"time"
)

type Message interface {
	Dest() string
	Src() string
	Time() time.Time
	Json() ([]byte, error)
	FromJson(bt []byte) error
}

type CommonMessage struct {
	From    string          `json:"src"`
	To      string          `json:"dest"`
	Content json.RawMessage `json:"content"`
}

func (c *CommonMessage) Src() string {
	return c.From
}

func (c *CommonMessage) Dest() string {
	return c.To
}

func (c *CommonMessage) Time() time.Time {
	return time.Now()
}

func (c *CommonMessage) Json() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CommonMessage) FromJson(bt []byte) error {
	return json.Unmarshal(bt, c)
}

type MessageParser func([] byte) (Message, error)

func DefaultParser(bt []byte) (Message, error) {
	msg := &CommonMessage{}
	err := msg.FromJson(bt)
	return msg, err
}
