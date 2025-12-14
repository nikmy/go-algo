package netx

import (
	"context"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Node interface {
	Addr() string
	Call(ctx context.Context, peer string, method string, arg proto.Message) (proto.Message, error)

	Init()
	Fail()
	Work(ctx context.Context, dialer Transport)
}

type Transport interface {
	Dial(ctx context.Context, peer string, method string, arg proto.Message, result proto.Message) error
	Broadcast(ctx context.Context, peers []Node, method string, message proto.Message) ([]Node, error)
}

type Ok = emptypb.Empty

type Topology struct {
	Adjacent map[string][]string
	NetNodes []Node
}
