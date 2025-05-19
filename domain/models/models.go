package models

import "github.com/serjnester/weather-subscription-service/domain/enums"

type Weather struct {
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
}

type Subscription struct {
	Email     string          `json:"email"`
	City      string          `json:"city"`
	Frequency enums.Frequency `json:"frequency"`
	Token     string          `json:"token"`
	Confirmed bool            `json:"confirmed"`
}
