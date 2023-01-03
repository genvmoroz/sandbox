package vesselinfo

import (
	"context"
	"fmt"

	"github.com/90poe/service-chassis/grpc/v3"
	"github.com/90poe/vessel-information-domain-service/v5/api"
)

type Client struct {
	conn api.VesselInformationServiceClient
}

func NewClient(ctx context.Context, host string, port int) (*Client, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClientConn(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to vessel info service: %w", err)
	}

	return &Client{
		conn: api.NewVesselInformationServiceClient(conn),
	}, nil
}

func (c *Client) GetShopTrialDataByIMO(ctx context.Context, id string) (*api.ShopTrialData, error) {
	return c.conn.GetShopTrialDataByIMO(ctx, &api.GetShopTrialDataByIMORequest{ID: id})
}

func (c *Client) GetVesselTypes(ctx context.Context) (*api.VesselTypes, error) {
	return c.conn.GetVesselTypes(ctx, &api.Empty{})
}
