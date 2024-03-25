package internal

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	listener "github.com/Kuchezai/iot-car/internal/handler/mqtt"
	repo "github.com/Kuchezai/iot-car/internal/repo/prometheus"
	service "github.com/Kuchezai/iot-car/internal/service"
	telegram "github.com/Kuchezai/iot-car/internal/telegram"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func App() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	metricRepo := repo.NewPrometheusMetricStorage()
	telegramChatID, err := strconv.Atoi(os.Getenv("CHAT_ID"))
	if err != nil {
		fmt.Println("Error loading chat_id:", err)
		os.Exit(1)
	}
	telegramSender := telegram.NewTelegramSender(telegramChatID, os.Getenv("API_KEY"))
	metricService := service.NewMetricService(metricRepo, telegramSender)
	mqttListener := listener.NewMQTTListener(metricService)
	setupMQTTClient(mqttListener)

	// Setup HTTP server to expose metrics.
	http.Handle("/metrics", promhttp.Handler())
	httpPort := os.Getenv("HTTP_PORT")
	fmt.Println("Server listening on :" + httpPort)
	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type MQTTListener interface {
	MessageRouter(client mqtt.Client, msg mqtt.Message)
}

func setupMQTTClient(listener MQTTListener) mqtt.Client {
	mqttBroker := os.Getenv("MQTT_BROKER")
	mqttPort := os.Getenv("MQTT_PORT")

	opts := mqtt.NewClientOptions().AddBroker(mqttBroker + ":" + mqttPort)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(listener.MessageRouter)
	opts.SetPingTimeout(1 * time.Second)

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := mqttClient.SubscribeMultiple(map[string]byte{
		"test/topic":      0,
		"test/coordinate": 0,
		"test/pressure":   0,
		"test/moving":     0,
	}, listener.MessageRouter); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	return mqttClient
}
