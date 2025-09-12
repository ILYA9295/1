package bot

import (
	"context"
	"log"
	"time"

	"antiHate/config"
	"antiHate/nlp_service/nlp_pb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func Start(cfg *config.Config, db *sqlx.DB) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var conn *grpc.ClientConn
	for {
		conn, err = grpc.Dial("nlp_service:50051", grpc.WithInsecure())
		if err == nil {
			break
		}
		log.Println("Waiting for NLP service...")
		time.Sleep(5 * time.Second)
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
			log.Printf("NLP gRPC error: %v", err)
			continue
		}

		if resp.Score > 0.7 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 Ваше сообщение содержит недопустимую лексику или оскорбления")
			bot.Send(msg)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Принято: "+update.Message.Text)
		bot.Send(msg)
	}
}
