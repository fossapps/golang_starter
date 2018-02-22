build:
	go build -tags dev -o .build/Server cmd/server/main.go

build-prod:
	go build -ldflags="-s -w" -o .build/ProductionServer -tags prod cmd/server/main.go

build-staging:
	go build -tags staging -o .build/StagingServer cmd/server/main.go

serve:
	go run -tags dev cmd/server/main.go

build-cron:
	go build -tags dev -o .build/ProcessRegistrationQueue cmd/cron/ProcessRegistrationQueue/main.go

seed:
	go run -tags dev cmd/seeding/main.go

test:
	go test -tags dev ./...

install:
	sudo mkdir -p /opt/api_server
	sudo systemctl stop api.service
	sudo cp ./main /opt/api_server/server
	sudo cp ./.build/api.service /etc/systemd/system/api.service
	sudo systemctl daemon-reload
	sudo systemctl start api.service
clean:
	rm -f .build/Server .build/ProductionServer .build/ProcessRegistrationQueue
