package collector

import (
	"gogarden/sensor"
	"encoding/json"
	"fmt"
)

type DS18B20 struct{}

func NewDS18B20() Collector {
	return new(DS18B20)
}

func (ds *DS18B20) Handle(b []byte) string {
	jk := new(sensor.DS18B20Message)
	json.Unmarshal(b, jk)

	return fmt.Sprintf("Device: %s, Temp: %.3f, Date: %s", jk.DeviceID, jk.Temperature, jk.Date)
}
