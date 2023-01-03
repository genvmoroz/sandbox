package itinerary

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/90poe/service-chassis/grpc/v3"
	"github.com/90poe/vessel-itinerary-domain-service/v3/api"
)

type Client struct {
	conn api.VesselItineraryDomainServiceClient
}

func NewClient(ctx context.Context, host string, port int) (*Client, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClientConn(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to vessel info service: %w", err)
	}

	return &Client{
		conn: api.NewVesselItineraryDomainServiceClient(conn),
	}, nil
}

func (c *Client) CreatePortCall(ctx context.Context, payload *api.PortCallPayload) (*api.PortCall, error) {
	req := &api.CreatePortCallRequest{Payload: payload}
	resp, err := c.conn.CreatePortCall(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create PortCall: %w", err)
	}
	if resp.GetPortCall() == nil {
		return nil, errors.New("create PortCall response returned nullable PortCall")
	}

	return resp.GetPortCall(), nil
}

func (c *Client) UpdatePortCall(ctx context.Context, call *api.PortCall) error {
	if call == nil {
		return errors.New("PortCall is missing")
	}

	req := portCallToUpdatePortCallRequest(call)
	resp, err := c.conn.UpdatePortCall(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update PortCall: %w", err)
	}

	if resp.GetPortCall() == nil {
		return errors.New("update PortCall response returned nullable PortCall")
	}

	return nil
}

func (c *Client) GetFutureDestinationItinerary(ctx context.Context, imo int32) (*api.PortCall, error) {
	req := &api.FutureItineraryRequest{
		Imo: imo,
	}
	resp, err := c.conn.GetFutureItineraryForVessel(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get future itinerary for vessel with imo [%d]: %w", imo, err)
	}

	lastIndex := len(resp.GetPortCalls()) - 1

	if lastIndex < 0 {
		return nil, nil
	}

	return resp.GetPortCalls()[lastIndex], nil
}

func (c *Client) GetLatestDestination(ctx context.Context, imo int32) (*api.PortCall, error) {
	req := &api.ItineraryRequest{
		Imo:   imo,
		Limit: 1,
	}
	resp, err := c.conn.GetItineraryForVessel(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get itinerary for the vessel with imo [%d]: %w", imo, err)
	}

	if len(resp.GetPortCalls()) == 0 {
		return nil, nil
	}

	return resp.GetPortCalls()[0], nil
}

func (c *Client) GetCurrentVoyageByIMO(ctx context.Context, imo int32) (*api.Voyage, error) {
	req := &api.GetCurrentVoyagesRequest{
		Pagination: &api.FullPagination{
			Limit:  -1,
			Offset: 0,
		},
	}
	resp, err := c.conn.GetCurrentVoyages(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get itinerary for the vessel with imo [%d]: %w", imo, err)
	}

	if len(resp.GetValues()) == 0 {
		return nil, nil
	}

	for _, voyage := range resp.GetValues() {
		if voyage.GetImo() == imo {
			return voyage, nil
		}
	}

	return nil, nil
}

func (c *Client) StayAtPortByIMO(ctx context.Context, imo int32) (bool, error) {
	voyage, err := c.GetCurrentVoyageByIMO(ctx, imo)
	if err != nil {
		return false, fmt.Errorf("failed to get current voyage for vessel by imo: %v: %w", imo, err)
	}

	if voyage != nil {
		return c.currentTimeBetweenArrivalAndDeparture(voyage)
	}

	return false, nil
}

func (c *Client) currentTimeBetweenArrivalAndDeparture(voyage *api.Voyage) (bool, error) {
	arrival, err := time.Parse(time.RFC3339, voyage.GetArrival())
	if err != nil {
		return false, fmt.Errorf("failed to parse arrival time: %w", err)
	}

	departure, err := time.Parse(time.RFC3339, voyage.GetDeparture())
	if err != nil {
		return false, fmt.Errorf("failed to parse arrival time: %w", err)
	}

	return time.Now().UTC().After(arrival) &&
			time.Now().UTC().Before(departure),
		nil
}

func portCallToUpdatePortCallRequest(call *api.PortCall) *api.UpdatePortCallRequest {
	if call == nil {
		return nil
	}

	return &api.UpdatePortCallRequest{
		Id:      call.GetId(),
		Version: call.GetVersion(),
		Payload: &api.PortCallPayload{
			Imo:              call.GetImo(),
			Arrival:          call.GetArrival(),
			Departure:        call.GetDeparture(),
			Berthing:         call.GetBerthing(),
			DestinationId:    call.GetDestination().GetId(),
			DestinationCode:  call.GetDestination().GetCode(),
			Activities:       call.GetActivities(),
			PlannedStatuses:  call.GetPlannedStatuses(),
			AgentAssignments: call.GetAgentAssignments(),
			BerthNotes:       call.GetBerthNotes(),
			VesselId:         call.GetVesselId(),
		},
	}
}
