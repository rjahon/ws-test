package main

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Ascendex struct {
	url  string
	conn *websocket.Conn
}

type Message struct {
	Topic  string `json:"m"`
	Symbol string `json:"symbol"`
	Data   any    `json:"data"`
}

type BestOrderBook struct {
	Ask Order `json:"ask"`
	Bid Order `json:"bid"`
}

type Order struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

func NewAscendex(grp string) *Ascendex {
	return &Ascendex{url: "wss://ascendex.com/" + grp + "/api/pro/v1/stream"}
}

func (a *Ascendex) Connection() error {
	dialer := websocket.Dialer{
		Subprotocols: []string{"json"},
	}

	conn, _, err := dialer.Dial(a.url, nil)
	if err != nil {
		return err
	}

	a.conn = conn

	return nil
}

func (a *Ascendex) Disconnect() {
	if a.conn == nil {
		return
	}

	if err := a.conn.Close(); err != nil {
		log.Println("Error while disconnecting:", err)
		return
	}

	a.conn = nil
}

func (a *Ascendex) SubscribeToChannel(symbol string) error {
	pattern, err := regexp.Compile("^[A-Z]+_[A-Z]+$")
	if err != nil {
		return err
	}

	if !pattern.MatchString(symbol) {
		return errors.New("strings did not match")
	}

	symbol = strings.ReplaceAll(symbol, "_", "/")

	msg := map[string]string{
		"op": "sub",
		"ch": "bbo:" + symbol,
	}

	if err := a.conn.WriteJSON(msg); err != nil {
		a.Disconnect()
		return err
	}

	return nil
}

func (a *Ascendex) ReadMessagesFromChannel(ch chan<- BestOrderBook) {
	for {
		var msg Message
		if err := a.conn.ReadJSON(&msg); err != nil {
			log.Println("error while reading JSON:", err)
		}
		if msg.Topic == "bbo" {
			var (
				bbo BestOrderBook
				err error
			)

			bbo.Ask, err = Parse("ask", msg.Data)
			if err != nil {
				log.Println("Error while parsing data:", err)
				return
			}

			bbo.Bid, err = Parse("bid", msg.Data)
			if err != nil {
				log.Println("Error while parsing data:", err)
				return
			}

			ch <- bbo
		}
	}
}

func (a *Ascendex) WriteMessagesToChannel() {
	for {
		msg := map[string]string{"op": "ping"}
		if err := a.conn.WriteJSON(msg); err != nil {
			log.Println("error while writing JSON:", err)
		}
		time.Sleep(15 * time.Second)
	}
}
