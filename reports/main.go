package main

import (
	"context"
	"log"
	"time"

	"github.com/90poe/voyage-monitor/reports-service/api"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/timestamppb"

	reports "sandbox/reports/service"
	"sandbox/tokenprovider"
)

const vesselID = "3b47ee32-d28f-647c-3df5-5949098ac0cd"

func main() {
	ctx := context.Background()

	tokenProvider, err := tokenprovider.Build(ctx, tokenprovider.Int)
	if err != nil {
		log.Panicln(err)
	}

	ctx, err = tokenProvider.PrepareAuthContext(ctx, tokenprovider.Zodiac)
	if err != nil {
		log.Panicln(err)
	}

	cli, err := reports.NewClient(ctx, "localhost", 8092)
	resp, err := cli.GetReportsSummaryForPeriod(ctx, &api.GetReportsSummaryForPeriodRequest{
		VesselId: &wrappers.StringValue{Value: vesselID},
		TimeFrom: timestamppb.New(time.Now().UTC().Add(-100 * time.Hour)),
		TimeTo:   timestamppb.New(time.Now().UTC()),
		Source: &api.Source{
			TenantUUID: tokenprovider.Ahoy,
			VesselUUID: vesselID,
			Tag:        nil,
		},
	})
	if err != nil {
		log.Panicln(err)
	}

	println(resp)
}
