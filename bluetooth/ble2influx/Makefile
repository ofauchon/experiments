all: build-current

build-current:
	@echo "Building for current arch (build/ble2influx)"
	@mkdir -p build
	@go build -o build/ble2influx ble2influx.go 

build-rpi: 
	@echo "Building for Raspberry ARM arch"
	@mkdir -p build
	@GOOS=linux GOARCH=arm GOARM=7 go build -o build/ble2influx-rpi ble2influx.go

	
clean: 
	rm -fr ./build

