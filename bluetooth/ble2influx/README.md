

# Build and run 

```
go get ./...
go build main.go 
sudo ./main -user olivier
```

# Cross compile for Raspberry Pi 3b+

GOOS=linux GOARCH=arm GOARM=7 go build -o ble2influx-rpi ble2influx.go


# Some words about the protocol 

You'll find all the details on the Xiaomi Mijia alternate driver repository:
https://github.com/pvvx/ATC_MiThermometer#bluetooth-advertising-formats
