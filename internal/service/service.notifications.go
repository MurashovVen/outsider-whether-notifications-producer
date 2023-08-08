package service

import (
	"context"
	"errors"
	"io"
	"time"

	grpcwhether "github.com/MurashovVen/outsider-proto/whether/golang"
	"github.com/MurashovVen/outsider-sdk/rabbitmq"
	"github.com/MurashovVen/outsider-sdk/rabbitmq/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func (s *Service) WhetherNotificationsProduce(ctx context.Context) error {
	stream, err := s.whetherService.Configurations(ctx, &grpcwhether.ConfigurationFilter{})
	if err != nil {
		return err
	}

loop:
	for {
		wc, err := stream.Recv()
		switch {
		case errors.Is(err, io.EOF):
			break loop
		case err != nil:
			return err
		}

		wcRabbit := &models.WhetherConfig{
			ChatID:      wc.ChatId,
			Temperature: wc.Temperature,
		}

		msgBody, err := wcRabbit.Marshal()
		if err != nil {
			s.log.Error("can't marshal message", zap.Int64("chat_id", wc.ChatId))
			continue
		}

		err = s.rabbit.PublishWithContext(
			ctx,
			amqp.ExchangeDirect,
			rabbitmq.WhetherNotificationsKey,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				//Expiration:   s.rabbitMessageExpirationTime.String(), todo
				Timestamp: time.Now(),
				Body:      msgBody,
			},
		)
		if err != nil {
			s.log.Error("can't publish message", zap.Int64("chat_id", wc.ChatId))
			continue
		}
	}

	return nil
}
