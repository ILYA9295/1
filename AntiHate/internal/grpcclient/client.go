package grpcclient

import (
	"context"
	"time"

	pb "antiHate/nlp_service"

	"google.golang.org/grpc"
)

type NLPClient struct {
	conn   *grpc.ClientConn
	client pb.NLPServiceClient
}

func NewClient(addr string) (*NLPClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &NLPClient{
		conn:   conn,
		client: pb.NewNLPServiceClient(conn),
	}, nil
}

func (c *NLPClient) Close() {
	c.conn.Close()
}

func (c *NLPClient) Analyze(text string) (*pb.AnalyzeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Analyze(ctx, &pb.AnalyzeRequest{Text: text})
}
