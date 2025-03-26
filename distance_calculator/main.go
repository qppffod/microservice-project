package main

import "log"

func main() {
	kafkaConsumer, err := NewKafkaConsumer("obudata")
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
