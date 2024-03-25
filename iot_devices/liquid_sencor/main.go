package main

import (
	"encoding/json"
	"fmt"
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

	for range time.Tick(1 * time.Second) {
		for t := 0.0; t <= 1.0; t += 0.01 {
			publishMessage(client, 100*(1-t*t))
			time.Sleep(5 * time.Second)
		}
		for t := 0.0; t <= 1.0; t += 0.2 {
			publishMessage(client, 10+90*t)
			time.Sleep(5 * time.Second)
		}
	}
}

type LiquidLevelData struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		LiquidLevel float64 `json:"liquid_level"`
	} `json:"data"`
}

func publishMessage(client mqtt.Client, level float64) {
	liquidLevelData := LiquidLevelData{
		SensorID:  "liquid_level_sensor_01",
		Timestamp: time.Now(),
		Data: struct {
			LiquidLevel float64 `json:"liquid_level"`
		}{
			LiquidLevel: level,
		},
	}

	jsonData, err := json.Marshal(liquidLevelData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	token := client.Publish("test/topic", 0, false, string(jsonData))
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}
	fmt.Println("Successfully published message:", string(jsonData))
}
