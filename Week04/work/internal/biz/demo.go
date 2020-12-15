package biz

import (
	"context"
)

type HelloRequest interface {
	Hello
	SetName(string)
	GetName() string
	BizHello(ctx context.Context, request HelloReq) HelloResponse
}

type Hello interface {
	ReSendMessage(context.Context, string) string
}

func NewHelloRequest(hello Hello) HelloRequest {
	return HelloRequest(&HelloReq{hello, ""})
}

type HelloReq struct {
	Hello
	name string
}

func (h *HelloReq) SetName(s string) {
	h.name = s
}

func (h *HelloReq) GetName() string {
	return h.name
}

func NewHelloReq(name string) HelloReq {
	return HelloReq{name: name}
}

type HelloResponse interface {
	GetMessage() string
}

type helloResponse struct {
	message string
}

func (h *helloResponse) GetMessage() string {
	return h.message
}

func (h *HelloReq) BizHello(ctx context.Context, request HelloReq) HelloResponse {
	message := request.ReSendMessage(ctx, request.GetName())
	res := helloResponse{message}
	return HelloResponse(&res)
}
