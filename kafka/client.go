package kafka

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type KafkaClient struct {
	Client   sarama.Client
	consumer *Consumer
	producer *Producer
}

func NewKafka(c *Config) *KafkaClient {

	// init (custom) config, enable errors and notifications
	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_2
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	// config.Group.Return.Notifications = true

	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	config.Metadata.Full = true

	err := config.Validate()
	if err != nil {
		log.Fatalf("sarama config is invalid: %v", err)
	}

	// Start with a client
	client, err := sarama.NewClient(c.Brokers, config)
	if err != nil {
		log.Fatalf("Failed to initiate kafka client: %v. config=%v, %v", err, c, config)
	}

	out := &KafkaClient{Client: client}

	return out
}

func (k *KafkaClient) GetProducer() (*Producer, error) {
	if k.producer != nil {
		return k.producer, nil
	}

	syncProducer, err := sarama.NewSyncProducerFromClient(k.Client)
	if err != nil {
		return nil, fmt.Errorf("Failed to start Sarama sync producer: %v", err)
	}
	asyncProducer, err := sarama.NewAsyncProducerFromClient(k.Client)
	if err != nil {
		return nil, fmt.Errorf("Failed to start Sarama async producer: %v", err)
	}
	k.producer = &Producer{
		k, syncProducer, asyncProducer,
	}
	return k.producer, nil
}

func (k *KafkaClient) GetConsumer(groupID string, saver Saver) (*Consumer, error) {
	if k.consumer != nil {
		return k.consumer, nil
	}

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(groupID, k.Client)
	if err != nil {
		return nil, fmt.Errorf("Failed to start Sarama consumer group: %v", err)
	}

	k.consumer = &Consumer{
		ConsumerGroup: group,
		Client:        k,
	}
	// consume errors
	go func() {
		for err := range k.consumer.ConsumerGroup.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	return k.consumer, nil
}

func (k *KafkaClient) Close() {
	k.consumer.ConsumerGroup.Close()
	k.producer.Close()
	k.Client.Close()
}
