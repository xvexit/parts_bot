package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct{
	api *tgbotapi.BotAPI
	router *Router
}

func NewBot(token string, router *Router) (*Bot, error){
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil{
		return nil, err
	}

	api.Debug = false

	return &Bot{
		api: api,
		router: router,
	}, nil
}

func (b *Bot) Start(){
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates{
		if update.Message == nil{
			continue
		}

		b.router.Handle(b.api, update.Message)
	}
}