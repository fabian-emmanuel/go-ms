package order

import (
	"context"
	"github.com/segmentio/ksuid"
	"time"
)

type Service interface {
	CreateOrder(ctx context.Context, accountId string, orderedProducts []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error)
}

type orderService struct {
	repo Repository
}

func NewOrderService(repo Repository) Service {
	return &orderService{repo}
}

func (s *orderService) CreateOrder(ctx context.Context, accountId string, orderedProducts []OrderedProduct) (*Order, error) {
	order := &Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountId: accountId,
		Products:  orderedProducts,
	}

	order.TotalAmount = 0.0

	for _, orderedProduct := range orderedProducts {
		order.TotalAmount += orderedProduct.Price * float64(orderedProduct.Quantity)
	}

	err := s.repo.CreateOrder(ctx, *order)
	if err != nil {
		return nil, err
	}

	return order, nil

}

func (s *orderService) GetOrdersForAccount(ctx context.Context, accountId string) ([]*Order, error) {
	return s.repo.GetOrdersForAccount(ctx, accountId)
}
