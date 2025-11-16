.PHONY: run
.PHONY: proto
.PHONY: new-service

run:
	go run ./services/$(SERVICE)/main.go

new-service:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please provide a service name. Usage: make new-service name=<service-name>"; \
		echo "Example: make new-service name=payment-service"; \
		exit 1; \
	fi
	@echo "Generating new service: $(name)"
	@if ! go run core/boilerplate/generate_service.go $(name); then \
		echo "Error: Failed to generate service."; \
		exit 1; \
	fi
	@echo "Service '$(name)' successfully created!"

SERVICES := $(shell find ./services -mindepth 1 -maxdepth 1 -type d)

proto:
	@for dir in $(SERVICES); do \
		if ls $$dir/api/grpc/pb/*.proto 1> /dev/null 2>&1; then \
			mkdir -p $$dir/api/grpc/pb/src/golang; \
			mkdir -p $$dir/docs; \
			protoc -I $$dir/api/grpc/pb -I core/grpc/proto/googleapis \
				-I core/grpc/proto \
				--go_out=paths=source_relative:$$dir/api/grpc/pb/src/golang \
				--go-grpc_out=paths=source_relative:$$dir/api/grpc/pb/src/golang \
				--grpc-gateway_out=paths=source_relative:$$dir/api/grpc/pb/src/golang \
				--openapiv2_out=$$dir/docs \
				$$dir/api/grpc/pb/*.proto; \
		fi; \
	done



print:
	@for dir in $(SERVICES); do \
		echo "Folder: $$dir"; \
	done