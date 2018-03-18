#!make
include .env
export

build:
	go build -o .build/Server cmd/server/main.go

serve:
	go run cmd/server/main.go

build-cron:
	go build -o .build/ProcessRegistrationQueue cmd/cron/ProcessRegistrationQueue/main.go

migrate:
	go run cmd/migrations/main.go

test:
	go test ./...

test-cover:
	go test -coverprofile=c.out ./... && go tool cover -html=c.out && rm c.out

test-unit:
	go test ./...

test-integration:
	go test -cover -tags integration ./...

install:
	sudo mkdir -p /opt/api_server
	sudo systemctl stop api.service
	sudo cp ./main /opt/api_server/server
	sudo cp ./.build/api.service /etc/systemd/system/api.service
	sudo systemctl daemon-reload
	sudo systemctl start api.service
clean:
	rm -f .build/Server .build/ProductionServer .build/ProcessRegistrationQueue
