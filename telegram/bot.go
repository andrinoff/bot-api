package telegram_bot

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Send(ctx context.Context, message string, image io.Reader) error {
	fmt.Println("Telegram sending...")
	var channel_ids = []int64{
		-1002556120690,
	}

	opts := []bot.Option{}
	fmt.Print(os.Getenv("BOT_TOKEN"))
	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		return fmt.Errorf("error creating bot: %w", err)
	}

	defer b.Close(ctx)
	fmt.Println(message, "connected to bot")

	for _, chat_id := range channel_ids {
		if image != nil {
			fmt.Printf("sending image to %d\n", chat_id)

			if _, err := b.SendPhoto(ctx, &bot.SendPhotoParams{
				ChatID:    chat_id,
				Photo:     &models.InputFileUpload{Data: image},
				Caption:   message,
				ParseMode: "Markdown",
			}); err != nil {
				return fmt.Errorf("error sending photo: %w", err)
			}
		} else {
			fmt.Printf("sending to %d\n", chat_id)

			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chat_id,
				Text:      message,
				ParseMode: "MarkdownV2",
			}); err != nil {
				return fmt.Errorf("error sending message: %w", err)
			}
		}
	}

	fmt.Println("Message sent")

	return nil
}
func SendNoImage(ctx context.Context, message string) {
	fmt.Println("Telegram sending...")
	var channel_ids = []int64{
		-1002556120690,
	}
	opts := []bot.Option{}
	b, err := bot.New(os.Getenv("BOT_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}
	for _, chat_id := range channel_ids {

		fmt.Printf("sending to %d\n", chat_id)

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chat_id,
			Text:      message,
			ParseMode: "Markdown",
		})
		if err != nil {
			fmt.Printf("error sending message: %s", err)
		}
	}
	fmt.Println("Message sent")
}
