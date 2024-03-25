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

	for range time.Tick(30 * time.Second) {
		latitude, longitude := generateSPBCoordinates()
		publishMessage(client, latitude, longitude)
	}
}

func generateSPBCoordinates() (float64, float64) {
	minLatitude := 59.8
	maxLatitude := 60.1
	minLongitude := 30.1
	maxLongitude := 30.6

	latitude := rand.Float64()*(maxLatitude-minLatitude) + minLatitude
	longitude := rand.Float64()*(maxLongitude-minLongitude) + minLongitude

	return latitude, longitude
}

type CoordinatesData struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	Data      struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"data"`
}

func publishMessage(client mqtt.Client, latitude, longitude float64) {
	coordinatesData := CoordinatesData{
		SensorID:  "coordinates_sensor_01",
		Timestamp: time.Now(),
		Data: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  latitude,
			Longitude: longitude,
		},
	}

	jsonData, err := json.Marshal(coordinatesData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	token := client.Publish("test/coordinate", 0, false, string(jsonData))
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}
	fmt.Println("Successfully published message:", string(jsonData))
}
