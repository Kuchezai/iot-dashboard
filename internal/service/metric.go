package service

import (
	"fmt"
	"time"
)

type MetricStorage interface {
	SetFuelLevel(fuelLevel float64, dateTime time.Time) error
	IncrementMQTTEvents(dateTime time.Time) error
	SetTirePressure(pressure float64, dateTime time.Time) error
	SetLocation(latitude, longitude float64, dateTime time.Time) error
	SetIsMoving(isMoving bool, dateTime time.Time) error
}

type TelegramSender interface {
	SendMessage(messageText string) error
}

type MetricService struct {
	telegramSender TelegramSender
	storage        MetricStorage
}

func NewMetricService(storage MetricStorage, telegramSender TelegramSender) *MetricService {
	return &MetricService{
		storage:        storage,
		telegramSender: telegramSender,
	}
}

func (s *MetricService) SetFuelLevel(fuelLevel float64, dateTime time.Time) error {
	return s.storage.SetFuelLevel(fuelLevel, dateTime)
}

func (s *MetricService) IncrementMQTTEvents(dateTime time.Time) error {
	return s.storage.IncrementMQTTEvents(dateTime)
}

func (s *MetricService) SetTirePressure(tire float64, dateTime time.Time) error {
	return s.storage.SetTirePressure(tire, dateTime)
}

func (s *MetricService) SetLocation(latitude, longitude float64, dateTime time.Time) error {
	return s.storage.SetLocation(latitude, longitude, dateTime)
}

func (s *MetricService) SetIsMoving(isMoving bool, dateTime time.Time) error {
	if isMoving && !s.isValidDateTime(dateTime) {
		message := fmt.Sprintf("❗❗❗Alert! The provided time (%s) is not within the allowed range.", dateTime.Format(time.RFC3339))
		s.telegramSender.SendMessage(message)
	}

	return s.storage.SetIsMoving(isMoving, dateTime)
}

func (s *MetricService) isValidDateTime(dateTime time.Time) bool {
	startTime := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 8, 0, 0, 0, dateTime.Location())
	endTime := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 9, 0, 0, 0, dateTime.Location())
	return dateTime.After(startTime) && dateTime.Before(endTime)
}
