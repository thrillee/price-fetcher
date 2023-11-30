package main

import (
	// "context"
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/thrillee/price-fetcher/client"
	"github.com/thrillee/price-fetcher/proto"
	// "fmt"
	// "log"
	// "github.com/thrillee/price-fetcher/client"
)

func main() {
	// cl := client.New("http://localhost:3000")
	// price, err := cl.FetchPrice(context.Background(), "ETT")
	// if err != nil {
	// log.Fatal(err)
	// }
	// fmt.Printf("%v\n", price)
	// return

	var (
		jsonAPiListenAddr = flag.String("json", ":3000", "Listen address of the  json service is running on")
		grpcAddr          = flag.String("grpc", ":4000", "Listen address of the  grpc service is running on")
		ctx               = context.Background()
	)

	flag.Parse()

	svc := NewLoggingService(NewMetricService(&priceFetcher{}))

	grpcClient, err := client.NewGRPCClient(":4000")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {

			time.Sleep(3 * time.Second)
			resp, err := grpcClient.FetchPrice(ctx, &proto.PriceRequest{Ticker: "BTC"})
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Fee -> %v\n", resp)
		}
	}()

	go GRPCServerlistenAndServe(*grpcAddr, svc)

	server := NewJsonAPIServer(*jsonAPiListenAddr, svc)
	server.Run()
}
