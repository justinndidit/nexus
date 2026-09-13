package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/justinndidit/nexus/ledger/internal/ledger/repository"
	"github.com/justinndidit/nexus/ledger/internal/platform/broker"
	"github.com/rs/zerolog"
)

type RelayWorker struct {
	logger    *zerolog.Logger
	stores    *repository.PostgresStores
	publisher broker.Publisher
}

func NewRelayWorker(logger *zerolog.Logger, stores *repository.PostgresStores, publisher *broker.KafkaProducer) *RelayWorker {
	return &RelayWorker{
		logger:    logger,
		stores:    stores,
		publisher: publisher,
	}
}

func (w *RelayWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processBatch(ctx)
		}
	}
}

func (w *RelayWorker) processBatch(ctx context.Context) {
	events, err := w.stores.OutboxStore.GetOutBoxEventsForUpdate(ctx)
	if err != nil {
		w.logger.Error().Err(err).Msg("failed to fetch outbox events")
		return
	}

	for _, event := range events {
		payload := broker.PublisherPayload{
			EventID: event.ID,
			Payload: event.Payload,
		}
		byte, err := json.Marshal(payload)
		if err != nil {
			w.logger.Error().Err(err).Msgf("failed to stringify payload for %s", event.ID)
			continue
		}

		err = w.publisher.Publish(ctx, event.Topic, event.ID.String(), byte)

		if err != nil {
			w.logger.Error().Err(err).Str("event_id", event.ID.String()).Msg("publish failed")
			_ = w.stores.OutboxStore.IncrementRetryCount(ctx, event.ID.String(), err.Error())
			continue
		}

		if err := w.stores.OutboxStore.MarkEventProcessed(ctx, event.ID.String()); err != nil {
			w.logger.Error().Err(err).Msg("failed to mark event as processed")
		}
	}
}
