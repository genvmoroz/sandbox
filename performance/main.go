package main

import (
	"context"
	"github.com/90poe/performance/vessel-performance-information-service/api"
	vesselinfoAPI "github.com/90poe/vessel-information-domain-service/v5/api"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v2"
	"log"
	performance "sandbox/performance/service"
	"sandbox/tokenprovider"
	vesselinfo "sandbox/vesselinfo/service"
	"time"
)

const vesselID = "cb0a8d59-753a-4990-b2ba-456b9bea364b"

type VesselType string

const (
	BULKCARRIER VesselType = "BULK CARRIER"
	CHEMICAL    VesselType = "CHEMICAL"
	LPGCARRIER  VesselType = "LPG CARRIER"
	OILCRUDE    VesselType = "OIL CRUDE"
	OILPRODUCT  VesselType = "OIL PRODUCT"
	CONTAINER   VesselType = "CONTAINER"
	PCTC        VesselType = "PCTC"
)

var respMap = map[VesselType]string{
	BULKCARRIER: "3b47ee32-d28f-647c-3df5-5949098ac0cd",
	CHEMICAL:    "970cf5e9-3267-7e1c-fc71-ecf340a933d6",
	LPGCARRIER:  "396f4891-b8d5-6ed6-c422-443263198aac",
	OILCRUDE:    "b94a41c4-5bb1-a7ce-5fe3-d6cbe84ae04b",
	OILPRODUCT:  "ab90577f-7803-840f-094c-58579915c671",
	CONTAINER:   "531931ca-31c5-231c-409b-5498a77041f7",
	PCTC:        "1436874f-4171-01b8-ae4a-7d4019f723df",
}

type StaticData struct {
	PerformanceData *api.VesselPerformanceData
	ShopTrialData   *vesselinfoAPI.ShopTrialData
}

var StaticDatas map[VesselType]StaticData

func main() {
	ctx := context.Background()

	tokenProvider, err := tokenprovider.Build(ctx, tokenprovider.Int)
	if err != nil {
		log.Panicln(err)
	}

	ctx, err = tokenProvider.PrepareAuthContext(ctx, tokenprovider.Ahoy)
	if err != nil {
		log.Panicln(err)
	}

	cli, err := performance.NewClient(ctx, "localhost", 8092)
	if err != nil {
		log.Panicln(err)
	}
	viCli, err := vesselinfo.NewClient(ctx, "localhost", 8082)
	if err != nil {
		log.Panicln(err)
	}

	r, _ := cli.GetVesselPerformanceInformation(ctx, &api.GetVesselPerformanceInformationRequest{
		VesselId: vesselID,
		Date:     timestamppb.New(time.Now().UTC()),
	})

	println(r)

	resp, err := viCli.GetShopTrialDataByIMO(ctx, vesselID)
	if err != nil {
		log.Panicln(err)
	}

	println(resp)

	t, _ := viCli.GetVesselTypes(ctx)
	for _, tt := range t.GetVesselTypes() {
		println(tt.GetName())
	}

	StaticDatas = make(map[VesselType]StaticData)
	for key, val := range respMap {
		resp, err := cli.GetVesselPerformanceInformation(ctx, &api.GetVesselPerformanceInformationRequest{
			VesselId: val,
		})
		if err != nil {
			log.Panicln(err)
		}
		if resp == nil {
			log.Panicln("resp is nil")
		}

		shopTri, err := viCli.GetShopTrialDataByIMO(ctx, val)
		if err != nil {
			log.Panicln(err)
		}

		StaticDatas[key] = StaticData{
			PerformanceData: resp,
			ShopTrialData:   shopTri,
		}
	}

	content, err := yaml.Marshal(StaticDatas)
	if err != nil {
		log.Panicln(err)
	}

	println(content)
}
