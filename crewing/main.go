package main

import (
	"context"
	"fmt"
	"log"
	crewing "sandbox/crewing/service"
	"sandbox/tokenprovider"
)

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

	cli, err := crewing.NewClient(ctx, "localhost", 8086)
	ranks, err := cli.GetAllRanks(ctx)
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println(len(ranks))
}
