package main

type Order struct {
	Item string
}

type OrderRepo interface {
	SaveOrder(*Order)
}

func NewOrderUseCase(repo OrderRepo) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

type OrderUseCase struct {
	repo OrderRepo
}

func (uc *OrderUseCase) Buy(o *Order) {
	// logic business
	uc.repo.SaveOrder(o)
}

func main() {

}
