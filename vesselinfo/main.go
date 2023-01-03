package main

import (
	"context"
	"log"

	"sandbox/tokenprovider"
	vesselinfo "sandbox/vesselinfo/service"
)

const vesselID = "ba962077-13af-43b5-8d8f-44b333d30d23"

func main() {
	ctx := context.Background()

	tokenProvider, err := tokenprovider.Build(ctx, tokenprovider.Dev)
	if err != nil {
		log.Panicln(err)
	}

	ctx, err = tokenProvider.PrepareAuthContext(ctx, tokenprovider.Ahoy)
	if err != nil {
		log.Panicln(err)
	}

	cli, err := vesselinfo.NewClient(ctx, "localhost", 8080)
	resp, err := cli.GetShopTrialDataByIMO(ctx, vesselID)
	if err != nil {
		log.Panicln(err)
	}

	println(resp)
}
