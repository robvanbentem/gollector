package collector

import (
	"encoding/json"
	"gogarden/sensor"
	"gollector/common"
	"time"

	"github.com/robvanbentem/gocmn"
)

type DS18B20 struct {
	Source string
}

func NewDS18B20() Collector {
	h := new(DS18B20)
	h.Source = common.ConfigRoot.Source

	return h
}

func (ds *DS18B20) Handle(b []byte) error {
	jk := new(sensor.DS18B20Message)
	json.Unmarshal(b, jk)

	dt, _ := time.Parse(time.RFC3339, jk.Date)

	tx := gocmn.GetDB().MustBegin()
	tx.MustExec("INSERT INTO `data` (type, device, value, date, source) VALUES(?, ?, ?, ?, ?)", jk.Type, jk.DeviceID, jk.Temperature, dt, ds.Source)
	err := tx.Commit()

	if err != nil {
		gocmn.Log.Error("Error inserting DS18B20 data: " + err.Error())
	}

	return err
}
