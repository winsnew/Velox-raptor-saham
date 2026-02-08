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

	// Handle Command /start
	if text == "/start" {
		worker.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: telego.ChatID{ID: chatID},
			Text:   "👋 Halo! Saya JEM_B0T.\n\nKetik `/analisis KODE.SAHAM` untuk memprediksi harga.\nContoh: `/analisis BBCA.JK`",
		})
		return
	}

	if strings.HasPrefix(text, "/analisis") {
		parts := strings.Split(text, " ")
		if len(parts) < 2 {
			worker.Bot.SendMessage(ctx, &telego.SendMessageParams{
				ChatID: telego.ChatID{ID: chatID},
				Text:   "⚠️ Format salah. Gunakan: `/analisis KODE.SAHAM`",
			})
			return
		}
		symbol := strings.ToUpper(parts[1])

		// Kirim ke Worker Queue
		job := queue.Job{
			ChatID: chatID,
			Symbol: symbol,
		}
		worker.EnqueueJob(job)

		worker.Bot.SendMessage(ctx, &telego.SendMessageParams{
			ChatID: telego.ChatID{ID: chatID},
			Text:   "⏳ Sedang menganalisis data di engine Python... Mohon tunggu sebentar ( estimasi 10-30 detik ).",
		})
	}
}
