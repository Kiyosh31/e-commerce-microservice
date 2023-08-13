SWAGGER_FOLDER = docs/swagger
CUSTOMER_SERVICE_PROTO_DIR = customer-service/proto
CUSTOMER_SERVICE_PB_DIR = customer-service/proto/pb

dev:
	skaffold dev

start-minikube:
	minikube start
	minikube tunnel

clean-docker:
	minikube ssh -- docker system prune
	docker system prune --force --all

dependencies:
	cd customer-service
	go mod download

common:
	chmod +x update-common.sh
	./update-common.sh

setup-linux:
	chmod +x setup-linux.sh
	./setup-linux.sh

dev-tools:
	chmod +x linux-dev-tools.sh
	./linux-dev-tools.sh

customer-proto:
	rm -rf ${CUSTOMER_SERVICE_PB_DIR}
	mkdir ${CUSTOMER_SERVICE_PB_DIR}
	rm -rf ${SWAGGER_FOLDER}
	mkdir ${SWAGGER_FOLDER}
	protoc --proto_path=${CUSTOMER_SERVICE_PROTO_DIR} --go_out=${CUSTOMER_SERVICE_PB_DIR} --go_opt=paths=source_relative \
    --go-grpc_out=${CUSTOMER_SERVICE_PB_DIR} --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=${CUSTOMER_SERVICE_PB_DIR} --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=${SWAGGER_FOLDER} --openapiv2_opt=allow_merge=true,merge_file_name=customer-service \
    ${CUSTOMER_SERVICE_PROTO_DIR}/*.proto