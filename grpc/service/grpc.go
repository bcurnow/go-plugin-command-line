package service

import (
	"context"

	"github.com/bcurnow/go-plugin-command-line/grpc/proto"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type Plugin struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl service.Service
}

func (p *Plugin) GRPCServer(broker *goplugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterServiceServer(s, &GRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *Plugin) GRPCClient(ctx context.Context, broker *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: proto.NewServiceClient(c),
		broker: broker,
	}, nil
}

type GRPCClient struct {
	broker *goplugin.GRPCBroker
	client proto.ServiceClient
}

func (c *GRPCClient) Name() string {
	resp, err := c.client.Name(context.Background(), &proto.Empty{})
	if err != nil {
		panic(err)
	}
	return resp.Name
}

func (c *GRPCClient) Log(val string) {
	_, err := c.client.Log(context.Background(), &proto.LogRequest{
		Val: val,
	})
	if err != nil {
		panic(err)
	}
}

type GRPCServer struct {
	Impl   service.Service
	broker *goplugin.GRPCBroker
}

func (s *GRPCServer) Name(ctx context.Context, req *proto.Empty) (*proto.NameResponse, error) {
	return &proto.NameResponse{
		Name: s.Impl.Name(),
	}, nil
}

func (s *GRPCServer) Log(ctx context.Context, req *proto.LogRequest) (*proto.Empty, error) {
	s.Impl.Log(req.Val)
	return &proto.Empty{}, nil
}

// Ensure we're properly implementing the interface
var _ goplugin.GRPCPlugin = &Plugin{}
