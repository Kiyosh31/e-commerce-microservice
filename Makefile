dev:
	skaffold dev

dependencies:
	@cd customer-service
	@go mod tidy

common:
	chmod +x update-common.sh
	./update-common.sh

setup-linux:
	chmod +x setup-linux.sh
	./setup-linux.sh

config-linux:
	minikube start
	eval $(minikube docker-env)
	minikube addons enable ingress
	minikube addons enable ingress-dns
	sudo chmod 666 /var/run/docker.sock