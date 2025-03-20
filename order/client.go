package order

import (
	"context"
	"github.com/fabian-emmanuel/go-ms/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewOrderServiceClient(conn)
	return &Client{conn, client}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}

func (c *Client) CreateOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error) {
	var productsProto []*pb.OrderProduct

	for _, product := range products {
		productsProto = append(productsProto, &pb.OrderProduct{
			ProductId: product.ID,
			Quantity:  product.Quantity,
		})
	}

	resp, err := c.service.CreateOrder(ctx, &pb.CreateOrderRequest{
		AccountId:     accountId,
		OrderProducts: productsProto,
	})

	if err != nil {
		return nil, err
	}

	newOrder := resp.Order
	newOrderCreatedAt := time.Time{}
	err = newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:          newOrder.Id,
		CreatedAt:   newOrderCreatedAt,
		AccountId:   newOrder.AccountId,
		TotalAmount: newOrder.TotalAmount,
		Products:    products,
	}, nil

}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error) {
	resp, err := c.service.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{
		AccountId: accountId,
	})

	if err != nil {
		return nil, err
	}

	var orders []*Order

	for _, op := range resp.Orders {
		newOrder := &Order{
			ID:          op.Id,
			AccountId:   op.AccountId,
			TotalAmount: op.TotalAmount,
		}

		newOrderCreatedAt := time.Time{}
		err = newOrderCreatedAt.UnmarshalBinary(op.CreatedAt)
		if err != nil {
			return nil, err
		}
		newOrder.CreatedAt = newOrderCreatedAt

		var products []OrderedProduct

		for _, p := range op.OrderedProducts {
			products = append(products, OrderedProduct{
				ID:          accountId,
				Quantity:    p.Quantity,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})

		}

		newOrder.Products = products
		orders = append(orders, newOrder)
	}

	return orders, nil

}
