package weatherapi

import (
	"context"

	"github.com/serjnester/weather-subscription-service/domain/models"
)

type MockWeatherClient struct {
	GetForecastFn func(ctx context.Context, city string) (models.Weather, error)
}

func (m MockWeatherClient) GetForecast(ctx context.Context, city string) (models.Weather, error) {
	return m.GetForecastFn(ctx, city)
}
