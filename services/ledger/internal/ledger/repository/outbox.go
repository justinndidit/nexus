package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/justinndidit/nexus/ledger/internal/ledger/domain"
	"github.com/rs/zerolog"
)

type OutboxRepo struct {
	repo   Repository
	logger *zerolog.Logger
}

func NewOutBoxRepo(r Repository, logger *zerolog.Logger) *OutboxRepo {
	return &OutboxRepo{
		repo:   r,
		logger: logger,
	}
}

func (o *OutboxRepo) CreateOutBoxEvent(ctx context.Context, payload domain.CreateOutboxEventRequest) error {
	stmt := `
		INSERT INTO outbox_events(event_type, payload, status, idempotency_key, priority, producer, queue_topic)
		VALUES (@eventType, @payload, @status, @idempotency_key, @priority, @producer, @queueTopic)
	`
	cmd, err := o.repo.Exec(ctx, stmt, pgx.NamedArgs{
		"payload":         payload.Payload,
		"idempotency_key": payload.IdempotencyKey,
		"eventType":       payload.EventType,
		"producer":        payload.Producer,
		"priority":        int(payload.Priority),
		"status":          payload.Status,
		// "queueTopic":      payload.Topic,
	})
	if err != nil {
		o.logger.Error().Err(err).Str("func", "CreateOutBoxEvent").Msg("failed to insert event")
		return err
	}

	if cmd.RowsAffected() == 0 {
		o.logger.Error().Str("func", "CreateOutBoxEvent").Msg("failed to ccreate outbox event")
		return fmt.Errorf("No rows affected")
	}

	return nil
}

func (o *OutboxRepo) GetOutBoxEventsForUpdate(ctx context.Context) ([]domain.OutBoxEvent, error) {
	stmt := `
			UPDATE outbox_events
			SET status = 'PROCESSING',
					locked_at = NOW()
			WHERE id IN (
					SELECT id FROM outbox_events
					WHERE status = 'PENDING'
					ORDER BY created_at ASC
					LIMIT 100
					FOR UPDATE SKIP LOCKED
			)
			RETURNING id, event_type, payload, status, idempotency_key, priority, producer, created_at, updated_at
	`
	rows, err := o.repo.Query(ctx, stmt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			o.logger.Info().Str("func", "GetOutBoxEventsForUpdate").Msg("no pending tasks....")
			return nil, nil
		}
		o.logger.Error().Err(err).Str("func", "GetOutBoxEventsForUpdate").Msgf("failed to execute statement %s", stmt)
		return nil, err
	}

	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.OutBoxEvent])
	if err != nil {
		o.logger.Error().Err(err).Str("func", "GetOutBoxEventsForUpdate").Msg("failed to collect events")
		return nil, err
	}

	return events, nil
}

func (o *OutboxRepo) IncrementRetryCount(ctx context.Context, id string, err_string string) error {
	stmt := `
		UPDATE outbox_events
		SET retry_count = retry_count + 1,
		    error_string = @err_string,
		    status = 'PENDING'
		WHERE id = @id
	`

	rows, err := o.repo.Exec(ctx, stmt, pgx.NamedArgs{
		"id":         id,
		"err_string": err_string,
	})
	if err != nil {
		o.logger.Error().Err(err).Str("func", "IncrementRetryCount").Msgf("failed to update %s for retry", id)
		return err
	}

	if rows.RowsAffected() == 0 {
		o.logger.Warn().Str("func", "IncrementRetryCount").Msg("no rows updated")
		return fmt.Errorf("failed to update %s for retry", id)
	}

	return nil
}

func (o *OutboxRepo) MarkEventProcessed(ctx context.Context, id string) error {
	stmt := `
			UPDATE outbox_events
			SET status = 'PROCESSED',
					locked_at = NOW()
			WHERE id = @id
	`
	rows, err := o.repo.Exec(ctx, stmt, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		o.logger.Error().Err(err).Msgf("failed to update %s to \"PROCESSED\"", id)
		return err
	}

	if rows.RowsAffected() == 0 {
		o.logger.Warn().Msg("no rows updated")
		return fmt.Errorf("failed to mark %s as processed", id)
	}

	return nil
}
