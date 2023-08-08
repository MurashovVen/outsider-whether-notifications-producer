package main

import (
	"time"

	"github.com/MurashovVen/outsider-sdk/app/configuration"
)

type config struct {
	configuration.Default
	configuration.RabbitMQ

	// todo move to sdk
	WhetherGRPCClientAddr string        `desc:"Address of whether service" default:"whether:5000" split_words:"true"`
	NotificationsTimer    time.Duration `desc:"Period of sending notifications" default:"5m" split_words:"true"`
}
