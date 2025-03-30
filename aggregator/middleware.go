package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/qppffod/microservice-project/types"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
	next           Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &MetricsMiddleware{
		next:           next,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
	}
}

// func (m *MetricsMiddleware) pushMetrics(start time.Time, err error) {
// }

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())

	err = m.next.AggregateDistance(distance)
	return err
}

func (m *MetricsMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()
		if err != nil {
			m.errCounterCalc.Inc()
		}
	}(time.Now())

	inv, err = m.next.CalculateInvoice(obuID)

	return inv, err
}

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

func (l *LogMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			TotalDistance = 0
			TotalAmount   = 0
		)
		if inv != nil {
			TotalDistance = int(inv.TotalDistance)
			TotalAmount = int(inv.TotalAmount)
		}
		logrus.WithFields(logrus.Fields{
			"took":        time.Since(start),
			"err":         err,
			"obuID":       obuID,
			"totalDist":   TotalDistance,
			"totalAmount": TotalAmount,
		}).Info("Calculate invoice")
	}(time.Now())

	inv, err = l.next.CalculateInvoice(obuID)

	return inv, err
}
