package bot

import (
	"context"
	"strings"
	"velox-raptor-saham/include/queue"

	"github.com/mymmrac/telego"
)

func HandleUpdate(update telego.Update, worker *queue.Worker) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	ctx := context.Background()

	if text == "/start" {
		worker.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: telego.ChatID{ID: chatID},
			Text:   "👋 HELLO! I'm JEM_B0T.\n\nMsg `/analisis KODE.SAHAM` for price predict.\nExample: `/analisis BBCA.JK`",
		})
		return
	}

	if strings.HasPrefix(text, "/analisis") {
		parts := strings.Split(text, " ")
		if len(parts) < 2 {
			worker.Bot.SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: chatID},
				Text:   "⚠️ Wrong Format. Use: `/analisis KODE.SAHAM`",
			})
			return
		}
		symbol := strings.ToUpper(parts[1])

		job := queue.Job{
			ChatID: chatID,
			Symbol: symbol,
		}
		worker.EnqueueJob(job)

		worker.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: telego.ChatID{ID: chatID},
			Text:   "Analyze Processing...",
		})
	}
}
