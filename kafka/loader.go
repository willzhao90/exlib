package kafka

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	pb "gitlab.com/sdce/protogo"
)

// LoaderStateType is type return by loader for other entities to know its state
type LoaderStateType int

// LoaderNameType used in name send to handler
type LoaderNameType string

const (
	// LoaderStateNone is a initial state or exit state
	LoaderStateNone LoaderStateType = iota
	// LoaderStateIdle when loader has nothing to process
	LoaderStateIdle
	// LoaderStateBusy when loader has >=1 message to process.
	LoaderStateBusy
)

const (
	// DefLoaderName is the default loader name send to handler
	DefLoaderName LoaderNameType = "loader"
)

// Loader interface returned by New Loader
type Loader interface {
	// Run blocked function
	Run(handler Handler) error
	GetState() LoaderStateType
	Close()
}

type loader struct {
	consumer *Consumer

	startTime   time.Time
	stopTime    time.Time
	topicPrefix string
	msgPool     chan *sarama.ConsumerMessage
	poolSig     chan int
}

func (ld *loader) Run(handler Handler) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return ld.handleSavedMessages(ctx, handler)
}

func (ld *loader) GetState() LoaderStateType {
	if ld.startTime.After(ld.stopTime) {
		return LoaderStateBusy
	}
	if ld.stopTime.After(time.Now().Add(-30 * time.Second)) {
		return LoaderStateBusy
	} else {
		return LoaderStateIdle
	}
}

func (ld *loader) Close() {
	if ld.consumer != nil {
		ld.consumer.ConsumerGroup.Close()
	}
}

func (ld *loader) handleSavedMessages(ctx context.Context, handler Handler) (err error) {
	// go routine to dispatch msg
	go func() {
		saramaHandler := NewSaramaHandler(HandlerFunc(func(ctx context.Context, msg *sarama.ConsumerMessage) error {
			select {
			case <-ld.poolSig: // channel failed
				return fmt.Errorf("loader msg pool channel closed")
			default:
			}
			ld.msgPool <- msg
			return nil
		}), nil, true, nil, nil)

		ld.consumer.HandleMessages(ctx, nil, regexp.MustCompile(fmt.Sprintf("^%v", ld.topicPrefix)), nil, saramaHandler)
	}()

	// handle msg
	defer close(ld.poolSig)
	ctx = context.WithValue(ctx, DefLoaderName, ld)
	for {
		func() {
			select {
			case msg := <-ld.msgPool:
				ld.startTime = time.Now()
				defer func() { ld.stopTime = time.Now() }()
				record := new(pb.FailedRecord)
				err = record.Unmarshal(msg.Value)
				if err != nil {
					log.Errorf("cannot unmarshal record (Tp:%v, Pt:%v, Of:%v): %v",
						msg.Topic, msg.Partition, msg.Offset, err)
					break
				}

				log.Infof(">>>> Handling saved record at (TPt:%v Of:%v) -> (Err:%v, TPt:%v, Of:%v)", msg.Partition, msg.Offset, record.Reason, record.Partition, record.Offset)
				err = handler.Dispatch(ctx, &sarama.ConsumerMessage{
					Partition: record.Partition,
					Offset:    record.Offset,
					Topic:     record.Topic,
					Value:     record.Value,
				})
				if err != nil {
					err = fmt.Errorf("Failed to handle message %v: %v", msg.Topic, err)
					log.Println(err)
					break
				}
			}
		}()
	}
}

// NewDefaultLoader create a default loader
func NewDefaultLoader(kafkaClient *KafkaClient, consumerGroup, topicPrefix string) Loader {
	consumer, err := kafkaClient.GetConsumer(consumerGroup, nil)
	if err != nil {
		log.Fatalf("Failed to get kafka consumer: %v", err)
	}
	now := time.Now()
	return &loader{
		consumer:    consumer,
		topicPrefix: topicPrefix,
		startTime:   now,
		stopTime:    now,
		msgPool:     make(chan *sarama.ConsumerMessage),
		poolSig:     make(chan int),
	}
}
