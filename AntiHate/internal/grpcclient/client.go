package grpcclient

import (
	"context"
	"time"

	pb "antiHate/proto/nlp" // путь на сгенерированный код из nlp.proto

	"google.golang.org/grpc"
)

// NLPClient — обёртка над gRPC-клиентом
type NLPClient struct {
	conn   *grpc.ClientConn
	client pb.NLPServiceClient
}

// NewClient создаёт подключение к NLP-сервису
func NewClient(addr string) (*NLPClient, error) {
	// В реальном продакшене лучше использовать WithTransportCredentials (TLS),
	// но для локального запуска и MVP оставляем Insecure.
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &NLPClient{
		conn:   conn,
		client: pb.NewNLPServiceClient(conn),
	}, nil
}

// Close — закрыть соединение
func (c *NLPClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// Analyze — отправляет текст на проверку NLP-сервису
func (c *NLPClient) Analyze(text string) (*pb.AnalyzeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.Analyze(ctx, &pb.AnalyzeRequest{Text: text})
}
