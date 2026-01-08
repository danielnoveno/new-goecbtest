RPI_GOOS ?= linux
RPI_GOARCH ?= arm
RPI_GOARM ?= 7

build:
	@go build -o bin/ecom ./cmd

build-rpi:
	@GOOS=$(RPI_GOOS) GOARCH=$(RPI_GOARCH) GOARM=$(RPI_GOARM) go build -o bin/ecom ./cmd

build-rpi-optimized:
	@echo "optimasi binary Raspberry Pi.."
	@GOOS=$(RPI_GOOS) GOARCH=$(RPI_GOARCH) GOARM=$(RPI_GOARM) \
		go build -ldflags="-s -w" -trimpath -o bin/ecom-rpi ./cmd
	@echo "optimasi binary Raspberry Pi selesai"

build-rpi-profile:
	@GOOS=$(RPI_GOOS) GOARCH=$(RPI_GOARCH) GOARM=$(RPI_GOARM) \
		go build -ldflags="-s -w" -trimpath -race -o bin/ecom-rpi-profile ./cmd

test:
	@go test -v ./...

run: build
	@./bin/ecom

# Contoh: make migration create_users_table
migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

fyne-build:
	@mkdir -p bin
	@cd cmd && fyne package -os windows -icon logo.png -name "ECB Test"
	@mv "cmd/ECB Test.exe" bin/

fyne-build-rpi:
	@mkdir -p bin
	@cd cmd && GOOS=$(RPI_GOOS) GOARCH=$(RPI_GOARCH) GOARM=$(RPI_GOARM) fyne package -os linux -arch arm -icon logo.png -name "ECB Test"
	@mv "cmd/ECB Test.tar.xz" bin/

clean:
	@rm -rf bin/

run-monitor: build
	@ENABLE_RESOURCE_MONITOR=true ./bin/ecom

