package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sandbox/auxl"
	"time"

	geo "github.com/kellydunn/golang-geo"
)

const vesselIMO uint = 1234568

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

func main() {

	startP := geo.NewPoint(-6.4899833, -104.0776169)
	endP := geo.NewPoint(-6.3152985, -143.6190282)

	dist := startP.GreatCircleDistance(endP)
	hours := hoursToDuration(dist / 60)

	start := time.Now().UTC().Add(-1 * hours)

	step := 10 * time.Minute

	commands := append([]CommandEnvelope{}, createVessel(start))
	commands = append(commands, command(startP, start))

	c := count(hours, step)
	d := dist / float64(c)
	prevP := startP
	for i := 1; i <= c; i++ {
		bearing := prevP.BearingTo(endP)

		nextT := step * time.Duration(i)
		nextP := prevP.PointAtDistanceAndBearing(d, bearing)

		commands = append(commands, command(nextP, start.Add(nextT)))

		prevP = nextP
	}

	if err := os.WriteFile("out.json", auxl.MustParseJSON(commands), 777); err != nil {
		panic(err)
	}
}

func command(p *geo.Point, time time.Time) CommandEnvelope {
	return CommandEnvelope{
		VesselID:    "vesselID",
		CommandName: "vessel-position",
		Command: auxl.MustParseJSON(
			VesselPosition{
				IMO: vesselIMO,
				Position: Position{
					Lat: p.Lat(),
					Lon: p.Lng(),
				},
				NavStatus: "UnderWayUsingEngine",
				When:      time,
			}),
	}
}

func createVessel(time time.Time) CommandEnvelope {
	return CommandEnvelope{
		VesselID:    "vesselID",
		CommandName: "create-vessel-from-template",
		Command: auxl.MustParseJSON(
			CreateVessel{
				IMO:                       vesselIMO,
				MMSI:                      vesselIMO * 1000,
				Name:                      fmt.Sprintf("DELETE_ME_%d", vesselIMO),
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

func count(t, step time.Duration) int {
	c := int(t / step)
	if c < 1 {
		return 1
	}

	return c
}

func hoursToDuration(hours float64) time.Duration {
	return time.Duration(hours * float64(time.Hour))
}
