package telegram_bot

import (
	"context"
	"fmt"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)
func Send(message string, image string) {
	fmt.Println("Telegram sending...")
	var channel_ids = []int64{
		-1002556120690,
	}
	

	opts := []bot.Option{
	
	}
	fmt.Print(os.Getenv("BOT_TOKEN"))
	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}
	fmt.Println(message, "connected to bot")

	for _, chat_id := range channel_ids {
		if image != "" {
			fmt.Printf("sending image to %d\n", chat_id)
			file, err := os.Open(image)
			if err != nil {
				fmt.Printf("Error opening image file: %v\n", err)
				continue
			}
			defer file.Close()
			b.SendPhoto(context.Background(), &bot.SendPhotoParams{
			ChatID: chat_id,
			Photo:  &models.InputFileUpload{Data: file},
			Caption: message,
			ParseMode: "Markdown",
		})
		}else {
			fmt.Printf("sending to %d\n", chat_id)

			b.SendMessage(context.Background(), &bot.SendMessageParams{
				ChatID: chat_id,
				Text:   message,
				ParseMode: "MarkdownV2",
			})
		}
	}

	fmt.Println("Message sent")
}
