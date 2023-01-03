package portinfo

import (
	"context"
	"fmt"

	"github.com/90poe/port-information-domain-service/v3/api"
	chassisgrpc "github.com/90poe/service-chassis/grpc/v3"
)

type Client struct {
	conn api.PortInformationDomainServiceClient
}

func NewClient(ctx context.Context, host string, port int) (*Client, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := chassisgrpc.NewClientConn(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to vessel info service: %w", err)
	}

	return &Client{
		conn: api.NewPortInformationDomainServiceClient(conn),
	}, nil
}

func (c *Client) GetPort(ctx context.Context, name string) (*api.Port, error) {
	req := &api.SearchPortsRequest{Name: name, Pagination: &api.Pagination{Limit: 1}}
	ports, err := c.conn.SearchPorts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to search ports by name [%s]: %w", name, err)
	}

	if len(ports.GetPorts()) == 0 {
		return nil, fmt.Errorf("no ports found by name: %s", name)
	}

	return ports.GetPorts()[0], nil
}

func (c *Client) GetAllPorts(ctx context.Context) (api.Ports, error) {
	req := &api.GetPortsRequest{}
	resp, err := c.conn.GetPorts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all ports: %w", err)
	}

	return resp, nil
}

func (c *Client) GetPortName(ctx context.Context, code string) (string, error) {
	req := &api.ByCodeRequest{
		Code: &api.EntityCode{
			Port: code,
		},
	}
	port, err := c.conn.GetPortByCode(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to get port by code [%s]: %w", code, err)
	}

	if port == nil {
		return "", fmt.Errorf("no port found by code: %s", code)
	}

	return port.GetName(), nil
}
