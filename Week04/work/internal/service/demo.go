package service

import (
	"context"
	"work/api/work/v1"
	"work/internal/biz"
)

type Hello struct {
	biz.HelloRequest
}

func NewHello(request biz.HelloRequest) *Hello {
	return &Hello{request}
}

func (h *Hello) SayHello(ctx context.Context, request *work_v1.HelloRequest) (*work_v1.HelloResponse, error) {
	bizRes := h.BizHello(ctx, biz.NewHelloReq(request.Name))
	return &work_v1.HelloResponse{Message: bizRes.GetMessage()}, nil
}
