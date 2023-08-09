# Introduction

This project is to migreate the original amazeone-app from `docker-compose` to `kubernetes` to have an actual taste of real life microservices architecture

This project was made entirely on ubuntu which could make not work in other platforms without modification

# Pre-requisites

1. [Go v1.18.1](https://tip.golang.org/doc/go1.18)
2. [Docker Desktop](https://www.docker.com/products/docker-desktop/)
3. [Minikube](https://minikube.sigs.k8s.io/docs/start/)

# Instructions

1. Start docker desktop
2. Download and install minikube with other necessary tools

   ```console
   make setup-linux
   ```

3. Install dependencies

   ```console
   make dependencies
   ```

4. In folder `/k8s` create a `secrets.yaml` like the example file, this file contains all the needed secrets for the project to run

5. Run te project

   ```console
   make dev
   ```

6. Open a new terminal and run

   ```console
   minikube tunnel
   ```

> Important! since minikube don't have a way to create a tunnel in background the last command must keep up and running in order to develop correctly otherwise you will not be able to access kubernetes cluster

# Usage

This project counts with hot reload thanks to `skaffold` every time you make a change you will need to wait until the new build is complete and deployed before testing

To use the project run this command

```console
make dev
```

this will start all microservices in dev mode, with live realoding for any change made

In case you want to use the `development tools` used alongside this project you can run

```console
make dev-tools
```

This command will install all dev tools needed to test the project, I will let you a list next
| Tool | Description |
|----------|:-------------:|
| [Evans for gRPC](https://github.com/ktr0731/evans) | This tool acts as a gRPC client for you to test you API with no need of creating one in code

# Wiki!

Please don't forget to visit the [Wiki Page](https://github.com/Kiyosh31/e-commerce-microservice/wiki) to see diagrams, design docs and everything related to the design of this project! :)
