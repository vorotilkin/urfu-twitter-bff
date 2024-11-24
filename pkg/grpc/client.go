package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Address string
}

type Client struct {
	config Config
	conn   *grpc.ClientConn
}

func (c *Client) OnStart(_ context.Context) error {
	conn, err := grpc.Dial(c.config.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *Client) OnStop(_ context.Context) error {
	return c.conn.Close()
}

func (c *Client) Connection() *grpc.ClientConn {
	return c.conn
}

func NewClient(c Config) *Client {
	return &Client{
		config: c,
	}
}
