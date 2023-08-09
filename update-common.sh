#!/bin/bash

cd customer-service
go get github.com/Kiyosh31/e-commerce-microservice-common
go get -u=patch github.com/Kiyosh31/e-commerce-microservice-common
go mod tidy
