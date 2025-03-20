//go:generate protoc --go_out=./pb --go-grpc_out=./pb order.proto
package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/fabian-emmanuel/go-ms/account"
	"github.com/fabian-emmanuel/go-ms/catalog"
	"github.com/fabian-emmanuel/go-ms/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	pb.UnimplementedOrderServiceServer
}

func ListenGRPC(s Service, accountServiceUrl, catalogServiceUrl string, port int) error {
	accountClient, err := account.NewClient(accountServiceUrl)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogServiceUrl)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{s, accountClient, catalogClient, pb.UnimplementedOrderServiceServer{}})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, req.AccountId)
	if err != nil {
		log.Println("Error getting account: ", err)
		return nil, errors.New("account not found")
	}

	var productIds []string
	orderedProducts, err := s.catalogClient.GetProductsByIds(ctx, productIds, 0, 0)
	if err != nil {
		log.Println("Error getting products: ", err)
		return nil, errors.New("products not found")
	}

	var products []OrderedProduct
	for _, p := range orderedProducts {
		product := OrderedProduct{
			ID:          p.ID,
			Quantity:    0,
			Price:       p.Price,
			Name:        p.Name,
			Description: p.Description,
		}

		for _, rp := range req.OrderProducts {
			if rp.ProductId == p.ID {
				product.Quantity = rp.Quantity
				break
			}
		}

		if product.Quantity != 0 {
			products = append(products, product)
		}

	}

	order, err := s.service.CreateOrder(ctx, req.AccountId, products)

	if err != nil {
		log.Println("Error creating order: ", err)
		return nil, errors.New("error creating order")
	}

	orderProto := &pb.Order{
		Id:              order.ID,
		AccountId:       order.AccountId,
		TotalAmount:     order.TotalAmount,
		OrderedProducts: []*pb.OrderedProduct{},
	}

	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	for _, p := range order.Products {
		orderProto.OrderedProducts = append(orderProto.OrderedProducts, &pb.OrderedProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Quantity:    p.Quantity,
			Price:       p.Price,
		})
	}

	return &pb.CreateOrderResponse{Order: orderProto}, nil

}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, req *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, req.AccountId)
	if err != nil {
		log.Println("Error getting account: ", err)
		return nil, errors.New("account not found")
	}

	accountOrders, err := s.service.GetOrdersForAccount(ctx, req.AccountId)
	if err != nil {
		log.Println("Error getting orders: ", err)
		return nil, errors.New("error getting orders")
	}

	productIdMap := map[string]bool{}
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIdMap[p.ID] = true
		}
	}

	var productIds []string
	for id := range productIdMap {
		productIds = append(productIds, id)
	}

	products, err := s.catalogClient.GetProductsByIds(ctx, productIds, 0, 0)
	if err != nil {
		log.Println("Error getting products: ", err)
		return nil, errors.New("products not found")
	}

	var orders []*pb.Order
	for _, o := range accountOrders {
		op := &pb.Order{
			Id:              o.ID,
			AccountId:       o.AccountId,
			TotalAmount:     o.TotalAmount,
			OrderedProducts: []*pb.OrderedProduct{},
		}

		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		for _, product := range o.Products {
			for _, p := range products {
				if p.ID == product.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
				}
			}
			op.OrderedProducts = append(op.OrderedProducts, &pb.OrderedProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Quantity:    product.Quantity,
				Price:       product.Price,
			})
		}

		orders = append(orders, op)
	}

	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}
