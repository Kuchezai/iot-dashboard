package prometheus

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetricStorage struct {
	liquidMetric      prometheus.Gauge
	pressureMetric    prometheus.Gauge
	mqttMessageMetric prometheus.Counter
	geoMetric         *prometheus.GaugeVec
	isMovingMetric    prometheus.Gauge
}

func NewPrometheusMetricStorage() *PrometheusMetricStorage {
	mqttMessageMetric := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_mqtt_messages_received_total",
		Help: "Total number of MQTT messages received by myapp.",
	})
	prometheus.MustRegister(mqttMessageMetric)

	liquidMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "liquid_level",
			Help: "Liquid level at percent",
		},
	)
	prometheus.MustRegister(liquidMetric)

	geoMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "geo",
			Help: "Geographical coordinates.",
		},
		[]string{"latitude", "longitude"},
	)
	prometheus.MustRegister(geoMetric)

	pressureMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "pressure",
			Help: "Pressure in Pascal",
		},
	)
	prometheus.MustRegister(pressureMetric)

	isMovingMetric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "is_moving",
			Help: "Is car moving",
		},
	)
	prometheus.MustRegister(isMovingMetric)

	return &PrometheusMetricStorage{
		liquidMetric:      liquidMetric,
		geoMetric:         geoMetric,
		mqttMessageMetric: mqttMessageMetric,
		pressureMetric:    pressureMetric,
		isMovingMetric:    isMovingMetric,
	}
}

func (p *PrometheusMetricStorage) SetFuelLevel(fuelLevel float64, dateTime time.Time) error {
	p.liquidMetric.Set(fuelLevel)
	return nil
}

func (p *PrometheusMetricStorage) IncrementMQTTEvents(dateTime time.Time) error {
	p.mqttMessageMetric.Inc()
	return nil
}

func (p *PrometheusMetricStorage) SetTirePressure(tire float64, dateTime time.Time) error {
	p.pressureMetric.Set(tire)
	return nil
}

func (p *PrometheusMetricStorage) SetLocation(latitude, longitude float64, dateTime time.Time) error {
	p.geoMetric.WithLabelValues(fmt.Sprintf("%f", latitude), fmt.Sprintf("%f", longitude)).Set(1)
	return nil
}

func (p *PrometheusMetricStorage) SetIsMoving(isMoving bool, dateTime time.Time) error {
	var value float64
	if isMoving {
		value = 1
	} else {
		value = 0
	}

	go func() {
		time.Sleep(30 * time.Second)
		p.isMovingMetric.Set(0)
	}()

	p.isMovingMetric.Set(value)
	return nil
}
