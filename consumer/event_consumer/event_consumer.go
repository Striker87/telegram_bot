package event_consumer

import (
	"log"
	"time"

	"github.com/Striker87/telegram_bot/clients/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int // сколько событий мы будем обрабатывать за раз
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %v", err)
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		err = c.handleEvents(gotEvents)
		if err != nil {
			log.Printf("failed to handleEvents %v due error: %v", gotEvents, err)
			continue
		}
	}
}

func (c Consumer) handleEvents(events []events.Event) error {
	// todo: make async
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("failed to handle event: %v", err)
			continue
		}
	}

	return nil
}
