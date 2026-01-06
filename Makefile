RPI_GOOS ?= linux
RPI_GOARCH ?= arm
RPI_GOARM ?= 7

build:
	@go build -o bin/ecom ./cmd

build-rpi:
	@GOOS=$(RPI_GOOS) GOARCH=$(RPI_GOARCH) GOARM=$(RPI_GOARM) go build -o bin/ecom ./cmd

# Build optimized untuk Raspberry Pi (stripped binary, smaller size)
build-rpi-optimized:
	@echo "Building optimized binary for Raspberry Pi..."
	@GOOS=$(RPI_GOOS) GOARCH=$(RPI_GOARCH) GOARM=$(RPI_GOARM) \
		go build -ldflags="-s -w" -trimpath -o bin/ecom-rpi ./cmd
	@echo "✓ Optimized binary created: bin/ecom-rpi"

# Build dengan profiling enabled untuk analisis performa
build-rpi-profile:
	@echo "Building profiling binary for Raspberry Pi..."
	@GOOS=$(RPI_GOOS) GOARCH=$(RPI_GOARCH) GOARM=$(RPI_GOARM) \
		go build -ldflags="-s -w" -trimpath -race -o bin/ecom-rpi-profile ./cmd
	@echo "✓ Profiling binary created: bin/ecom-rpi-profile"

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

# Clean build artifacts
clean:
	@rm -rf bin/
	@echo "✓ Build artifacts cleaned"

# Build dan run dengan monitoring enabled
run-monitor: build
	@ENABLE_RESOURCE_MONITOR=true ./bin/ecom

