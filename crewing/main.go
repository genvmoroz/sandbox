package main

import (
	"context"
	"fmt"
	"log"
	"strings"

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

	for i := 0; i < len(ranks); i++ {
		for y := i + 1; y < len(ranks); y++ {
			if strings.EqualFold(ranks[i].Department, ranks[y].Department) {
				fmt.Println(ranks[i].Department)
				fmt.Println(ranks[y].Department)
			}
		}
	}

}
