package main

import (
	"os"
	"gollector/common"
	"gollector/net"
	"gollector/collector"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"encoding/json"
	"gollector/db"
)

type Payload struct {
	Type string
	Data json.RawMessage
}

var collectors map[string]collector.CollectorCreator

func main() {
	common.LoadConfig()
	common.InitLogger()
	db.InitDB(common.ConfigRoot.Database)
	defer db.CloseDB()

	if err := net.Connect(); err != nil {
		common.Log.Fatal("Could not connect to MQTT broker")
		os.Exit(1)
	}
	defer net.Disconnect()
	common.Log.Info("Connected to MQTT broker")

	comms := make(chan MQTT.Message)
	net.Subscribe(&comms)
	common.Log.Infof("Subscribed to %s path\n", common.ConfigRoot.MQTT.Path)

	collectors = map[string]collector.CollectorCreator{}
	collectors["DS18B20"] = collector.NewDS18B20

	for {
		select {
		case m := <-comms:
			handle(m)
		}
	}
}

func handle(m MQTT.Message) {
	i := new(Payload)
	err := json.Unmarshal(m.Payload(), i)
	if err != nil {
		common.Log.Errorf("Could not decode payload\n")
		return
	}

	common.Log.Infof("Handling type %s\n", i.Type)

	if val, ok := collectors[i.Type]; ok {
		val().Handle(m.Payload())
	} else {
		common.Log.Error("Unknown Message type")
	}
}
