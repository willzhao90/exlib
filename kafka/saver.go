package kafka

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/sdce/exlib/kafka/topic"
	pb "gitlab.com/sdce/protogo"
)

// Saver is used for save unhandled messages
type Saver interface {
	Save(*sarama.ConsumerMessage, error) error
}

type defaultSaver struct {
	producer sarama.SyncProducer
	client   sarama.Client
	prefix   string
}

func (ds *defaultSaver) Save(msg *sarama.ConsumerMessage, err error) error {
	// todo: need implement
	rec := &pb.FailedRecord{
		Topic:     msg.Topic,
		Value:     msg.Value,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Reason:    err.Error(),
		Time:      msg.Timestamp.UnixNano(),
	}
	data, err := rec.Marshal()
	if err != nil {
		return err
	}
	newTopic := ds.prefix + msg.Topic
	message := &sarama.ProducerMessage{
		Topic: newTopic,
		Value: sarama.ByteEncoder(data),
	}
	log.Infof("Save a topic %v, reason %v", newTopic, rec.Reason)
	_, _, err = ds.producer.SendMessage(message)
	return err
}

// NewDefaultSaver return default saver instance
func NewDefaultSaver(kafkaClient *KafkaClient, optPrefix *string) Saver {
	producer, err := kafkaClient.GetProducer()
	if err != nil {
		log.Fatalf("Failed to get producer: %v", err)
	}
	prefix := topic.SavedTopicPrefix
	if optPrefix != nil {
		prefix = *optPrefix
	}
	return &defaultSaver{
		producer: producer.SyncProducer,
		client:   kafkaClient.Client,
		prefix:   prefix,
	}
}
