package main

import (
	"time"

	"github.com/qppffod/microservice-project/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"took": time.Since(start),
		}).Info("Aggregate distance")
	}(time.Now())

	err = l.next.AggregateDistance(distance)

	return err
}
