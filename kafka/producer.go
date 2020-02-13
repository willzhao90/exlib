package kafka

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

// Producer containts sync and async producers
type Producer struct {
	Client        *KafkaClient
	SyncProducer  sarama.SyncProducer
	AsyncProducer sarama.AsyncProducer
}

// Close closes up all producers.
// TODO not sure whether to close client as well.
func (p *Producer) Close() {
	if err := p.SyncProducer.Close(); err != nil {
		log.Println("Failed to shut down data collector cleanly", err)
	}

	if err := p.AsyncProducer.Close(); err != nil {
		log.Println("Failed to shut down access log producer cleanly", err)
	}
}
