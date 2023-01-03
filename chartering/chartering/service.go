package chartering

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/90poe/chartering-domain-service/api"
	chassisgrpc "github.com/90poe/service-chassis/grpc/v3"
)

type Client struct {
	api api.CharteringDomainServiceClient
}

type Location struct {
	Code string
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Waypoint struct {
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"long"`
	PortName    string  `json:"port_name"`
	PortCountry string  `json:"port_country"`
}

// NewClient creates a new instance of charteringClient
func NewClient(ctx context.Context, host string, port int) (*Client, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := chassisgrpc.NewClientConn(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to chartering service: %w", err)
	}

	return &Client{
		api.NewCharteringDomainServiceClient(conn),
	}, nil
}

func (c *Client) BuildRoute(ctx context.Context, locations []Location) ([]Location, float64, error) {
	var totalDistance float64

	waypoints := transformLocationsToWaypoints(locations)

	res, err := c.api.Route(ctx, &api.RouteRequest{Points: waypoints})
	if err != nil {
		return nil, 0, fmt.Errorf("error calling route: %w", err)
	}

	locations2, totalDistance := transformRouteLegsToLocations(res.GetLegs())

	return MergeLegsLocations(locations2), totalDistance, nil
}

// transformRouteLegsToLocations transforms route legs into 2 dimensions locations slice and sum totalDistance
func transformRouteLegsToLocations(legs []*api.Leg) ([][]Location, float64) {
	var totalDistance float64
	locations2 := make([][]Location, 0, len(legs))
	for _, leg := range legs {
		legPoints := make([]Location, 0, len(leg.Points))
		for _, p := range leg.Points {
			legPoints = append(legPoints, Location{Lat: p.Latitude, Long: p.Longitude})
		}
		locations2 = append(locations2, legPoints)
		totalDistance += leg.GetMetrics().GetTotalDistance()
	}
	return locations2, totalDistance
}

// transformLocationsToWaypoints transforms points from model to api waypoints
func transformLocationsToWaypoints(locations []Location) []*api.Point {
	waypoints := make([]*api.Point, 0, len(locations))
	for i, l := range locations {
		waypoints = append(waypoints, &api.Point{
			Label:     strconv.Itoa(i),
			Latitude:  l.Lat,
			Longitude: l.Long,
		})
	}
	return waypoints
}

func (c *Client) SearchPort(ctx context.Context, name, countryISO string) (Location, error) {
	res, err := c.api.Port(ctx, &api.PortRequest{
		Pagination: &api.Pagination{
			Size: 1,
		},
		Filter: &api.PortFilter{
			Name:    name,
			Country: countryISO,
		},
	})

	if err != nil {
		return Location{}, err
	}

	if len(res.GetPorts()) != 0 {
		p := res.GetPorts()[0]
		return Location{Lat: p.GetLatitude(), Long: p.GetLongitude(), Code: p.Locode}, nil
	}

	return Location{}, errors.New("port not found")
}

func (c *Client) DiscoverLocations(ctx context.Context, waypoints []Waypoint) ([]Location, error) {
	locations := make([]Location, 0, len(waypoints))

	for _, v := range waypoints {
		isPort := v.Longitude == 0 && v.Latitude == 0 && v.PortName != ""
		if isPort {
			loc, err := c.SearchPort(ctx, v.PortName, v.PortCountry)
			if err != nil {
				return nil, fmt.Errorf("not found any port by '%s' on country '%s': %w", v.PortName, v.PortCountry, err)
			}
			locations = append(locations, loc)
			continue
		}
		locations = append(locations, Location{Lat: v.Latitude, Long: v.Longitude})
	}

	return locations, nil
}

func MergeLegsLocations(legsLocations [][]Location) []Location {
	var locations []Location

	// transform points
	for i, leg := range legsLocations {
		for j, pos := range leg {
			loc := Location{
				Lat:  pos.Lat,
				Long: pos.Long,
			}
			// if nth leg, first point equal to last point, discard
			if i > 0 && j == 0 && len(locations) > 0 && locations[len(locations)-1] == loc {
				continue
			}

			locations = append(locations, loc)
		}
	}

	return locations
}
