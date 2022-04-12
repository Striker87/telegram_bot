package telegram

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Striker87/telegram_bot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, userName string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command %q from %q", text, userName)

	// add page: http:// ...
	// rnd page: /rnd
	// help: /help
	// start: /start: hi + help

	if isAddCmd(text) {
		return p.savePage(chatId, text, userName)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatId, userName)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.tg.SendMessage(chatId, msgUnknownCmd)
	}
}

func (p *Processor) savePage(chatId int, pageUrl string, userName string) error {
	page := &storage.Page{
		URL:      pageUrl,
		UserName: userName,
	}
	//sendMsg := NewMessageSender(chatId, p.tg)

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return fmt.Errorf("failed to do command due error: %v", err)
	}
	if isExists {
		//return sendMsg(msgAlreadyExists)
		return p.tg.SendMessage(chatId, msgAlreadyExists)
	}

	err = p.storage.Save(page)
	if err != nil {
		return fmt.Errorf("failed to save page %v due error: %v", page, err)
	}

	err = p.tg.SendMessage(chatId, msgSaved)
	if err != nil {
		return fmt.Errorf("failed to SendMessage chatId = %d, msgSaved = %s due error: %v", chatId, msgSaved, err)
	}

	return nil
}

func (p *Processor) sendRandom(chatId int, userName string) error {
	page, err := p.storage.PickRandom(userName)
	if err != nil && !errors.Is(err, storage.ErrNotSavedPages) {
		return fmt.Errorf("failed to PickRandom userName = %s due error: %v", userName, err)
	}

	if errors.Is(err, storage.ErrNotSavedPages) {
		return p.tg.SendMessage(chatId, msgNoSavedPages)
	}

	err = p.tg.SendMessage(chatId, page.URL)
	if err != nil {
		return fmt.Errorf("failed to SendMessage chatId = %d, page.URL = %s due error: %v", chatId, page.URL, err)
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatId int) error {
	return p.tg.SendMessage(chatId, msgHelp)
}

func (p *Processor) sendHello(chatId int) error {
	return p.tg.SendMessage(chatId, msgHello)
}

//func NewMessageSender(chatId int, tg *telegram.Client) func(string) error {
//	return func(msg string) error {
//		return tg.SendMessage(chatId, msg)
//	}
//}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
