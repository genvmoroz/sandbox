package reports

import (
	"context"
	"fmt"

	"github.com/90poe/service-chassis/grpc/v3"
	"github.com/90poe/voyage-monitor/reports-service/api"
)

type Client struct {
	conn api.ReportsServiceClient
}

func NewClient(ctx context.Context, host string, port int) (*Client, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClientConn(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to reports service: %w", err)
	}

	return &Client{
		conn: api.NewReportsServiceClient(conn),
	}, nil
}

func (c *Client) GetReportsSummaryForPeriod(ctx context.Context, in *api.GetReportsSummaryForPeriodRequest) ([]*api.ReportSummaryInfo, error) {
	resp, err := c.conn.GetReportsSummaryForPeriod(ctx, in)
	if err != nil {
		return nil, err
	}

	return resp.GetReports(), nil
}
