package main

import (
	"context"
	"log"
	"velox-raptor-saham/include/bot"
	"velox-raptor-saham/include/queue"
	"velox-raptor-saham/include/storage"
	"velox-raptor-saham/include/util"

	"github.com/mymmrac/telego"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	util.InitLogger()
	util.Logger.Println("Starting Trading Bot...")

	// 2. Init Database
	if err := storage.InitDB(viper.GetViper()); err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	util.Logger.Println("Database Connected")

	// 3. Init Telegram Bot
	botToken := viper.GetString("telegram.bot_token")
	telegoBot, err := bot.InitBot(botToken)
	if err != nil {
		log.Fatalf("Failed to init bot: %v", err)
	}

	worker := queue.NewWorker(
		telegoBot,
		viper.GetInt("limit.bot_message_per_second"),
		viper.GetInt("limit.per_chat_delay_ms"),
	)
	worker.Start()
	util.Logger.Println("Worker Started")

	// 5. Long Polling Loop
	updatesChan := make(chan telego.Update, 100)

	// Goroutine Fetcher
	go func() {
		params := telego.GetUpdatesParams{
			Timeout: viper.GetInt("telegram.polling_timeout"),
		}

		ctx := context.Background()
		for {
			updates, err := telegoBot.GetUpdates(ctx, &params)
			if err != nil {
				util.Logger.Printf("Error getting updates: %v", err)
				continue
			}
			for _, update := range updates {
				if update.UpdateID >= params.Offset {
					params.Offset = update.UpdateID + 1
					updatesChan <- update
				}
			}
		}
	}()

	// Main Event Loop
	util.Logger.Println("Bot is running...")
	for update := range updatesChan {
		bot.HandleUpdate(update, worker)
	}
}
