package catalog

import (
	"context"
	"github.com/fabian-emmanuel/go-ms/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewCatalogServiceClient(conn)
	return &Client{conn, client}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}

func (c *Client) CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	res, err := c.service.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, err
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	res, err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, err
}

func (c *Client) GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error) {
	res, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{Skip: skip, Take: take})
	if err != nil {
		return nil, err
	}

	var products []*Product
	for _, product := range res.Products {
		products = append(products, &Product{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return products, err
}

func (c *Client) GetProductsByIds(ctx context.Context, ids []string, skip, take uint64) ([]*Product, error) {
	res, err := c.service.GetProductsWithIds(ctx, &pb.GetProductsWithIdsRequest{Skip: skip, Take: take, Ids: ids})
	if err != nil {
		return nil, err
	}

	var products []*Product
	for _, product := range res.Products {
		products = append(products, &Product{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return products, err
}

func (c *Client) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error) {
	res, err := c.service.SearchProducts(ctx, &pb.SearchProductsRequest{Skip: skip, Take: take, Query: query})
	if err != nil {
		return nil, err
	}

	var products []*Product
	for _, product := range res.Products {
		products = append(products, &Product{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return products, err
}
