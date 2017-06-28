package net

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"gollector/common"
)

type Message struct {
	Path    string
	Message []byte
}

var client MQTT.Client
var comms *chan MQTT.Message

var path string

func Connect() error {
	cfg := common.ConfigRoot.MQTT

	opts := MQTT.NewClientOptions().AddBroker(cfg.Broker)
	opts.SetClientID(cfg.Name)
	opts.SetDefaultPublishHandler(handleMessage)

	client = MQTT.NewClient(opts)
	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		return token.Error()
	}

	path = cfg.Path

	return nil
}

func Subscribe(c *chan MQTT.Message) error {
	if token := client.Subscribe(path, 2, handleMessage); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	comms = c

	return nil
}

func Disconnect() {
	if client.IsConnected() {
		client.Disconnect(250)
	}
}

func handleMessage(c MQTT.Client, m MQTT.Message) {
	*comms <- m
}
