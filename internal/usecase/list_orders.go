package usecase

import "github.com/leandrobraga/goexpert-desafio-cleanarch/internal/entity"

type ListOrdersInputDTO struct {
	Page  int
	Limit int
	Sort  string
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(orderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (l *ListOrdersUseCase) Execute(input ListOrdersInputDTO) ([]OrderOutputDTO, error) {
	orders, err := l.OrderRepository.List(input.Page, input.Limit, input.Sort)
	if err != nil {
		return nil, err
	}
	var ordersDTO []OrderOutputDTO
	for _, order := range orders {
		orderDTO := OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		ordersDTO = append(ordersDTO, orderDTO)
	}
	return ordersDTO, nil
}
