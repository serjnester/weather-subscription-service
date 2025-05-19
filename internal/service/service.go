package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/serjnester/weather-subscription-service/domain/models"
	"github.com/serjnester/weather-subscription-service/internal/clients/weatherapi"
	"github.com/serjnester/weather-subscription-service/internal/storage"
)

type Service interface {
	Subscribe(ctx context.Context, sub models.Subscription) error
	Confirm(ctx context.Context, token string) error
	Unsubscribe(ctx context.Context, token string) error

	WeatherForecast(ctx context.Context, city string) (models.Weather, error)
}

var _ Service = (*Impl)(nil)

type Impl struct {
	storage storage.SubscriptionStorage
	weather weatherapi.WeatherClient
}

func NewService(s storage.SubscriptionStorage, weatherCli weatherapi.WeatherClient) *Impl {
	return &Impl{storage: s, weather: weatherCli}
}

var (
	ErrAlreadySubscribed = errors.New("email already subscribed to city")
	ErrTokenNotFound     = errors.New("token not found")
)

func (s *Impl) Subscribe(ctx context.Context, sub models.Subscription) error {
	exists, err := s.storage.IsAlreadySubscribed(ctx, sub.Email, sub.City)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadySubscribed
	}

	sub.Token = generateToken()

	err = s.storage.Create(ctx, sub)
	if err != nil {
		return err
	}

	// todo send confirm email

	return nil
}

func (s *Impl) Confirm(ctx context.Context, token string) error {
	sub, err := s.storage.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTokenNotFound
		}
		return err
	}
	if sub.Confirmed {
		return nil
	}

	return s.storage.Confirm(ctx, token)
}

func (s *Impl) Unsubscribe(ctx context.Context, token string) error {
	_, err := s.storage.GetByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTokenNotFound
		}

		return err
	}
	return s.storage.Unsubscribe(ctx, token)
}

func (s *Impl) WeatherForecast(ctx context.Context, city string) (models.Weather, error) {
	forecast, err := s.weather.GetForecast(ctx, city)
	if err != nil {
		return models.Weather{}, fmt.Errorf("[Service.GetWeatherForecast] %w", err)
	}

	return forecast, nil
}

func generateToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
