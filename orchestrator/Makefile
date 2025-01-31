SERVER_NAME=grpc_server

PROTO_DIR=api/orchestrator_v1
OUT_DIR=pkg/grpc/pb

ifeq ($(OS),Windows_NT)
	RM = del /Q
else
	RM = rm -rf
endif

GO_DEPS=google.golang.org/grpc \
         google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc

deps:
	@echo "Installing dependencies..."
	go get -u $(GO_DEPS)

generate-proto:
	@echo "Generating gRPC code from proto files..."
	protoc	-I ./api \
			--go_out .\
			--go-grpc_out .\
			./$(PROTO_DIR)/orchestrator.proto

# Команда для сборки сервера
build:
	@echo "Building the server..."
	go build -o $(SERVER_NAME) orchestrator/cmd/main/.

# Запуск сервера
run: build
	@echo "Running the server..."
	./$(SERVER_NAME)

docker-run:
	@echo "Running the server..."
	docker-compose up -d

docker-stop:
	@echo "Stopping the server..."
	docker compose --file ./docker-compose.yml down
# Очистка сгенерированных файлов и исполняемого файла
clean:
	@echo "Cleaning up..."
	$(RM) "$(OUT_DIR)\petition.pb.go"
	$(RM) "$(OUT_DIR)\petition_grpc.pb.go"
	$(RM) "$(SERVER_NAME)"