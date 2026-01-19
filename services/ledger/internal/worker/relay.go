package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/justinndidit/nexus/ledger/internal/ledger"
	"github.com/justinndidit/nexus/ledger/internal/platform/broker"
	"github.com/rs/zerolog"
)

type RelayWorker struct {
	logger    *zerolog.Logger
	repo      ledger.Repository
	publisher broker.Publisher
}

func NewRelayWorker(logger *zerolog.Logger, repo ledger.PostgresRepository, publisher *broker.KafkaProducer) *RelayWorker {
	return &RelayWorker{
		logger:    logger,
		repo:      &repo,
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
	events, err := w.repo.GetOutBoxEventsForUpdate(ctx)
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
			_ = w.repo.IncrementRetryCount(ctx, event.ID.String(), err.Error())
			continue
		}

		if err := w.repo.MarkEventProcessed(ctx, event.ID.String()); err != nil {
			w.logger.Error().Err(err).Msg("failed to mark event as processed")
		}
	}
}
