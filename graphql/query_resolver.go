package main

import (
	"context"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	if id != nil {
		res, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			return nil, err
		}

		return []*Account{{
			ID:   res.ID,
			Name: res.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(10) // Default values

	if pagination != nil {
		skip, take = pagination.bounds()
	}

	res, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	var accounts []*Account
	for _, a := range res {
		accounts = append(accounts, &Account{
			ID:   a.ID,
			Name: a.Name,
		})
	}

	return accounts, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	if id != nil {
		res, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			return nil, err
		}

		return []*Product{{
			ID:          res.ID,
			Name:        res.Name,
			Description: res.Description,
			Price:       res.Price,
		}}, nil
	}

	skip, take := uint64(0), uint64(10) // Default values

	if pagination != nil {
		skip, take = pagination.bounds()
	}

	if query != nil {
		res, err := r.server.catalogClient.SearchProducts(ctx, *query, skip, take)
		if err != nil {
			return nil, err
		}
		var products []*Product
		for _, p := range res {
			products = append(products, &Product{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
		return products, nil
	}

	res, err := r.server.catalogClient.GetProducts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	var products []*Product
	for _, p := range res {
		products = append(products, &Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return products, nil
}

func (p *PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(0)

	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}

	if p.Take != nil {
		takeValue = uint64(*p.Take)
	}
	return skipValue, takeValue

}
