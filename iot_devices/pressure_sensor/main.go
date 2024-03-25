package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Client connected")

	for range time.Tick(1 * time.Minute) {
		pressure := generatePressure()
		publishMessage(client, pressure)
	}
}

type PressureData struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		Pressure float64 `json:"pressure"`
	} `json:"data"`
}

func publishMessage(client mqtt.Client, pressure float64) {
	pressureData := PressureData{
		SensorID:  "pressure_sensor_01",
		Timestamp: time.Now(),
		Data: struct {
			Pressure float64 `json:"pressure"`
		}{
			Pressure: pressure,
		},
	}

	jsonData, err := json.Marshal(pressureData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	token := client.Publish("test/pressure", 0, false, string(jsonData))
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}
	fmt.Println("Successfully published message:", string(jsonData))
}

func generatePressure() float64 {
	return 1900 + rand.Float64()*(300)
}
