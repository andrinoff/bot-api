package discord_bot

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var BotToken string = os.Getenv("DISCORD_TOKEN")

func Send(ctx context.Context, message string, image io.Reader) error {
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	if err = discord.Open(); err != nil {
		return fmt.Errorf("error opening Discord connection: %w", err)
	}

	defer discord.Close()

	_, err = discord.ChannelMessageSendComplex("1383195767623127111", &discordgo.MessageSend{
		Content: message,
		Files: []*discordgo.File{
			{
				Name:   "image.png",
				Reader: image,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	fmt.Println("Message sent")

	return nil
}

func SendNoImage(ctx context.Context, message string) {
	fmt.Println("Discord sending...")
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}
	discord.Open()

	defer discord.Close()
	fmt.Println(message)
	_, err = discord.ChannelMessageSend("1383195767623127111", message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	fmt.Println("Message sent")

}
