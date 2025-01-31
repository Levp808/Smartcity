POSTGRES_USER:=postgres
POSTGRES_PASSWORD:=password
POSTGRES_DB:=depart
POSTGRES_ADDRESS:=0.0.0.0:5433

build:
	docker-compose -f ./docker-compose-depart.yml build

up:
	docker-compose -f ./docker-compose-depart.yml up -d

migrate:up
	migrate -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DB)?sslmode=disable' -path ./migrations up

down:
	docker-compose -f ./docker-compose-petition.yml down

generate-proto:
	mkdir -p pkg/depart_v1
	protoc --proto_path ../api/depart_v1 \
	--go_out=pkg/depart_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/depart_v1 --go-grpc_opt=paths=source_relative \
	../api/depart_v1/depart.proto

	mkdir -p internal/grpc/clients/petition/pb
	protoc --proto_path ../api/petition_v1 --proto_path vendor.protogen \
	--go_out=internal/grpc/clients/petition/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/grpc/clients/petition/pb --go-grpc_opt=paths=source_relative \
	../api/petition_v1/petition.proto

vendor-proto:
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi

.PHONY: build up down generate-proto