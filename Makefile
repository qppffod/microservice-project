obu:
	@go build -o bin/obu obu/*.go
	@./bin/obu

receiver:
	@go build -o bin/rec data_receiver/main.go
	@./bin/rec

.PHONY: obu