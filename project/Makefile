FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp
MAILER_BINARY=mailServiceApp
LISTENER_BINARY=listenerApp
.DEFAULT_GOAL := all

## default
all: up_build start

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mailer build_listener
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## build_auth: builds the authentication binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"
build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"

build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"
build_mailer:
	@echo "Building mailer binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAILER_BINARY} ./cmd/api
	@echo "Done!"
## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

clean: stop down
	@echo "Deleting frontend"
	cd ../front-end && rm -rf ./${FRONT_END_BINARY}
	@echo "Deleting broker"
	cd ../broker-service && rm -rf ./${BROKER_BINARY}
	@echo "Deleting authentication"
	cd ../authentication-service && rm -rf ./${AUTH_BINARY}
	@echo "Deleting Logger"
	cd ../logger-service && rm -rf ./${LOGGER_BINARY}
	@echo "Deleting Mailer"
	cd ../mail-service && rm -rf ./${MAILER_BINARY}
	@echo "Deleting Listener"
	cd ../listener-service && rm -rf ./${LISTENER_BINARY}
	@echo "Delete DB"
	sudo rm -rf ./db-data
	mkdir -p db-data/postgres
	mkdir -p db-data/mongo

clean_cache: clean
	docker system prune -a --volumes
