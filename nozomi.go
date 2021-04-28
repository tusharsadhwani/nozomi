package main

import (
	"fmt"
	"log"
	"os"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

var HELP_MSG = `
Hi, I'm Nozomi.
I re-send any media that was forwarded from another channnel.

Just add me to your group, make me admin (to allow deleting the forwards), and I'll do my work.
`

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		fmt.Println("Unable to read bot token. Make sure you export $TOKEN in the environment.")
		os.Exit(1)
	}

	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalln(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.ForwardFromChat == nil {
			continue
		}

		msg := tgbot.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "html"
		msg.ReplyToMessageID = update.Message.MessageID

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = HELP_MSG
			case "help":
				msg.Text = HELP_MSG
			}
			bot.Send(msg)
		}

		if update.Message.Photo != nil {
			photoSizes := *update.Message.Photo
			if len(photoSizes) > 0 {
				photoMsg := tgbot.NewPhotoShare(
					update.Message.Chat.ID,
					photoSizes[len(photoSizes)-1].FileID,
				)
				_, err := bot.Send(photoMsg)
				if err == nil {
					bot.DeleteMessage(
						tgbot.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID),
					)
				}
			}
		}

		if update.Message.Video != nil {
			videoMsg := tgbot.NewVideoShare(
				update.Message.Chat.ID,
				update.Message.Video.FileID,
			)
			_, err := bot.Send(videoMsg)
			if err == nil {
				bot.DeleteMessage(
					tgbot.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID),
				)
			}
		}

		if update.Message.Animation != nil {
			gifMsg := tgbot.NewAnimationShare(
				update.Message.Chat.ID,
				update.Message.Animation.FileID,
			)
			_, err := bot.Send(gifMsg)
			if err == nil {
				bot.DeleteMessage(
					tgbot.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID),
				)
			}
		}
	}
}
