package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/tittuvarghese/core/logger"
	"github.com/tittuvarghese/order-management-service/core/database"
	"github.com/tittuvarghese/order-management-service/models"
	"github.com/tittuvarghese/order-management-service/proto"
	"github.com/tittuvarghese/order-management-service/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	proto.UnimplementedOrderServiceServer
	GrpcServer  *grpc.Server
	RdbInstance *database.RelationalDatabase
}

var log = logger.NewLogger("order-management-service")

func NewGrpcServer() *Server {
	return &Server{GrpcServer: grpc.NewServer()}
}

func (s *Server) Run(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error("Failed to listen", err)
	}

	proto.RegisterOrderServiceServer(s.GrpcServer, s)

	// Register reflection service on gRPC server
	reflection.Register(s.GrpcServer)
	log.Info("GRPC server is listening on port " + port)
	if err := s.GrpcServer.Serve(lis); err != nil {
		log.Error("failed to serve", err)
	}
}

func (s *Server) mustEmbedUnimplementedAuthServiceServer() {
	log.Error("implement me", nil)
}

func (s *Server) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	var order models.Order
	var totalPrice float64

	for _, item := range req.Items {
		totalPrice += item.Price * float64(item.Quantity)
		orderedItem := models.OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
		order.Items = append(order.Items, orderedItem)
	}
	order.TotalPrice = totalPrice

	buyerId, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return &proto.CreateOrderResponse{
			Message: "Unable to parse buyer id",
		}, err
	}
	order.CustomerID = buyerId
	order.Address = models.Address{
		AddressLine1: req.Address.AddressLine1,
		AddressLine2: req.Address.AddressLine2,
		City:         req.Address.City,
		State:        req.Address.State,
		Zip:          req.Address.Zip,
		Country:      req.Address.Country,
	}
	order.Phone = req.Phone

	err = service.CreateOrder(order, s.RdbInstance)

	if err != nil {
		return &proto.CreateOrderResponse{
			Message: "Failed to create the order. error: " + err.Error(),
		}, err
	}

	// Return the created product
	return &proto.CreateOrderResponse{Message: "Successfully created the order"}, nil
}
