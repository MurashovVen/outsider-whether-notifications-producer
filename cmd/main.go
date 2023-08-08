package main

import (
	"context"
	"errors"

	grpcwhether "github.com/MurashovVen/outsider-proto/whether/golang"
	"github.com/MurashovVen/outsider-sdk/app"
	"github.com/MurashovVen/outsider-sdk/app/configuration"
	"github.com/MurashovVen/outsider-sdk/app/logger"
	"github.com/MurashovVen/outsider-sdk/app/termination"
	"github.com/MurashovVen/outsider-sdk/grpc"
	"github.com/MurashovVen/outsider-sdk/rabbitmq"
	"go.uber.org/zap"

	"outsider-wthether-notifications-producer/internal/producer"
	"outsider-wthether-notifications-producer/internal/service"
)

func main() {
	var (
		cfg = new(config)

		ctx = context.Background()
	)

	configuration.MustProcessConfig(cfg)

	var (
		log = logger.MustCreateLogger(cfg.Env)

		whetherClientConn = grpc.MustConnect(cfg.WhetherGRPCClientAddr, grpc.DefaultDialOptions(log)...)

		rabbitClient = rabbitmq.MustConnect(cfg.RabbitMQURL, rabbitmq.ChannelWithLogger(log))
	)

	application := app.New(
		log,
		app.AppendWorks(
			producer.New(
				service.New(
					grpcwhether.NewWhetherClient(whetherClientConn),
					rabbitClient,
					cfg.NotificationsTimer,
					service.WithLogger(log),
				),
				cfg.NotificationsTimer,
				log,
			),
		),
	)

	if err := application.Run(ctx); err != nil && !errors.Is(err, termination.ErrStopped) {
		log.Error("running error", zap.Error(err))
	}
}
