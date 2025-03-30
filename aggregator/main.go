package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qppffod/microservice-project/types"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		store          = makeStore()
		svc            Aggregator
		GRPCListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		HTTPListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewInvoiceAggregator(store)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)

	go makeGRPCTransport(GRPCListenAddr, svc)
	makeHTTPTransport(HTTPListenAddr, svc)
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
		if r.Method != http.MethodGet {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "method not allowed"})
			return
		}
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
		if r.Method != http.MethodPost {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "method not allowed"})
			return
		}
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

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatal("Invalid store type giver %s", storeType)
		return nil
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
