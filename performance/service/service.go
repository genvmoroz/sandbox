package performance

import (
	"context"
	"fmt"

	"github.com/90poe/performance/vessel-performance-information-service/api"
	"github.com/90poe/service-chassis/grpc/v3"
)

type Client struct {
	conn api.VesselPerformanceInformationClient
}

func NewClient(ctx context.Context, host string, port int) (*Client, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClientConn(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to vessel info service: %w", err)
	}

	return &Client{
		conn: api.NewVesselPerformanceInformationClient(conn),
	}, nil
}

func (c *Client) GetVesselPerformanceInformation(ctx context.Context, in *api.GetVesselPerformanceInformationRequest) (*api.VesselPerformanceData, error) {
	return c.conn.GetVesselPerformanceInformation(ctx, in)
}
