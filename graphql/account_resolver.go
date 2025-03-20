package main

import (
	"context"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrdersForAccount(ctx, obj.ID)

	if err != nil {
		return nil, err
	}

	var orders []*Order

	for _, order := range orderList {
		var products []*OrderedProduct

		for _, product := range order.Products {
			products = append(products, &OrderedProduct{
				ID:          product.ID,
				Name:        product.Name,
				Price:       product.Price,
				Description: product.Description,
				Quantity:    int(product.Quantity),
			})
		}

		orders = append(orders, &Order{
			ID:         order.ID,
			CreatedAt:  order.CreatedAt,
			TotalPrice: order.TotalAmount,
			Products:   products,
		})
	}

	return orders, nil
}
