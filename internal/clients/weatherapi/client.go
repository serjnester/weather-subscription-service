package weatherapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/serjnester/weather-subscription-service/domain/models"
)

type WeatherClient interface {
	GetForecast(ctx context.Context, city string) (models.Weather, error)
}

var _ WeatherClient = (*Client)(nil)

var ErrCityNotFound = errors.New("city not found")

const apiErrCodeCityNotFound = 1006

type Client struct {
	baseURL string
	key     string
	resty   *resty.Client
}

func NewClient(baseURL, key string) *Client {
	return &Client{
		baseURL: baseURL,
		key:     key,
		resty:   resty.New(),
	}
}

type weatherAPIResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  int     `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		WindKph float64 `json:"wind_kph"`
	} `json:"current"`
}

type weatherAPIErr struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Client) GetForecast(ctx context.Context, city string) (models.Weather, error) {
	var apiResp weatherAPIResponse
	var apiErr weatherAPIErr

	resp, err := c.resty.R().SetContext(ctx).
		SetQueryParams(map[string]string{
			"key": c.key,
			"q":   city,
		}).
		SetResult(&apiResp).
		SetError(&apiErr).
		Get(fmt.Sprintf("%s/forecast.json", c.baseURL))

	if err != nil {
		return models.Weather{}, fmt.Errorf("weather api request error: %w", err)
	}

	if resp.IsError() {
		if apiErr.Error.Code == apiErrCodeCityNotFound {
			return models.Weather{}, ErrCityNotFound
		}
		return models.Weather{}, fmt.Errorf("weather api returned error: %s", apiErr.Error.Message)
	}

	return models.Weather{
		Temperature: apiResp.Current.TempC,
		Description: apiResp.Current.Condition.Text,
		Humidity:    apiResp.Current.Humidity,
	}, nil
}
