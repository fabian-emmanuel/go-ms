package catalog

import (
	"context"
	"github.com/segmentio/ksuid"
)

type Service interface {
	CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProductById(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error)
	GetProductsWithIds(ctx context.Context, ids []string, skip, take uint64) ([]*Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error)
}

type catalogService struct {
	repository Repository
}

func NewCatalogService(repository Repository) Service {
	return &catalogService{repository}
}

func (s *catalogService) CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	product := &Product{
		Name:        name,
		Description: description,
		ID:          ksuid.New().String(),
		Price:       price,
	}

	if err := s.repository.CreateProduct(ctx, *product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *catalogService) GetProductById(ctx context.Context, id string) (*Product, error) {
	return s.repository.GetProductById(ctx, id)
}

func (s *catalogService) GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.GetProducts(ctx, skip, take)
}

func (s *catalogService) GetProductsWithIds(ctx context.Context, ids []string, skip, take uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.GetProductsWithIds(ctx, ids, skip, take)
}

func (s *catalogService) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.SearchProducts(ctx, query, skip, take)
}
