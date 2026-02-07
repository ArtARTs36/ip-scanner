SERVICE_NAME=ip-scanner
API_PATH 		   = api/grpc/${SERVICE_NAME}
PROTO_API_DIR 	   = api/grpc/${SERVICE_NAME}
PROTO_OUT_DIR 	   = pkg/${SERVICE_NAME}-grpc-api
PROTO_API_OUT_DIR  = ${PROTO_OUT_DIR}

.DEFAULT_GOAL := help

help: ## Show help
	@printf "\033[33m%s:\033[0m\n" 'Available commands'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_\/-]+:.*?## / {printf "  \033[32m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

init: ## Init service
	docker exec -it infra-postgres13 psql -c "CREATE DATABASE ${SERVICE_NAME}"
	make migration/migrate

tidy: ## Add missing and remove unused GO modules
	go mod tidy

gen/proto: ## Generate gRPC structures
	mkdir -p ${PROTO_OUT_DIR}
	protoc \
		-I ${API_PATH}/v1 \
		--include_imports \
		--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
        --go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
        --descriptor_set_out=$(PROTO_OUT_DIR)/api.pb \
		--go-client-builder_out=$(PROTO_OUT_DIR) \
            --go-client-builder_opt=generate-mock-client=true \
            --go-client-builder_opt=embed-client=true \
        ./${PROTO_API_DIR}/v1/*.proto

gen/go: ## Generate go/mock structures
	go generate ./...
	go generate ./pkg/${SERVICE_NAME}-grpc-api/mock_client_gen.go

gen: ## Generate go/mock, gRPC structures
	make gen/proto
	make gen/go
	go mod vendor

test: ## Run go tests
	go test ./...

lint: ## Run linter
	golangci-lint run --fix

up: ## Up services (foreground)
	docker-compose up

up-d: ## Up services (background)
	docker-compose up
