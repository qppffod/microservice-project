package main

import (
	"log"

	"github.com/qppffod/microservice-project/aggregator/client"
)

const aggregatorEndpoint = "http://localhost:3000/aggregate"

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer("obudata", svc, client.NewClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
