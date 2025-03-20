package account

import (
	"context"
	"github.com/fabian-emmanuel/go-ms/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {

	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	client := pb.NewAccountServiceClient(conn)
	return &Client{conn, client}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	r, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})

	if err != nil {
		return nil, err
	}

	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip, take uint64) ([]*Account, error) {
	r, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: skip, Take: take})
	if err != nil {
		return nil, err
	}
	var accounts []*Account
	for _, a := range r.Accounts {
		accounts = append(accounts, &Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}
	return accounts, nil
}
