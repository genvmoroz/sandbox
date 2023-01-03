package auxl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math"
	"time"

	geo "github.com/kellydunn/golang-geo"

	"sandbox/chartering/chartering"
)

type MoveToPort struct {
	Port         string
	CountryISO2  string
	StartSpeedKM float64
	EndSpeedKM   float64
}

type Action struct {
	ID      string
	Pos     Position
	Actions Actions
	Time    time.Time
	Speed   float64
	Wait    bool
	Port    string
}

type Actions []string

func (aa Actions) String() string {
	res := ""
	for i, a := range aa {
		if i == 0 {
			res += a
			if len(aa) != 1 {
				res += ",\n"
			}
		} else if i == len(aa)-1 {
			res += "\t\t\t" + a
		} else {
			res += "\t\t\t" + a + "," + "\n"
		}
	}

	return res
}

func (a Action) String() string {
	return fmt.Sprintf("\t{\n\t\t\"id\": \"%s\",\n\t\t\"pos\": { \"lat\": %f, \"lon\": %f },\n\t\t\"actions\": [\n\t\t\t%v\n\t\t]\n\t}", a.ID, a.Pos.Lat, a.Pos.Lon, a.Actions)
}

var count = 0

type Position struct {
	Lat float64
	Lon float64
}

type VesselPosition struct {
	IMO       uint
	NavStatus string
	When      time.Time
	Position  Position
}

type CommandEnvelope struct {
	VesselID    string
	CommandName string
	Command     json.RawMessage
}

func MustParseJSON(o interface{}) json.RawMessage {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}

	return b
}

func BuildRoute(ctx context.Context, client *chartering.Client, prevT time.Time, move ...MoveToPort) []Action {
	var actions []Action

	ports := make([]chartering.Location, len(move))
	for i, m := range move {
		p, err := client.SearchPort(ctx, m.Port, m.CountryISO2)
		if err != nil {
			panic(err)
		}
		ports[i] = p
	}

	points := make([]routePoint, 0)

	for i := 0; i < len(ports)-1; i++ {
		var (
			currentP = ports[i]
			nextP    = ports[i+1]
		)
		locations, dist, err := client.BuildRoute(ctx, []chartering.Location{{Lat: currentP.Lat, Long: currentP.Long}, {Lat: nextP.Lat, Long: nextP.Long}})
		if err != nil {
			panic(err)
		}

		if len(locations) == 1 {
			panic("invalid route")
		}

		a := Action{
			ID:  uuid.NewString(),
			Pos: Position{Lat: locations[0].Lat, Lon: locations[0].Long},
			Actions: []string{
				fmt.Sprintf("{ \"action-name\": \"set-speed\", \"speed\": %f }", 60*0.5399568035),
				"{ \"action-name\": \"nav-status\", \"state\": \"UnderWayUsingEngine\" }",
				fmt.Sprintf("{ \"action-name\": \"vessel-ais-type-5\", \"eta\": \"%v\", \"destination\": \"%v\"}", prevT.Add(HoursToDuration((dist*1.852)/60)).Format(time.RFC3339), nextP.Code),
				"{ \"action-name\": \"vessel-position\" }",
			},
		}
		actions = append(actions, a)

		for _, ll := range locations {
			points = append(points, routePoint{
				Position: Position{Lat: ll.Lat, Lon: ll.Long},
				Speed:    60,
			})
		}

		if i != len(ports)-1 {
			pp, t := wait(prevT, Position{Lat: points[len(points)-1].Position.Lat, Lon: points[len(points)-1].Position.Lon}, time.Hour, currentP.Code)
			points = append(points, pp...)
			prevT = t
		}

		//prevSS := move[i].StartSpeedKM
		//prevES := move[i].EndSpeedKM
		//
		//for index := 0; index < len(locations)-1; index++ {
		//	currentL := locations[index]
		//	nextL := locations[index+1]
		//
		//	if index == 0 {
		//		pp, t, prevS := route(Position{Lat: currentL.Lat, Lon: currentL.Long}, Position{Lat: nextL.Lat, Lon: nextL.Long}, prevSS, prevES, true, prevT)
		//		points = append(points, pp...)
		//
		//		prevT = t
		//		prevSS = prevS
		//	} else if index == len(locations)-2 {
		//		pp, t, prevS := route(Position{Lat: currentL.Lat, Lon: currentL.Long}, Position{Lat: nextL.Lat, Lon: nextL.Long}, prevSS, 0, false, prevT)
		//		points = append(points, pp...)
		//
		//		prevT = t
		//		prevSS = prevS
		//	} else {
		//		pp, t, prevS := route(Position{Lat: currentL.Lat, Lon: currentL.Long}, Position{Lat: nextL.Lat, Lon: nextL.Long}, prevSS, prevES, true, prevT)
		//		points = append(points, pp...)
		//
		//		prevT = t
		//		prevSS = prevS
		//	}
		//}

		for _, p := range points {
			a := Action{
				ID:    uuid.NewString(),
				Pos:   p.Position,
				Time:  p.Time,
				Speed: p.Speed,
				Wait:  p.Wait,
				Port:  p.Port,
			}
			if p.Speed == 0 {
				a.Actions = append(a.Actions, "{ \"action-name\": \"nav-status\", \"state\": \"AtAnchor\" }")
				a.Actions = append(a.Actions, "{ \"action-name\": \"vessel-position\" }")
			} else {
				a.Actions = append(a.Actions, fmt.Sprintf("{ \"action-name\": \"set-speed\", \"speed\": %f }", p.Speed*0.5399568035))
				a.Actions = append(a.Actions, "{ \"action-name\": \"nav-status\", \"state\": \"UnderWayUsingEngine\" }")
				a.Actions = append(a.Actions, "{ \"action-name\": \"vessel-position\" }")
			}

			actions = append(actions, a)
		}
	}

	return actions
}

const SpeedAccelerationPeerKM = 1

type routePoint struct {
	Position Position
	Time     time.Time
	Dist     float64
	Speed    float64
	Wait     bool
	Port     string
}

func route(start, end Position, startSpeedKM, endSpeedKM float64, accelerationDirection bool, startTime time.Time) ([]routePoint, time.Time, float64) {
	var (
		startP = geo.NewPoint(start.Lat, start.Lon)
		endP   = geo.NewPoint(end.Lat, end.Lon)

		dist    = startP.GreatCircleDistance(endP)
		durStep = 10 * time.Minute
		prevP   = startP
		prevS   = startSpeedKM

		points []routePoint
	)

	points = append(points,
		routePoint{
			Position: Position{Lat: start.Lat, Lon: start.Lon},
			Time:     startTime,
			Dist:     0,
			Speed:    startSpeedKM,
		},
	)

	if accelerationDirection {
		for count := math.Abs(endSpeedKM-startSpeedKM) / SpeedAccelerationPeerKM; dist > 0 && count > 0; count-- {
			var (
				bearing   = prevP.BearingTo(endP)
				nextS     = prevS + SpeedAccelerationPeerKM
				nextD     = distKM(dist, SpeedAccelerationPeerKM)
				nextT     = duration(nextD, nextS)
				timestamp = startTime.Add(-nextT)
				nextP     = prevP.PointAtDistanceAndBearing(nextD, bearing)
			)

			points = append(points,
				routePoint{
					Position: Position{Lat: nextP.Lat(), Lon: nextP.Lng()},
					Time:     timestamp,
					Dist:     nextD,
					Speed:    nextS,
				},
			)

			prevP = nextP
			prevS = nextS
			dist = dist - nextD
			startTime = timestamp
		}

		for dist > 0 {
			var (
				bearing   = prevP.BearingTo(endP)
				nextD     = distKM(dist, distance(prevS, durStep))
				nextT     = duration(nextD, prevS)
				timestamp = startTime.Add(-nextT)
				nextP     = prevP.PointAtDistanceAndBearing(nextD, bearing)
			)

			points = append(points,
				routePoint{
					Position: Position{Lat: nextP.Lat(), Lon: nextP.Lng()},
					Time:     timestamp,
					Dist:     nextD,
					Speed:    prevS,
				},
			)

			prevP = nextP
			dist = dist - nextD
			startTime = timestamp
		}
	} else {
		previosDist := dist
		if dist < math.Abs(startSpeedKM-endSpeedKM-1)*SpeedAccelerationPeerKM {
			for dist > 0 {
				var (
					acceleration = math.Abs(startSpeedKM-endSpeedKM-1) / previosDist
					bearing      = prevP.BearingTo(endP)
					nextS        = speedKMPerHR(prevS, acceleration)
					nextD        = distKM(dist, SpeedAccelerationPeerKM)
					nextT        = duration(nextD, nextS)
					timestamp    = startTime.Add(-nextT)
					nextP        = prevP.PointAtDistanceAndBearing(nextD, bearing)
				)

				points = append(points,
					routePoint{
						Position: Position{Lat: nextP.Lat(), Lon: nextP.Lng()},
						Time:     timestamp,
						Dist:     nextD,
						Speed:    nextS,
					},
				)

				prevP = nextP
				prevS = nextS
				dist = dist - nextD
				startTime = timestamp
			}
		}
		for distWithoutAcceleration := dist - (math.Abs(startSpeedKM-endSpeedKM-1) * SpeedAccelerationPeerKM); dist > 0 && distWithoutAcceleration > 0; {
			var (
				bearing   = prevP.BearingTo(endP)
				nextD     = distKM(distWithoutAcceleration, distance(prevS, durStep))
				nextT     = duration(nextD, prevS)
				timestamp = startTime.Add(-nextT)
				nextP     = prevP.PointAtDistanceAndBearing(nextD, bearing)
			)

			points = append(points,
				routePoint{
					Position: Position{Lat: nextP.Lat(), Lon: nextP.Lng()},
					Time:     timestamp,
					Dist:     nextD,
					Speed:    prevS,
				},
			)

			prevP = nextP
			dist = dist - nextD
			distWithoutAcceleration = distWithoutAcceleration - nextD
			startTime = timestamp
		}

		for dist > 0 {
			var (
				bearing   = prevP.BearingTo(endP)
				nextS     = prevS - SpeedAccelerationPeerKM
				nextD     = distKM(dist, SpeedAccelerationPeerKM)
				nextT     = duration(nextD, nextS)
				timestamp = startTime.Add(-nextT)
				nextP     = prevP.PointAtDistanceAndBearing(nextD, bearing)
			)

			points = append(points,
				routePoint{
					Position: Position{Lat: nextP.Lat(), Lon: nextP.Lng()},
					Time:     timestamp,
					Dist:     nextD,
					Speed:    nextS,
				},
			)

			prevP = nextP
			prevS = nextS
			dist = dist - nextD
			startTime = timestamp
		}
	}

	return points, startTime, prevS
}

func wait(now time.Time, pos Position, wait time.Duration, port string) ([]routePoint, time.Time) {
	stepDur := 4 * time.Hour

	prevT := now

	var points []routePoint
	for wait > 0 {

		nextD := stepDur
		if wait < stepDur {
			nextD = wait
		}
		points = append(points, routePoint{Position: pos, Time: prevT, Wait: true, Port: port})

		wait = wait - nextD
		prevT = prevT.Add(-nextD)
	}

	return points, prevT
}

func duration(distKM, speedKMPerHR float64) time.Duration {
	return HoursToDuration(distKM / speedKMPerHR)
}

func distKM(distKM, diffKM float64) float64 {
	res := distKM - diffKM
	if res < 0 {
		return distKM
	}
	return diffKM
}

func speedKMPerHR(speedKM, diffKM float64) float64 {
	res := speedKM - diffKM
	if res < 0 {
		return speedKM
	}
	return res
}

func distance(speedKMPerHR float64, dur time.Duration) float64 {
	return speedKMPerHR * dur.Hours()
}

func HoursToDuration(hours float64) time.Duration {
	return time.Duration(hours * float64(time.Hour))
}

func Command(imo uint, p *geo.Point, time time.Time) CommandEnvelope {
	return CommandEnvelope{
		VesselID:    "vesselID",
		CommandName: "vessel-position",
		Command: MustParseJSON(
			VesselPosition{
				IMO: imo,
				Position: Position{
					Lat: p.Lat(),
					Lon: p.Lng(),
				},
				NavStatus: "UnderWayUsingEngine",
				When:      time,
			}),
	}
}

func createVessel(imo uint, time time.Time) CommandEnvelope {
	return CommandEnvelope{
		VesselID:    "vesselID",
		CommandName: "create-vessel-from-template",
		Command: MustParseJSON(
			CreateVessel{
				IMO:                       imo,
				MMSI:                      imo * 1000,
				Name:                      fmt.Sprintf("DELETE_ME_%d", imo),
				DeliveredToPrincipalsDate: time,
			}),
	}
}

type CreateVessel struct {
	IMO  uint
	MMSI uint
	Name string
	Type string

	DeliveredToPrincipalsDate time.Time
}

func Count(t, step time.Duration) int {
	c := int(t / step)
	if c < 1 {
		return 1
	}

	return c
}

func getCount() int {
	count++

	return count
}

type DomainPosition struct {
	Lat       float64
	Long      float64
	Timestamp int64 // in unix millis
	SOG       float64
	COG       float64
}

func CalculateDistanceInKM(positions []DomainPosition) float64 {
	if len(positions) < 2 {
		return 0
	}

	var (
		first  = geo.NewPoint(positions[0].Lat, positions[0].Long)
		second = geo.NewPoint(positions[1].Lat, positions[1].Long)
		dist   = first.GreatCircleDistance(second)
	)

	return dist + CalculateDistanceInKM(positions[1:])
}
