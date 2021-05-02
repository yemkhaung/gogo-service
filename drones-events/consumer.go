package dronesevents

import (
	"log"
)

type eventConsumer interface {
	Dequeue() (<-chan interface{}, error)
}

type consumerRegistry struct {
	consumers  map[string]eventConsumer
	processors map[string]eventProcessor
}

func newConsumerRegistry() *consumerRegistry {
	return &consumerRegistry{
		consumers:  make(map[string]eventConsumer),
		processors: make(map[string]eventProcessor),
	}
}

func (reg *consumerRegistry) RegisterConsumer(eventType string, consumer eventConsumer) {
	reg.consumers[eventType] = consumer
}

func (reg *consumerRegistry) GetProcessor(eventType string) (processor eventProcessor) {
	processor, ok := reg.processors[eventType]
	if ok {
		return processor
	}
	log.Printf("No event-processor configured for eventType:%s", eventType)
	return nil
}

func (reg *consumerRegistry) RegisterProcessor(eventType string, processor eventProcessor) {
	reg.processors[eventType] = processor
}

func (reg *consumerRegistry) consumeEvents() error {
	for eventType, consumer := range reg.consumers {
		eventsQ, err := consumer.Dequeue()
		if err != nil {
			return err
		}
		processor := reg.GetProcessor(eventType)
		if processor == nil {
			continue
		}
		log.Printf("Started consuming (%s) events...", eventType)
		go func(msgs <-chan interface{}, proc eventProcessor) {
			for event := range msgs {
				err := proc.Process(event)
				if err != nil {
					log.Printf("Error processing event: %s", err)
				}
			}
		}(eventsQ, processor)
	}

	return nil
}
