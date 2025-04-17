package main

import (
	"net"
	"time"
)

type Message struct {
	Type    string
	Payload interface{}
}

type MessageQueue struct {
	events []Message
	conn   net.Conn
	buffer []byte
}

func NewMessageQueue(conn net.Conn) *MessageQueue {
	return &MessageQueue{
		events: make([]Message, 0),
		conn:   conn,
		buffer: make([]byte, 1024),
	}
}

func (q *MessageQueue) ReadFromConnection() bool {
	q.conn.SetReadDeadline(time.Now().Add(time.Millisecond))
	n, err := q.conn.Read(q.buffer)

	if err == nil {
		message := string(q.buffer[:n])
		q.events = append(q.events, Message{
			Type:    "MESSAGE",
			Payload: message,
		})
		return true
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return false
	}

	q.events = append(q.events, Message{
		Type:    "ERROR",
		Payload: err.Error(),
	})

	return true
}

func (q *MessageQueue) Pop() *Message {
	if len(q.events) == 0 {
		return nil
	}

	evt := q.events[0]
	q.events = q.events[1:]
	return &evt
}
