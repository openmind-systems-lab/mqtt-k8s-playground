package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/openmind-systems-lab/mqtt-k8s-playground/internal/mqtt"
	paho "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	brokerURL := os.Getenv("MQTT_BROKER_URL")
	if brokerURL == "" {
		brokerURL = "tcp://omsl-emqx-listeners.default.svc.cluster.local:1883"
	}
	brokerTopic := "omsl/playground/test"

	subClient := mqtt.NewClient("omsl-subscriber", brokerURL)
	if token := subClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Subscriber connection error: %v", token.Error())
	}

	subClient.Subscribe(brokerTopic, 1, func(client paho.Client, msg paho.Message) {
		fmt.Printf("[Subscriber] Received message: %s on topic: %s\n", string(msg.Payload()), msg.Topic())
	})

	pubClient := mqtt.NewClient("omsl-publisher", brokerURL)
	if token := pubClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Publisher connection error: %v", token.Error())
	}

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		count := 1
		for range ticker.C {
			text := fmt.Sprintf("OMSL Test Message #%d", count)
			token := pubClient.Publish(brokerTopic, 1, false, text)
			token.Wait()
			fmt.Printf("[Publisher] Sent: %s\n", text)
			count++
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down clients...")
	subClient.Disconnect(250)
	pubClient.Disconnect(250)
}
