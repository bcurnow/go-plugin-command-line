package service

import (
	"context"

	"github.com/bcurnow/go-plugin-command-line/grpc/proto"
	grpcservice "github.com/bcurnow/go-plugin-command-line/grpc/service"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type Plugin struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl command.Command
}

func (p *Plugin) GRPCServer(broker *goplugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterCommandServer(s, &GRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *Plugin) GRPCClient(ctx context.Context, broker *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: proto.NewCommandClient(c),
		broker: broker,
	}, nil
}

// make sure we correct implement plugin.GRPCPlugin
var _ goplugin.GRPCPlugin = &Plugin{}

type GRPCClient struct {
	broker *goplugin.GRPCBroker
	client proto.CommandClient
}

func (c *GRPCClient) Help() string {
	resp, err := c.client.Help(context.Background(), &proto.Empty{})
	if err != nil {
		panic(err)
	}
	return resp.Help
}

func (c *GRPCClient) Execute(args []string) error {
	_, err := c.client.Execute(context.Background(), &proto.ExecuteRequest{
		Args: args,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *GRPCClient) SetServices(reattaches map[string]plugin.Reattach) error {
	_, err := c.client.SetServices(context.Background(), &proto.SetServicesRequest{
		Reattaches: grpcservice.ToProtoReattaches(reattaches),
	})
	if err != nil {
		return err
	}
	return nil
}

type GRPCServer struct {
	Impl   command.Command
	broker *goplugin.GRPCBroker
}

func (s *GRPCServer) Help(ctx context.Context, req *proto.Empty) (*proto.HelpResponse, error) {
	return &proto.HelpResponse{
		Help: s.Impl.Help(),
	}, nil
}

func (s *GRPCServer) Execute(ctx context.Context, req *proto.ExecuteRequest) (*proto.Empty, error) {
	s.Impl.Execute(req.Args)
	return &proto.Empty{}, nil
}

func (s *GRPCServer) SetServices(ctx context.Context, req *proto.SetServicesRequest) (*proto.Empty, error) {
	return &proto.Empty{}, s.Impl.SetServices(grpcservice.ToReattaches(req.Reattaches))
}

// Ensure we're properly implementing the interface
var _ goplugin.GRPCPlugin = &Plugin{}
