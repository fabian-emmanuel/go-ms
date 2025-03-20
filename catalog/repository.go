package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"io"
	"log"
)

type Repository interface {
	Close()
	CreateProduct(ctx context.Context, product Product) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error)
	GetProductsWithIds(ctx context.Context, ids []string, skip, take uint64) ([]*Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(elastic.Config{Addresses: []string{url}})

	if err != nil {
		return nil, err
	}

	return &elasticRepository{client}, nil
}

func (r *elasticRepository) Close() {
	log.Println("Elasticsearch client does not require explicit closing")
}

func (r *elasticRepository) CreateProduct(ctx context.Context, product Product) error {
	body, err := json.Marshal(productDocument{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	})

	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      "catalog",
		DocumentID: product.ID,
		Body:       bytes.NewReader(body),
	}

	res, err := req.Do(ctx, r.client)

	if err != nil {
		return fmt.Errorf("failed to index product: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body")
		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("error indexing product: %s", res.String())
	}

	return nil
}

func (r *elasticRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	req := esapi.GetRequest{
		Index:      "catalog",
		DocumentID: id,
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := res.Body.Close()
		if err != nil {
			log.Printf("failed to close response body")
		}
	}(res.Body)

	if res.IsError() {
		return nil, fmt.Errorf("error getting product: %s", res.String())
	}
	var product Product
	if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product: %w", err)
	}
	return &product, err
}

func (r *elasticRepository) GetProducts(ctx context.Context, skip, take uint64) ([]*Product, error) {
	// Construct the search query with pagination
	query := map[string]interface{}{
		"from": skip,
		"size": take,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	// Convert the query to JSON
	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Create the search request
	req := esapi.SearchRequest{
		Index: []string{"catalog"},
		Body:  bytes.NewReader(body),
	}

	// Execute the request
	res, err := req.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body")
		}
	}(res.Body)

	// Check for response errors
	if res.IsError() {
		return nil, fmt.Errorf("error searching products: %s", res.String())
	}

	// Parse the response
	var sr searchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	// Convert the results to []*Product
	var products []*Product
	for _, hit := range sr.Hits.Hits {
		product := hit.Source
		products = append(products, &product)
	}

	return products, err
}

func (r *elasticRepository) GetProductsWithIds(ctx context.Context, ids []string, skip, take uint64) ([]*Product, error) {
	// Construct the query to filter by IDs
	query := map[string]interface{}{
		"from": skip,
		"size": take,
		"query": map[string]interface{}{
			"ids": map[string]interface{}{
				"values": ids,
			},
		},
	}

	// Convert query to JSON
	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// Create the search request
	req := esapi.SearchRequest{
		Index: []string{"catalog"},
		Body:  bytes.NewReader(body),
	}

	// Execute the request
	res, err := req.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("failed to search products by IDs: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body")
		}
	}(res.Body)

	// Check for response errors
	if res.IsError() {
		return nil, fmt.Errorf("error searching products by IDs: %s", res.String())
	}

	// Parse response
	var sr searchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	// Convert results to []*Product
	var products []*Product
	for _, hit := range sr.Hits.Hits {
		products = append(products, &hit.Source)
	}

	return products, err
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]*Product, error) {
	// Construct a match query for full-text search
	searchQuery := map[string]interface{}{
		"from": skip,
		"size": take,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name", "description"}, // Searching in name and description fields
			},
		},
	}

	// Convert query to JSON
	body, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search query: %w", err)
	}

	// Create the search request
	req := esapi.SearchRequest{
		Index: []string{"catalog"},
		Body:  bytes.NewReader(body),
	}

	// Execute the request
	res, err := req.Do(ctx, r.client)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body")
		}
	}(res.Body)

	// Check for response errors
	if res.IsError() {
		return nil, fmt.Errorf("error executing search: %s", res.String())
	}

	// Parse response
	var sr searchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	// Convert results to []*Product
	var products []*Product
	for _, hit := range sr.Hits.Hits {
		products = append(products, &hit.Source)
	}

	return products, err
}
