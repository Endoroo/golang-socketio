package main

import (
	"fmt"
	"github.com/endoroo/golang-socketio"
	"github.com/endoroo/golang-socketio/transport"
	"log"
	"time"
)

const (
	userId      = 424242
	pairId      = 424
	accessToken = "access_token"
)

type HeadersArgs struct {
	Authorization string `json:"Authorization"`
}

type AuthArgs struct {
	Headers HeadersArgs `json:"headers"`
}

type SubscribeArgs struct {
	Channel string   `json:"channel"`
	Auth    AuthArgs `json:"auth"`
}

type TickerMessage struct {
	ClosedOrders    int `json:"closedOrders"`
	Id              int
	LastPrice       string `json:"lastPrice"`
	LastPriceDayAgo string `json:"lastPriceDayAgo"`
	MarketVolume    string `json:"market_volume"`
	MaxBuy          string `json:"maxBuy"`
	MinSell         string `json:"minSell"`
	Precision       int
	Socket          string
	Spread          string
	VolumeSum       string `json:"volumeSum"`
}

type OrderMessage struct {
	Id             int
	UserId         int `json:"user_id"`
	CurrencyPairId int `json:"currency_pair_id"`
	Amount         string
	Price          string
	Amount2        string
	Type           string
	Socket         string
}

func main() {
	client, err := gosocketio.Dial(
		gosocketio.GetUrl("socket.stex.com", 443, true),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = client.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Printf("connected %s", h.Id())
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	err = subscribeOnPublicChannelExample(client)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = subscribeOnPrivateChannelExample(client)
	if err != nil {
		log.Fatal(err)
		return
	}

	time.Sleep(1 * time.Minute)
	client.Close()
}

func subscribeOnPublicChannelExample(client *gosocketio.Client) error {
	err := client.On("App\\\\Events\\\\Ticker", func(h *gosocketio.Channel, message TickerMessage) {
		log.Println("tick ", message)
	})
	if err != nil {
		return err
	}

	err = client.Emit("subscribe", SubscribeArgs{
		Channel: "rate",
	})
	if err != nil {
		return err
	}

	return nil
}

func subscribeOnPrivateChannelExample(client *gosocketio.Client) error {
	err := client.On("App\\\\Events\\\\UserOrder", func(h *gosocketio.Channel, message OrderMessage) {
		log.Println("order status ", message)
	})
	if err != nil {
		return err
	}

	err = client.Emit("subscribe", SubscribeArgs{
		Channel: "private-sell_user_data_u" + fmt.Sprintf("%d", userId) + "c" + fmt.Sprintf("%d", pairId),
		Auth:    AuthArgs{Headers: HeadersArgs{Authorization: "Bearer " + accessToken}},
	})
	if err != nil {
		return err
	}
	return nil
}
