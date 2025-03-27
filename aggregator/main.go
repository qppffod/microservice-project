package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/qppffod/microservice-project/types"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address of the aggreagtor HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	var (
		svc Aggregator
	)
	svc = NewInvoiceAggregator(store)

	makeHTTPTransport(*listenAddr, svc)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
