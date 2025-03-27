package main

import "log"

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer("obudata", svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
