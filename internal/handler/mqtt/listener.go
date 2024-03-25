package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MetricService interface {
	SetFuelLevel(fuelLevel float64, dateTime time.Time) error
	IncrementMQTTEvents(dateTime time.Time) error
	SetTirePressure(tire float64, dateTime time.Time) error
	SetLocation(latitude, longitude float64, dateTime time.Time) error
	SetIsMoving(isMoving bool, dateTime time.Time) error
}

func NewMQTTListener(service MetricService) *MQTTListener {
	return &MQTTListener{
		service: service,
	}
}

type MQTTListener struct {
	service MetricService
}

func (l *MQTTListener) MessageRouter(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	l.service.IncrementMQTTEvents(time.Now())

	switch msg.Topic() {
	case "test/topic":
		l.liquidLevelHandler(msg)
	case "test/coordinate":
		l.coordinatesHandler(msg)
	case "test/pressure":
		l.pressureHandler(msg)
	case "test/moving":
		l.movingHandler(msg)
	}

}

type LiquidLevelData struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		LiquidLevel float64 `json:"liquid_level"`
	} `json:"data"`
}

func (l *MQTTListener) liquidLevelHandler(msg mqtt.Message) {
	var data LiquidLevelData
	err := json.Unmarshal(msg.Payload(), &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	l.service.SetFuelLevel(data.Data.LiquidLevel, data.Timestamp)
}

type CoordinatesData struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"data"`
}

func (l *MQTTListener) coordinatesHandler(msg mqtt.Message) {
	var coordinates CoordinatesData
	err := json.Unmarshal(msg.Payload(), &coordinates)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	l.service.SetLocation(coordinates.Data.Latitude, coordinates.Data.Longitude, coordinates.Timestamp)
}

type PressureData struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		Pressure float64 `json:"pressure"`
	} `json:"data"`
}

func (l *MQTTListener) pressureHandler(msg mqtt.Message) {
	var pressureData PressureData
	err := json.Unmarshal(msg.Payload(), &pressureData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	l.service.SetTirePressure(pressureData.Data.Pressure, pressureData.Timestamp)
}

type MovingData struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		IsMoving bool `json:"is_moving"`
	} `json:"data"`
}

func (l *MQTTListener) movingHandler(msg mqtt.Message) {
	var movingData MovingData
	err := json.Unmarshal(msg.Payload(), &movingData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	l.service.SetIsMoving(movingData.Data.IsMoving, movingData.Timestamp)
}
