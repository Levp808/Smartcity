POSTGRES_USER:=postgres
POSTGRES_PASSWORD:=password
POSTGRES_DB_PETITION:=petition
POSTGRES_ADDRESS_PETITION:=0.0.0.0:5432
POSTGRES_ADDRESS_DEPART_COMM:=0.0.0.0:5431
POSTGRES_ADDRESS_DEPART_ROAD:=0.0.0.0:5430
POSTGRES_DB_DEPART:=report

LOCAL_DEPS:=$(CURDIR)/deps

run: get-deps build-app build-kafka up-kafka delay-kafka up-app delay-app migrate

stop:
	docker compose --file ./docker-compose.yml down --volumes --rmi all
	docker compose --file ./docker-compose-kafka.yml down --rmi all

get-deps:
	GOBIN=$(LOCAL_DEPS) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

build-app:
	docker-compose -f ./docker-compose.yml build

build-kafka:
	docker-compose -f ./docker-compose-kafka.yml build

up-kafka:
	docker-compose -f ./docker-compose-kafka.yml up -d

up-app:
	docker-compose -f ./docker-compose.yml up -d

migrate:
	$(LOCAL_DEPS)/migrate -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS_PETITION)/$(POSTGRES_DB_PETITION)?sslmode=disable' -path ./petition/migrations up
	$(LOCAL_DEPS)/migrate -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS_DEPART_COMM)/$(POSTGRES_DB_DEPART)?sslmode=disable' -path ./depart/migrations up
	$(LOCAL_DEPS)/migrate -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS_DEPART_ROAD)/$(POSTGRES_DB_DEPART)?sslmode=disable' -path ./depart/migrations up

delay-app:
	@if [ "$$(uname)" = "Linux" ] || [ "$$(uname)" = "Darwin" ]; then \
		sleep 20; \
	else \
		powershell -Command "Start-Sleep -Seconds 20"; \
	fi

delay-kafka:
	@if [ "$$(uname)" = "Linux" ] || [ "$$(uname)" = "Darwin" ]; then \
		sleep 20; \
	else \
		powershell -Command "Start-Sleep -Seconds 20"; \
	fi

.PHONY: run stop
