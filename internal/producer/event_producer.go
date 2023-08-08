package producer

import (
	"context"
	"time"

	"github.com/MurashovVen/outsider-sdk/app/logger"
	"go.uber.org/zap"

	"outsider-wthether-notifications-producer/internal/service"
)

type EventProducer struct {
	svc *service.Service

	notificationTime time.Duration

	log *logger.Logger
}

func New(svc *service.Service, notificationTime time.Duration, log *logger.Logger) *EventProducer {
	return &EventProducer{
		svc:              svc,
		notificationTime: notificationTime,
		log:              log.Named("EventProducer"),
	}
}

func (ep *EventProducer) start(ctx context.Context) error {
	timer := time.NewTicker(ep.notificationTime)

	ep.do(ctx)

	for {
		select {
		case <-timer.C:
			ep.do(ctx)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ep *EventProducer) do(ctx context.Context) {
	ep.log.Info("starting producing events")

	if err := ep.svc.WhetherNotificationsProduce(ctx); err != nil {
		ep.log.Error("finished producing events", zap.Error(err))
		return
	}

	ep.log.Info("finished producing events")
}
