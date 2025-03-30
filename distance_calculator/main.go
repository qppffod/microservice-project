package main

import (
	"log"

	"github.com/qppffod/microservice-project/aggregator/client"
)

const (
	aggregatorEndpoint = "http://localhost:3000"
	grpcEndpoint       = ":3001"
)

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	httpclient := client.NewHTTPClient(aggregatorEndpoint)
	// grpcClient, err := client.NewGRPCClient(grpcEndpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	kafkaConsumer, err := NewKafkaConsumer("obudata", svc, httpclient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
