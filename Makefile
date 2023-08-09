CUSTOMER_SERVICE_PROTO_DIR = customer-service/proto
CUSTOMER_SERVICE_USER_PB_DIR = customer-service/proto/pb

dev:
	skaffold dev

dependencies:
	

common:
	chmod +x update-common.sh
	./update-common.sh

setup-linux:
	chmod +x setup-linux.sh
	./setup-linux.sh

customer-proto:
	rm -rf ${CUSTOMER_SERVICE_USER_PB_DIR}
	mkdir ${CUSTOMER_SERVICE_USER_PB_DIR}
	protoc --proto_path=${CUSTOMER_SERVICE_PROTO_DIR} --go_out=${CUSTOMER_SERVICE_USER_PB_DIR} --go_opt=paths=source_relative \
    --go-grpc_out=${CUSTOMER_SERVICE_USER_PB_DIR} --go-grpc_opt=paths=source_relative \
    ${CUSTOMER_SERVICE_PROTO_DIR}/*.proto