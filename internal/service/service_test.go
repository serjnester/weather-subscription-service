package service

import (
	"context"
	"errors"
	"github.com/serjnester/weather-subscription-service/internal/clients/weatherapi"
	"github.com/serjnester/weather-subscription-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/serjnester/weather-subscription-service/domain/models"
)

func TestService_WeatherForecast(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		mockWeather func() *weatherapi.MockWeatherClient
		city        string
		wantErr     error
		wantTemp    float64
	}{
		{
			name: "success",
			mockWeather: func() *weatherapi.MockWeatherClient {
				return &weatherapi.MockWeatherClient{
					GetForecastFn: func(ctx context.Context, city string) (models.Weather, error) {
						return models.Weather{
							Temperature: 25.0,
							Description: "Sunny",
							Humidity:    40,
						}, nil
					},
				}
			},
			city:     "Kyiv",
			wantErr:  nil,
			wantTemp: 25.0,
		},
		{
			name: "city not found",
			mockWeather: func() *weatherapi.MockWeatherClient {
				return &weatherapi.MockWeatherClient{
					GetForecastFn: func(ctx context.Context, city string) (models.Weather, error) {
						return models.Weather{}, weatherapi.ErrCityNotFound
					},
				}
			},
			city:     "Nowhere",
			wantErr:  weatherapi.ErrCityNotFound,
			wantTemp: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Impl{
				storage: nil,
				weather: tt.mockWeather(),
			}

			result, err := s.WeatherForecast(ctx, tt.city)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected err = %v, got = %v", tt.wantErr, err)
			}

			if result.Temperature != tt.wantTemp {
				t.Errorf("expected temperature = %v, got = %v", tt.wantTemp, result.Temperature)
			}
		})
	}
}

func TestService_Subscribe(t *testing.T) {
	weatherClient := weatherapi.MockWeatherClient{
		GetForecastFn: func(ctx context.Context, city string) (models.Weather, error) {
			return models.Weather{
				Temperature: 25.0,
				Description: "Sunny",
				Humidity:    40,
			}, nil
		},
	}
	t.Run("success", func(t *testing.T) {
		st := &storage.MockStorage{
			IsAlreadySubscribedFn: func(ctx context.Context, email, city string) (bool, error) {
				return false, nil
			},
			CreateFn: func(ctx context.Context, sub models.Subscription) error {
				return nil
			},
		}

		svc := NewService(st, weatherClient)

		err := svc.Subscribe(context.Background(), models.Subscription{
			Email:     "test@example.com",
			City:      "Kyiv",
			Frequency: "daily",
		})
		assert.NoError(t, err)
	})

	t.Run("already subscribed", func(t *testing.T) {
		st := &storage.MockStorage{
			IsAlreadySubscribedFn: func(ctx context.Context, email, city string) (bool, error) {
				return true, nil
			},
		}
		svc := NewService(st, weatherClient)
		err := svc.Subscribe(context.Background(), models.Subscription{
			Email: "test@example.com",
			City:  "Kyiv",
		})
		assert.ErrorContains(t, err, "already subscribed")
	})
}

func TestService_Confirm(t *testing.T) {
	weatherClient := weatherapi.MockWeatherClient{
		GetForecastFn: func(ctx context.Context, city string) (models.Weather, error) {
			return models.Weather{
				Temperature: 25.0,
				Description: "Sunny",
				Humidity:    40,
			}, nil
		},
	}

	t.Run("success", func(t *testing.T) {
		st := &storage.MockStorage{
			GetByTokenFn: func(ctx context.Context, token string) (models.Subscription, error) {
				return models.Subscription{Token: token, Confirmed: false}, nil
			},
			ConfirmFn: func(ctx context.Context, token string) error {
				return nil
			},
		}

		svc := NewService(st, weatherClient)
		err := svc.Confirm(context.Background(), "some-token")
		assert.NoError(t, err)
	})

	t.Run("already confirmed", func(t *testing.T) {
		st := &storage.MockStorage{
			GetByTokenFn: func(ctx context.Context, token string) (models.Subscription, error) {
				return models.Subscription{Token: token, Confirmed: true}, nil
			},
		}

		svc := NewService(st, weatherClient)
		err := svc.Confirm(context.Background(), "some-token")
		assert.NoError(t, err)
	})
}

func TestService_Unsubscribe(t *testing.T) {
	weatherClient := weatherapi.MockWeatherClient{
		GetForecastFn: func(ctx context.Context, city string) (models.Weather, error) {
			return models.Weather{
				Temperature: 25.0,
				Description: "Sunny",
				Humidity:    40,
			}, nil
		},
	}

	t.Run("success", func(t *testing.T) {
		st := &storage.MockStorage{
			GetByTokenFn: func(ctx context.Context, token string) (models.Subscription, error) {
				return models.Subscription{Token: token, Confirmed: true}, nil
			},
			UnsubscribeFn: func(ctx context.Context, token string) error {
				return nil
			},
		}

		svc := NewService(st, weatherClient)
		err := svc.Unsubscribe(context.Background(), "some-token")
		assert.NoError(t, err)
	})

	t.Run("not confirmed", func(t *testing.T) {
		st := &storage.MockStorage{
			GetByTokenFn: func(ctx context.Context, token string) (models.Subscription, error) {
				return models.Subscription{Token: token, Confirmed: false}, nil
			},
			UnsubscribeFn: func(ctx context.Context, token string) error {
				return nil
			},
		}

		svc := NewService(st, weatherClient)
		err := svc.Unsubscribe(context.Background(), "some-token")
		assert.NoError(t, err)
	})
}
