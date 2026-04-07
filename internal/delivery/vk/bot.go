package vk

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

type Bot struct {
	vk     *api.VK
	lp     *longpoll.LongPoll
	router *Router
}

func NewBot(token string, groupID int) (*Bot, error) {

	vk := api.NewVK(token)

	lp, err := longpoll.NewLongPoll(vk, groupID)
	if err != nil {
		return nil, err
	}

	return &Bot{
		vk:     vk,
		lp:     lp,

	}, nil
}

func (b *Bot) SetRouter(r *Router){
	b.router = r
}

func (b *Bot) Start() {
	log.Println("VK bot started!")

	b.lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		msg := obj.Message

		// защита от пустых сообщений
		if msg.Text == "" {
			return
		}

		log.Printf("IN [%d]: %s", msg.FromID, msg.Text)

		// ВАЖНО: теперь идёт в router
		b.router.Handle(msg.FromID, msg.Text)
	})

	if err := b.lp.Run(); err != nil {
		log.Fatal(err)
	}
}

func (b *Bot) sendMessage(userID int, text string) {

	rand.Seed(time.Now().UnixNano())

	_, err := b.vk.MessagesSend(api.Params{
		"user_id":   userID,
		"message":   text,
		"random_id": rand.Int(),
	})

	if err != nil {
		log.Println("send error:", err)
	}
}
