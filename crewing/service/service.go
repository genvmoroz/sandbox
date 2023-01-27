package portinfo

import (
	"context"
	"fmt"

	"github.com/90poe/atlantic/crewing-domain-service/api"
	chassisgrpc "github.com/90poe/service-chassis/grpc/v3"
)

type Client struct {
	conn api.CrewingDomainServiceClient
}

func NewClient(ctx context.Context, host string, port int) (*Client, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := chassisgrpc.NewClientConn(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to crewing service: %w", err)
	}

	return &Client{
		conn: api.NewCrewingDomainServiceClient(conn),
	}, nil
}

func (r *Client) GetAllRanks(ctx context.Context) ([]*api.Rank, error) {
	resp, err := r.conn.GetAllRanks(ctx, &api.Empty{})
	if err != nil {
		return nil, fmt.Errorf("grpc [GetRanks]: %w", err)
	}

	return resp.GetRanks(), nil
}
