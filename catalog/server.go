//go:generate protoc --go_out=./pb --go-grpc_out=./pb catalog.proto

package catalog

import (
	"context"
	"fmt"
	"github.com/fabian-emmanuel/go-ms/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type grpcServer struct {
	service Service
	pb.UnimplementedCatalogServiceServer
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterCatalogServiceServer(serv, &grpcServer{s, pb.UnimplementedCatalogServiceServer{}})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product, err := s.service.CreateProduct(ctx, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProductResponse{Product: &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}}, err
}

func (s *grpcServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := s.service.GetProductById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetProductResponse{Product: &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}}, err
}

func (s *grpcServer) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	response, err := s.service.GetProducts(ctx, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}

	var products []*pb.Product
	for _, product := range response {
		products = append(products, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return &pb.GetProductsResponse{Products: products}, err
}

func (s *grpcServer) GetProductsWithIds(ctx context.Context, req *pb.GetProductsWithIdsRequest) (*pb.GetProductsResponse, error) {
	response, err := s.service.GetProductsWithIds(ctx, req.Ids, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}

	var products []*pb.Product
	for _, product := range response {
		products = append(products, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return &pb.GetProductsResponse{Products: products}, err

}

func (s *grpcServer) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.GetProductsResponse, error) {
	response, err := s.service.SearchProducts(ctx, req.Query, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}

	var products []*pb.Product
	for _, product := range response {
		products = append(products, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return &pb.GetProductsResponse{Products: products}, err
}
