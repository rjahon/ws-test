package main

import (
	"testing"
	"time"
)

var a = func() *Ascendex {
	a := NewAscendex("0")
	a.Connection()
	return a
}()

func TestConnection(t *testing.T) {
	a := NewAscendex("0")
	err := a.Connection()
	defer a.Disconnect()
	if err != nil {
		t.Error("Connection() failed", err.Error())
	}
}

func TestWriteMessagesToChannel(t *testing.T) {
	ch := make(chan bool)
	go func() {
		for {
			var msg Message
			a.conn.ReadJSON(&msg)
			if msg.Topic == "pong" {
				ch <- true
				break
			}
		}
	}()

	go a.WriteMessagesToChannel()
	select {
	case <-ch:
		return
	case <-time.After(5 * time.Second):
		t.Errorf("WriteMessageToChannel() failed")
	}
}

func TestReadMessagesFromChannel(t *testing.T) {
	a.SubscribeToChannel("BTC_USDT")

	ch := make(chan BestOrderBook)
	go a.ReadMessagesFromChannel(ch)
	select {
	case <-ch:
		return
	case <-time.After(4 * time.Second):
		t.Errorf("ReadMessagesFromChannel() failed")
	}
}

func TestSubscribeToChannel(t *testing.T) {
	if err := a.SubscribeToChannel("BTC_USDT"); err != nil {
		t.Errorf("SubscribeToChannel() failed")
	}
}
