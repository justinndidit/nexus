package broker

import "context"

type KafkaProducer struct {
}

func NewKafkaProducer() *KafkaProducer {
	return &KafkaProducer{}
}

// TODO:Implement kafka producer
func (k *KafkaProducer) Publish(ctx context.Context, topic, key string, payload []byte) error {
	return nil
}
