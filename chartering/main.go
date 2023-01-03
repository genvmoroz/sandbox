package main

import (
	"context"
	"fmt"
	"sandbox/auxl"
	"sandbox/chartering/chartering"
	"sandbox/tokenprovider"
	"time"
)

func main() {

	dist := auxl.CalculateDistanceInKM([]auxl.DomainPosition{
		{Lat: 58.9824238, Long: -3.3411979},
		{Lat: 51.9111400, Long: 4.3502600},
	})

	to := time.Now().UTC().Add(auxl.HoursToDuration(dist / 60))

	from := to.Add(auxl.HoursToDuration(dist / 60))

	fmt.Println("to", to.Format(time.RFC3339))
	fmt.Println("from", from.Format(time.RFC3339))

	tokenProvider, err := tokenprovider.Build(context.Background(), tokenprovider.Dev)
	if err != nil {
		panic(err)
	}

	ctx, err := tokenProvider.PrepareAuthContext(context.Background(), tokenprovider.Ahoy)
	if err != nil {
		panic(err)
	}

	c, err := chartering.NewClient(ctx, "localhost", 8087)

	move := []auxl.MoveToPort{
		{Port: "SOUTHAMPTON", CountryISO2: "GB", StartSpeedKM: 0, EndSpeedKM: 40},
		{Port: "ROTTERDAM", CountryISO2: "NL", StartSpeedKM: 0, EndSpeedKM: 40},
		{Port: "SOUTHAMPTON", CountryISO2: "GB", StartSpeedKM: 0, EndSpeedKM: 40},
	}

	_ = auxl.BuildRoute(ctx, c, time.Now().UTC(), move...)

	//for i, j := 0, len(actions)-1; i < j; i, j = i+1, j-1 {
	//	actions[i], actions[j] = actions[j], actions[i]
	//}

	//var startWait []int
	//var endWait []int
	//
	//for index, a := range actions {
	//	if index == 0 {
	//		continue
	//	}
	//	if a.Wait && !actions[index-1].Wait {
	//		startWait = append(startWait, index)
	//	}
	//	if !a.Wait && actions[index-1].Wait {
	//		endWait = append(endWait, index)
	//	}
	//}

	//for index, i := range startWait {
	//	if index == 0 {
	//		actions[0].Actions = append(actions[0].Actions, fmt.Sprintf("{ \"action-name\": \"vessel-ais-type-5\", \"eta\": \"%v\", \"destination\": \"%s\"}", actions[i].Time.Format(time.RFC3339), move[index+1].Port))
	//	} else {
	//		actions[endWait[index-1]].Actions = append(actions[endWait[index-1]].Actions, fmt.Sprintf("{ \"action-name\": \"vessel-ais-type-5\", \"eta\": \"%v\", \"destination\": \"%v\"}", actions[i].Time.Format(time.RFC3339), move[index+1].Port))
	//	}
	//}
	//actions[endWait[len(endWait)-1]].Actions = append(actions[endWait[len(endWait)-1]].Actions, fmt.Sprintf("{ \"action-name\": \"vessel-ais-type-5\", \"eta\": \"%v\", \"destination\": \"%v\"}", actions[len(actions)-1].Time.Format(time.RFC3339), move[len(endWait)+1].Port))

	//res := ""
	//for i, a := range actions {
	//	if i == len(actions)-1 {
	//		res += a.String() + "\n"
	//	} else {
	//		res += a.String() + ",\n"
	//	}
	//}
	//
	//if err = os.WriteFile("actions.json", []byte(res), 777); err != nil {
	//	panic(err)
	//}
}
