package main

import "context"

// ShopService 服务结构体
type ShopService struct {
	xxx interface{}
}

// NewShopService 依赖反转
func NewShopService(x interface{}) interface{} {
	return &ShopService{xxx: x}
}

// CreateOrder 实现业务逻辑
func (svr *ShopService) CreateOrder(ctx context.Context, r interface{}) (interface{}, error) {
	return nil, nil
}

func main() {

}
