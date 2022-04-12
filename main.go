package main

import (
	"flag"
	"log"

	"github.com/Striker87/telegram_bot/clients/events/telegram"
	tgClient "github.com/Striker87/telegram_bot/clients/telegram"
	"github.com/Striker87/telegram_bot/consumer/event_consumer"
	"github.com/Striker87/telegram_bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	// fetcher := fetcher.New() // получает события
	// processor := processor.New() // обрабатывает события
	eventsProcessor := telegram.New(
		tgClient.New(mustToken(), tgBotHost),
		files.New(storagePath),
	)

	log.Println("service started")
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatalf("consumer failed due error: %v", err)
	}
}

func mustToken() string {
	// bot -tg-bot-token 'my-token'
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
