package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qppffod/microservice-project/types"
	"google.golang.org/grpc"
)

func main() {
	HTTPListenAddr := flag.String("HTTPAddr", ":4000", "listen address of the aggreagtor HTTP server")
	GRPCListenAddr := flag.String("GRPCAddr", ":3001", "listen address of the aggreagtor GRPC server")
	flag.Parse()

	store := NewMemoryStore()
	var (
		svc Aggregator
	)
	svc = NewInvoiceAggregator(store)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)

	go makeGRPCTransport(*GRPCListenAddr, svc)
	makeHTTPTransport(*HTTPListenAddr, svc)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport running on port", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	server := grpc.NewServer()
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU ID"})
			return
		}

		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU ID"})
			return
		}

		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		WriteJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
