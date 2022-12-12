
# Project description

Ble2influx decodes Mijia Bluetooth frames and sends metrics to influxdb.

# Build

Just run :
```
$ make
```

Or do it by hand:

```
go get ./...
go build ble2influx.go
```

You can even cross compile for RPI with
```
$ make build-rpi
```

# Run

```
$ make
$ sudo ./build/ble2influx -user simpleuser
```
note: ble2influx needs root permissions to open /dev/hci device, but it'll drop to unprivileged if needed


# Some words about the Mijia protocol

You'll find all the details on the Xiaomi Mijia alternate driver repository:
https://github.com/pvvx/ATC_MiThermometer#bluetooth-advertising-formats
