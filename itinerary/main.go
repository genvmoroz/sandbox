package main

import (
	"context"

	itinerary "sandbox/itinerary/service"
	"sandbox/tokenprovider"
)

func main() {
	tokenProvider, err := tokenprovider.Build(context.Background(), tokenprovider.Dev)
	if err != nil {
		panic(err)
	}

	ctx, err := tokenProvider.PrepareAuthContext(context.Background(), tokenprovider.Ahoy)
	if err != nil {
		panic(err)
	}

	cli, err := itinerary.NewClient(ctx, "localhost", 8084)
	if err != nil {
		panic(err)
	}

	stay, err := cli.StayAtPortByIMO(ctx, 5992963)
	if err != nil {
		panic(err)
	}

	println(stay)
}
