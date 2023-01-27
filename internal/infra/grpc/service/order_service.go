package service

import (
	"context"

	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/infra/grpc/pb"
	"github.com/leandrobraga/goexpert-desafio-cleanarch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewOrderService(
	createOrderUseCase usecase.CreateOrderUseCase,
	ListOrdersUseCase usecase.ListOrdersUseCase,
) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  ListOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	order := &pb.Order{Id: output.ID, Price: float32(output.Price), Tax: float32(output.Tax), FinalPrice: float32(output.FinalPrice)}
	return &pb.CreateOrderResponse{Order: order}, nil
}

func (s *OrderService) ListOrder(ctx context.Context, in *pb.ListOrderRequest) (*pb.ListOrderResponse, error) {
	dto := usecase.ListOrdersInputDTO{
		Limit: int(in.Limit),
		Page:  int(in.Page),
		Sort:  in.Sort,
	}
	output, err := s.ListOrdersUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	var orders []*pb.Order
	for _, orderDTO := range output {
		order := &pb.Order{
			Id:         orderDTO.ID,
			Price:      float32(orderDTO.Price),
			Tax:        float32(orderDTO.Tax),
			FinalPrice: float32(orderDTO.FinalPrice),
		}
		orders = append(orders, order)
	}
	return &pb.ListOrderResponse{Orders: orders}, nil
}
