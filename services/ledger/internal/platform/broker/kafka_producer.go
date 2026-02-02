package broker

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	kafka "github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	logger *zerolog.Logger
}

func NewKafkaProducer(brokers []string, topic string, logger *zerolog.Logger) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			MaxAttempts:  5,
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			RequiredAcks: kafka.RequiredAcks(1),
			Async:        false,
			Compression:  kafka.Snappy,
		},
		logger: logger,
	}
}

func (p *KafkaProducer) PublishBatch(ctx context.Context, messages []kafka.Message) error {
	start := time.Now()

	err := p.writer.WriteMessages(ctx, messages...)
	if err != nil {
		p.logger.Error().Err(err).Msg("failed to publish batch to kafka")
		return fmt.Errorf("kafka publish error: %w", err)
	}

	p.logger.Info().
		Int("count", len(messages)).
		Dur("latency_ms", time.Since(start)).
		Msg("successfully published batch")

	return nil
}

func (p *KafkaProducer) Publish(ctx context.Context, topic string, key string, message []byte) error {
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
