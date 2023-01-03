package main

import (
	"context"
	"fmt"
	"log"

	portinfo "sandbox/portinfo/service"
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

	cli, err := portinfo.NewClient(ctx, "localhost", 8080)
	resp, err := cli.GetAllPorts(ctx)
	if err != nil {
		log.Panicln(err)
	}

	for _, port := range resp.Ports {
		fmt.Println(port.GetName())
	}
}
