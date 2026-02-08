package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
	"velox-raptor-saham/include/service"
	"velox-raptor-saham/include/util"

	"github.com/mymmrac/telego"
)

type Job struct {
	ChatID int64
	Symbol string
}

type Worker struct {
	Bot           *telego.Bot
	jobQueue      chan Job
	limiter       <-chan time.Time // Global limiter
	lastChatTimes map[int64]time.Time
	chatMutex     sync.Mutex
	delayPerChat  time.Duration
}

func NewWorker(bot *telego.Bot, msgPerSec int, delayMs int) *Worker {
	interval := time.Second / time.Duration(msgPerSec)
	limiter := time.NewTicker(interval).C

	return &Worker{
		Bot:           bot,
		jobQueue:      make(chan Job, 100),
		limiter:       limiter,
		lastChatTimes: make(map[int64]time.Time),
		delayPerChat:  time.Duration(delayMs) * time.Millisecond,
	}
}

func (w *Worker) Start() {
	go func() {
		for range w.limiter {
			select {
			case job := <-w.jobQueue:
				w.processJob(job)
			default:
				// Tidak ada job, skip tick
			}
		}
	}()
}

func (w *Worker) EnqueueJob(job Job) {
	w.jobQueue <- job
}

func (w *Worker) processJob(job Job) {
	w.chatMutex.Lock()
	lastSent, exists := w.lastChatTimes[job.ChatID]
	timeSinceLast := time.Since(lastSent)

	if exists && timeSinceLast < w.delayPerChat {
		time.Sleep(w.delayPerChat - timeSinceLast)
	}
	w.lastChatTimes[job.ChatID] = time.Now()
	w.chatMutex.Unlock()

	util.Logger.Printf("Processing %s for chat %d", job.Symbol, job.ChatID)

	result, err := service.AnalyzeStock(job.ChatID, job.Symbol)
	var messageText string

	if err != nil {
		messageText = fmt.Sprintf("❌ *Terjadi Kesalahan*\n\n`%s`", err.Error())
	} else {
		emoji := "⚪"
		if result.Signal == "BUY 🟢" {
			emoji = "🟢"
		} else if result.Signal == "SELL 🔴" {
			emoji = "🔴"
		}

		messageText = fmt.Sprintf(
			"📊 *ANALISIS SAHAM %s*\n\n"+
				"💰 Harga Saat Ini: `%.2f`\n"+
				"🔮 Prediksi Besok: `%.2f`\n"+
				"📈 Perubahan: `%.2f%%`\n\n"+
				"%s *Sinyal: %s*\n\n"+
				"🛡️ Support: `%.2f`\n"+
				"🚀 Resistance: `%.2f`\n\n"+
				"📉 *Confidence: %s*\n"+
				"📊 Interval: `%s`\n"+
				"⚡ RMSE Error: `%.2f`\n\n"+
				"_Powered by VeloxRaptor_",
			result.Symbol, result.CurrentPrice, result.Prediction, result.ChangePct,
			emoji, result.Signal, result.Support, result.Resistance,
			result.ConfidenceLevel, result.ConfidenceInterval, result.RMSE,
		)
	}

	ctx := context.Background()

	_, err = w.Bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: job.ChatID},
		Text:      messageText,
		ParseMode: telego.ModeMarkdown,
	})

	if err != nil {
		util.Logger.Printf("Gagal mengirim pesan: %v", err)
	}
}
