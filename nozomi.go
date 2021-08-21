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

		msg := tgbot.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "html"
		msg.ReplyToMessageID = update.Message.MessageID

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "resend":
				resendMedia(bot, update.Message.ReplyToMessage)
				continue
			case "start":
				msg.Text = HELP_MSG
			case "help":
				msg.Text = HELP_MSG
			}
			bot.Send(msg)
		}

		if update.Message.ForwardFromChat == nil {
			continue
		}
		resendMedia(bot, update.Message)
	}
}

func resendMedia(bot *tgbot.BotAPI, message *tgbot.Message) {
	if message.Photo != nil {
		photoSizes := *message.Photo
		if len(photoSizes) > 0 {
			photoMsg := tgbot.NewPhotoShare(
				message.Chat.ID,
				photoSizes[len(photoSizes)-1].FileID,
			)
			_, err := bot.Send(photoMsg)
			if err == nil {
				bot.DeleteMessage(
					tgbot.NewDeleteMessage(message.Chat.ID, message.MessageID),
				)
			}
		}
	}

	if message.Video != nil {
		videoMsg := tgbot.NewVideoShare(
			message.Chat.ID,
			message.Video.FileID,
		)
		_, err := bot.Send(videoMsg)
		if err == nil {
			bot.DeleteMessage(
				tgbot.NewDeleteMessage(message.Chat.ID, message.MessageID),
			)
		}
	}

	if message.Animation != nil {
		gifMsg := tgbot.NewAnimationShare(
			message.Chat.ID,
			message.Animation.FileID,
		)
		_, err := bot.Send(gifMsg)
		if err == nil {
			bot.DeleteMessage(
				tgbot.NewDeleteMessage(message.Chat.ID, message.MessageID),
			)
		}
	}
}
