package kafka

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

const (
	// ConsumerStateFail mark the message with fail
	FAIL_METADATA = "{\"success\": false}"
	// ConsumerStateDone mark the message with done
	SUCCESS_METADATA = "{\"success\": true}"
)

// Handler that actually handles incoming messages
type Handler interface {
	Dispatch(ctx context.Context, msg *sarama.ConsumerMessage) error
}

// HandlerFunc single function handler
type HandlerFunc func(context.Context, *sarama.ConsumerMessage) error

// Dispatch performs dispatch
func (h HandlerFunc) Dispatch(ctx context.Context, msg *sarama.ConsumerMessage) error {
	return h(ctx, msg)
}

// ConsumerSessionFunc is defined for implementing setup or cleanup func
type ConsumerSessionFunc func(sarama.ConsumerGroupSession) error

// MsgOption used for additional options of msg
// every option occupy 1 bit
type MsgOption int32

const (
	// MarkMsgAfterProcessing share bit 0
	MarkMsgAfterProcessing MsgOption = 1
	// MarkMsgBeforeProcessing share bit 1
	MarkMsgBeforeProcessing MsgOption = 2
)

type saramaGroupHandler struct {
	handler  Handler
	saver    Saver
	errFatal bool                // assume errors returned from the dispatcher as fatal and return immediately
	setup    ConsumerSessionFunc // Can be nil
	cleanup  ConsumerSessionFunc // Can be nil
	option   MsgOption
}

// NewSaramaHandler Wraps our handler and implements sarama.ConsumerGroupHandler
func NewSaramaHandler(ourHandler Handler, saver Saver, errFatal bool, setup, cleanup ConsumerSessionFunc, opts ...MsgOption) sarama.ConsumerGroupHandler {
	option := MarkMsgAfterProcessing
	if len(opts) > 0 {
		option = 0
		for _, opt := range opts {
			option = option | opt
		}
	}
	return &saramaGroupHandler{ourHandler, saver, errFatal, setup, cleanup, option}
}

func (h saramaGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	if h.setup == nil {
		return nil
	}
	return h.setup(session)
}
func (h saramaGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	if h.cleanup == nil {
		return nil
	}
	return h.cleanup(session)
}

func (h saramaGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()
	tracer := opentracing.GlobalTracer()
	for msg := range claim.Messages() {
		err := func() error {
			metadata := SUCCESS_METADATA
			if h.option&MarkMsgBeforeProcessing != 0 {
				session.MarkMessage(msg, metadata)
			}
			defer func() {
				if h.option&MarkMsgAfterProcessing != 0 {
					session.MarkMessage(msg, metadata)
				}
			}()
			span := tracer.StartSpan("Kafka inflow- " + msg.Topic)
			span.SetTag("partition", msg.Partition)
			span.SetTag("offset", msg.Offset)
			span.SetTag("topic", msg.Topic)
			defer span.Finish()
			err := h.handler.Dispatch(opentracing.ContextWithSpan(ctx, span), msg)
			if err != nil {
				span.LogFields(otlog.Error(err))
				metadata = FAIL_METADATA
				err = fmt.Errorf("Failed to handle message %v.%v.%v: %v", msg.Topic, msg.Partition, msg.Offset, err)
				log.Error(err)
				if h.saver != nil {
					err = h.saver.Save(msg, err)
					span.LogFields(otlog.String("saver", fmt.Sprintf("saved to saver: %v", err)))
					if err != nil {
						log.Errorf("Failed to perform fail save on message %v.%v.%v: %v", msg.Topic, msg.Partition, msg.Offset, err)
					}
				}
				if h.errFatal {
					return err
				}
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
