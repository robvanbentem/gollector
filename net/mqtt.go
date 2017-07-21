package net

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"gollector/common"
	"time"
	"gocmn"
	"go/token"
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
	opts.SetDefaultPublishHandler(defaultMessageHandler)

	client = MQTT.NewClient(opts)
	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		return token.Error()
	}

	path = cfg.Path

	go checkConnection()

	return nil
}

func Subscribe(c *chan MQTT.Message) error {
	if tkn := client.Subscribe(path, 2, func(c MQTT.Client, m MQTT.Message) { *comms <- m }); tkn.Wait() && tkn.Error() != nil {
		return tkn.Error()
	}

	comms = c
	return nil
}

func Disconnect() {
	if client.IsConnected() {
		client.Disconnect(250)
	}
}

func defaultMessageHandler(c MQTT.Client, m MQTT.Message) {
	gocmn.Log.Warning("Got message on default message handler")
}

func checkConnection() {
	ticker := time.NewTicker(time.Second * 60).C
	for {
		select {
		case <-ticker:
			if !client.IsConnected() {
				gocmn.Log.Info("MQTT disconnected, reconnecting..")
				Connect()
			} else {
				gocmn.Log.Debug("MQTT connection OK")
			}

			Subscribe(comms)
		}
	}
}
