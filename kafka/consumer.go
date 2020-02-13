package kafka

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

// Consumer wrapper
type Consumer struct {
	ConsumerGroup sarama.ConsumerGroup
	Client        *KafkaClient
}

func (c *Consumer) GetTopics(topics []string, topicWhiteList, topicBlackList *regexp.Regexp) ([]string, error) {
	allTopics, err := c.Client.Client.Topics()
	if err != nil {
		return nil, fmt.Errorf("Failed to get topics: %v", err)
	}
	if topicWhiteList == nil && topicBlackList == nil {
		return topics, nil
	}
	for _, topic := range allTopics {
		if (topicWhiteList == nil || topicWhiteList.Match([]byte(topic))) &&
			(topicBlackList == nil || !topicBlackList.Match([]byte(topic))) {
			topics = append(topics, topic)
		}
	}
	return topics, nil
}

// HandleMessages registers handler function for Kafka message
// It should block the calling thread.
func (c *Consumer) HandleMessages(ctx context.Context, topics []string, topicWhiteList, topicBlackList *regexp.Regexp,
	handler sarama.ConsumerGroupHandler) {
	for {
		fullTopics, err := c.GetTopics(topics, topicWhiteList, topicBlackList)
		if err != nil {
			log.Fatalf("Failed to load topics to consume: %v", err)
		}
		if len(fullTopics) == 0 {
			// no related topic created for consuming
			log.Infof("no topics available, wait 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		// TODO some timeout mechanism
		err = c.ConsumerGroup.Consume(ctx, fullTopics, handler)
		if err != nil {
			log.Fatalf("Failed to consume: %v", err)
		}
	}
}
