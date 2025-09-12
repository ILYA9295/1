proto:
	protoc --go_out=. --go-grpc_out=. nlp_service/nlp.proto

run:
	go run ./cmd/main.go

docker-up:
	docker-compose up --build
