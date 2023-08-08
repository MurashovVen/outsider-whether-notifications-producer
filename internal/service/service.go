package service

import (
	"time"

	grpcwhether "github.com/MurashovVen/outsider-proto/whether/golang"
	"github.com/MurashovVen/outsider-sdk/app/logger"
	"github.com/MurashovVen/outsider-sdk/rabbitmq"
)

type (
	Service struct {
		whetherService grpcwhether.WhetherClient

		rabbit                      *rabbitmq.Channel
		rabbitMessageExpirationTime time.Duration

		log *logger.Logger
	}

	Option func(*Service)
)

func New(
	whetherService grpcwhether.WhetherClient, rabbit *rabbitmq.Channel,
	rabbitMessageExpirationTime time.Duration, options ...Option,
) *Service {
	svc := &Service{
		whetherService:              whetherService,
		rabbit:                      rabbit,
		rabbitMessageExpirationTime: rabbitMessageExpirationTime,
		log:                         logger.NewNop(),
	}

	for _, opt := range options {
		opt(svc)
	}

	return svc
}

func WithLogger(log *logger.Logger) Option {
	return func(service *Service) {
		service.log = log
	}
}
