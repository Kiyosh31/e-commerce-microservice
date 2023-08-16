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


SWAGGER_FOLDER = docs/swagger

CUSTOMER_SERVICE_PROTO_DIR = customer-service/proto
CUSTOMER_SERVICE_PB_DIR = customer-service/proto/pb
CUSTOMER_SERVICE_SWAGGER_FILE = docs/swagger/customer.swagger.json

INVENTORY_SERVICE_PROTO_DIR = inventory-service/proto
INVENTORY_SERVICE_PB_DIR = inventory-service/proto/pb
INVENTORY_SERVICE_SWAGGER_FILE = docs/swagger/inventory.swagger.json

customer-proto:
	rm -rf ${CUSTOMER_SERVICE_PB_DIR}
	mkdir ${CUSTOMER_SERVICE_PB_DIR}
	rm -f ${CUSTOMER_SERVICE_SWAGGER_FILE}
	protoc --proto_path=${CUSTOMER_SERVICE_PROTO_DIR} --go_out=${CUSTOMER_SERVICE_PB_DIR} --go_opt=paths=source_relative \
		--go-grpc_out=${CUSTOMER_SERVICE_PB_DIR} --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=${CUSTOMER_SERVICE_PB_DIR} --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=${SWAGGER_FOLDER} \
		${CUSTOMER_SERVICE_PROTO_DIR}/*.proto

inventory-proto:
	rm -rf ${INVENTORY_SERVICE_PB_DIR}
	mkdir ${INVENTORY_SERVICE_PB_DIR}
	rm -f ${INVENTORY_SERVICE_SWAGGER_FILE}
	protoc --proto_path=${INVENTORY_SERVICE_PROTO_DIR} --go_out=${INVENTORY_SERVICE_PB_DIR} --go_opt=paths=source_relative \
		--go-grpc_out=${INVENTORY_SERVICE_PB_DIR} --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=${INVENTORY_SERVICE_PB_DIR} --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=${SWAGGER_FOLDER} \
		${INVENTORY_SERVICE_PROTO_DIR}/*.proto