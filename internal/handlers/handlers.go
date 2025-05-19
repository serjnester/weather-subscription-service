package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/serjnester/weather-subscription-service/domain/enums"
	"github.com/serjnester/weather-subscription-service/domain/models"
	"github.com/serjnester/weather-subscription-service/internal/clients/weatherapi"
	"github.com/serjnester/weather-subscription-service/internal/service"
	"net/http"
)

type Handler interface {
	GetWeather(c *gin.Context)
	Subscribe(c *gin.Context)
	Unsubscribe(c *gin.Context)
	ConfirmSubscription(c *gin.Context)
}

type CreditLimitsParams struct {
}

func NewHandler(service service.Service) Handler {
	return &handler{
		Service: service,
	}
}

type handler struct {
	Service service.Service
}

// GetWeather godoc
//
//	@Summary		Get current weather for a city
//	@Description	Returns the current weather forecast for the specified city using WeatherAPI.com.
//	@Tags			weather
//	@Accept			json
//	@Produce		json
//	@Param			city	query		string	true	"City name for weather forecast"
//	@Success		200		{object}	models.Weather
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		404		{string}	string	"City not found"
//	@Router			/api/weather [get]
func (h *handler) GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, "Missing 'city' query parameter")
		return
	}

	weather, err := h.Service.WeatherForecast(c.Request.Context(), city)
	if err != nil {
		if errors.Is(err, weatherapi.ErrCityNotFound) {
			c.JSON(http.StatusNotFound, "City not found")
		} else {
			c.JSON(http.StatusInternalServerError, "Failed to fetch weather")
		}
		return
	}

	c.JSON(http.StatusOK, weather)
}

type subscribeForm struct {
	Email     string          `form:"email" binding:"required,email"`
	City      string          `form:"city" binding:"required"`
	Frequency enums.Frequency `form:"frequency" binding:"required,oneof=hourly daily"`
}

// Subscribe godoc
//
//	@Summary		Subscribe to weather updates
//	@Description	Subscribe an email to receive weather updates for a specific city with chosen frequency.
//	@Tags			subscription
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			email		formData	string	true	"Email address to subscribe"
//	@Param			city		formData	string	true	"City for weather updates"
//	@Param			frequency	formData	string	true	"Frequency of updates (hourly or daily)" Enums(hourly, daily)
//	@Success		200			{string}	string	"Subscription successful. Confirmation email sent."
//	@Failure		400			{string}	string	"Invalid input"
//	@Failure		409			{string}	string	"Email already subscribed"
//	@Router			/api/subscribe [post]
func (h *handler) Subscribe(c *gin.Context) {
	var form subscribeForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid input")
		return
	}

	err := h.Service.Subscribe(c.Request.Context(), models.Subscription{
		Email:     form.Email,
		City:      form.City,
		Frequency: form.Frequency,
	})
	if err != nil {
		if errors.Is(err, service.ErrAlreadySubscribed) {
			c.JSON(http.StatusConflict, "Email already subscribed")
			return
		}

		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, "Subscription successful. Confirmation email sent.")
}

// ConfirmSubscription godoc
//
//	@Summary		Confirm email subscription
//	@Description	Confirms a subscription using the token sent in the confirmation email.
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Confirmation token"
//	@Success		200		{string}	string	"Subscription confirmed successfully"
//	@Failure		400		{string}	string	"Invalid token"
//	@Failure		404		{string}	string	"Token not found"
//	@Router			/api/confirm/{token} [get]
func (h *handler) ConfirmSubscription(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, "Invalid token")
		return
	}

	err := h.Service.Confirm(c.Request.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTokenNotFound):
			c.JSON(http.StatusNotFound, "Token not found")
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, "Subscription confirmed successfully")
}

// Unsubscribe godoc
//
//	@Summary		Unsubscribe from weather updates
//	@Description	Unsubscribes an email from weather updates using the token sent in emails.
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Unsubscribe token"
//	@Success		200		{string}	string	"Unsubscribed successfully"
//	@Failure		400		{string}	string	"Invalid token"
//	@Failure		404		{string}	string	"Token not found"
//	@Router			/api/unsubscribe/{token} [get]
func (h *handler) Unsubscribe(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, "Invalid token")
		return
	}

	err := h.Service.Unsubscribe(c.Request.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTokenNotFound):
			c.JSON(http.StatusNotFound, "Token not found")
		default:

			c.JSON(http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, "Unsubscribed successfully")
}
