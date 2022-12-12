
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


# Sensor json descriptor file (optional)

```
[
        {"mac": "a4c1381c1390","model": "xiaomi_mijia","name": "mijia_outside","desc":"Sensor outside"},
        {"mac": "a4c1386b3dc6","model": "xiaomi_mijia","name": "mijia_room1","desc":"Sensor 1"},
        {"mac": "a4c1382b4044","model": "xiaomi_mijia","name": "mijia_room2","desc":"Sensor 2"},
]
```

# Some words about the Mijia protocol

You'll find all the details on the Xiaomi Mijia alternate driver repository:
https://github.com/pvvx/ATC_MiThermometer#bluetooth-advertising-formats
