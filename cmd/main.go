package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/serjnester/weather-subscription-service/docs"
	"github.com/serjnester/weather-subscription-service/internal/clients/weatherapi"

	"github.com/serjnester/weather-subscription-service/internal/configs"
	"github.com/serjnester/weather-subscription-service/internal/handlers"
	"github.com/serjnester/weather-subscription-service/internal/service"
	"github.com/serjnester/weather-subscription-service/internal/storage"
	"github.com/serjnester/weather-subscription-service/pkg/logging"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "serve",
		Usage: "Weather Service",
		Action: func(c *cli.Context) error {
			ctx := context.Background()

			cfg, err := configs.Load()
			if err != nil {
				log.Fatalf("config load error: %v", err)
			}

			logger, err := logging.NewLogger()
			if err != nil {
				log.Fatalf("logger init error: %v", err)
			}
			defer logger.Sync()
			dbConn, err := storage.NewDBConn(ctx, cfg.DB, logger)
			if err != nil {
				logger.Fatal("db init error", zap.Error(err))
			}
			defer dbConn.Close()

			weatherClient := weatherapi.NewClient(cfg.WeatherAPI.BaseURL, cfg.WeatherAPI.Key)
			subscriptionStorage := storage.NewSubscriptionStorage(storage.New(dbConn))
			weatherService := service.NewService(subscriptionStorage, weatherClient)

			router := handlers.NewRouter(handlers.RouterParams{Config: *cfg})
			handlers.RegisterHandlers(router, handlers.RegisterHandlersParams{
				MainHandler: handlers.NewHandler(weatherService),
			})

			addr := fmt.Sprintf(":%s", cfg.Port)
			logger.Info("starting HTTP server", zap.String("addr", addr))

			if err := http.ListenAndServe(addr, router); err != nil {
				logger.Fatal("server error", zap.Error(err))
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
