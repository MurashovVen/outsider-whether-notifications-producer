package producer

import (
	"context"

	"github.com/MurashovVen/outsider-sdk/app"
)

var (
	_ app.Work = (*EventProducer)(nil)
)

func (ep *EventProducer) Runner(ctx context.Context) func() error {
	return func() error {
		return ep.start(ctx)
	}
}

func (ep *EventProducer) Name() string {
	return "EventProducer"
}
