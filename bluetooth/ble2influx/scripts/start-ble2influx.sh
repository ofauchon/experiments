#!/bin/sh
influx_token="YOUR INFLUXDB TOKEN"
influx_server="http://localhost:8086"
influx_measurement="YOUR MEASUREMENT"
influx_org="YOUR ORGANISATION"
influx_bucket="YOUR BUCKET"

#ble2influx will drop root privileges, and switch back to USER
usr=$USER

sudo ./ble2influx \
	-influx_server ${influx_server} \
	-influx_token=${influx_token} \
	-influx_measurement=${influx_measurement} \
	-influx_bucket=${influx_bucket} \
	-influx_org=${influx_org} \
	-user=${usr}

.
