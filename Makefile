obu:
	@go build -o bin/obu obu/*.go
	@./bin/obu

receiver:
	@go build -o bin/rec ./data_receiver
	@./bin/rec

calculator:
	@go build -o bin/calculator ./distance_calculator
	@./bin/calculator

agg:
	@go build -o bin/agg ./aggregator
	@./bin/agg

.PHONY: obu