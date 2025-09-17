package bot

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	"antiHate/config"
	"antiHate/nlp_service/nlp_pb"
)

// Start запускает Telegram-бота
func Start(cfg *config.Config, db *sqlx.DB) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании бота: %v", err)
	}
	log.Printf("🤖 Бот запущен: %s", bot.Self.UserName)

	conn, err := grpc.Dial("nlp_service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("❌ Ошибка при подключении к NLP gRPC: %v", err)
	}
	defer conn.Close()
	client := nlp_pb.NewNLPServiceClient(conn)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := client.Analyze(ctx, &nlp_pb.AnalyzeRequest{
			Text: update.Message.Text,
		})
		cancel()
		if err != nil {
			log.Printf("⚠️ Ошибка NLP gRPC: %v", err)
			continue
		}

		if resp.Score > 0.7 {
			// удаляем сообщение
			deleteConfig := tgbotapi.DeleteMessageConfig{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.MessageID,
			}
			if _, err := bot.Request(deleteConfig); err != nil {
				log.Printf("⚠️ Ошибка при удалении сообщения: %v", err)
			}

			// сохраняем в БД как удалённое
			db.Exec(
				"INSERT INTO messages (chat_id, user_id, message, deleted, created_at) VALUES ($1,$2,$3,$4,NOW())",
				update.Message.Chat.ID,
				update.Message.From.ID,
				update.Message.Text,
				true,
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"🚫 Сообщение нарушало правила и было удалено")
			bot.Send(msg)
			continue
		}

		// чистое сообщение → сохраняем
		db.Exec(
			"INSERT INTO messages (chat_id, user_id, message, deleted, created_at) VALUES ($1,$2,$3,$4,NOW())",
			update.Message.Chat.ID,
			update.Message.From.ID,
			update.Message.Text,
			false,
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Принято: "+update.Message.Text)
		bot.Send(msg)
	}
}
