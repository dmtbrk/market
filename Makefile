all: build run

build:
	docker-compose build

run:
	docker-compose up

test:
	go test ./...

gen:
	go generate ./...

protoc:
	 protoc -I api/ api/product.proto --go_out=plugins=grpc:grpc --experimental_allow_proto3_optional

gqlgen:
	gqlgen generate
