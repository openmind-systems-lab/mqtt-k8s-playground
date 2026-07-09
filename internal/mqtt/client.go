package mqtt

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewClient(clientID, broker string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)
	opts.SetConnectTimeout(5 * time.Second)

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Printf("[%s] Connected to broker\n", clientID)
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Printf("[%s] Connection lost: %v\n", clientID, err)
	}

	return mqtt.NewClient(opts)
}
