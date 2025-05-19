package weatherapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_GetForecast(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/forecast.json", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"location": {"name": "Kyiv"},
				"current": {
					"temp_c": 22.5,
					"humidity": 60,
					"condition": {"text": "Cloudy"},
					"wind_kph": 15
				}
			}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "dummy-key")
		weather, err := c.GetForecast(context.Background(), "Kyiv")
		require.NoError(t, err)
		require.Equal(t, 22.5, weather.Temperature)
		require.Equal(t, "Cloudy", weather.Description)
		require.Equal(t, 60, weather.Humidity)
	})

	t.Run("city not found error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{
				"error": {
					"code": 1006,
					"message": "No matching location found."
				}
			}`))
		}))
		defer server.Close()

		c := NewClient(server.URL, "dummy-key")
		_, err := c.GetForecast(context.Background(), "UnknownCity")
		require.ErrorIs(t, err, ErrCityNotFound)
	})
}
