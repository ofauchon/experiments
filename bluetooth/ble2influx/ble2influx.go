package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/pkg/errors"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

type MijiaDeviceConfig struct {
	Mac   string `json:"mac"`
	Model string `json:"model"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
}

var nodeConfigJson string = `[
{"mac": "a4c1381c0390","model": "xiaomi_mijia","name": "mijia_salon1","desc":"Capteur dans le salon"},
{"mac": "a4c1386b6dc6","model": "xiaomi_mijia","name": "mijia_capteur1","desc":"Capteur 1"},
{"mac": "a4c1382b1044","model": "xiaomi_mijia","name": "mijia_capteur2","desc":"Capteur 2"},
{"mac": "a4c138e73381","model": "xiaomi_mijia","name": "mijia_capteur3","desc":"Capteur 3"},
{"mac": "a4c138634fff","model": "xiaomi_mijia","name": "mijia_capteur4","desc":"Capteur 4"},
{"mac": "a4c1382c1d48","model": "xiaomi_mijia","name": "mijia_capteur5","desc":"Capteur 5"},
{"mac": "a4c138a23543","model": "xiaomi_mijia","name": "mijia_capteur6","desc":"Capteur 6"},
{"mac": "a4c138a09c8e","model": "xiaomi_mijia","name": "mijia_capteur7","desc":"Capteur 7"},
{"mac": "a4c1384f7a79","model": "xiaomi_mijia","name": "mijia_capteur8","desc":"Capteur 8"}
]`
var mijiaConfig = make([]MijiaDeviceConfig, 0)

var (
	// Args
	influx_server       = flag.String("influx_server", "http://localhost:8086", "Sets the influxDB server")
	influx_token        = flag.String("influx_token", "", "Sets the influxDB token")
	influx_org          = flag.String("influx_org", "your.org", "Sets the influxDB organization")
	influx_bucket       = flag.String("influx_bucket", "default", "Sets the influxDB bucket")
	influx_measurement  = flag.String("influx_measurement", "metro", "Sets the influxDB measurement name")
	dropuser            = flag.String("user", "default", "Drop privileges to <user>")
	device              = flag.String("device", "", "implementation of ble")
	sensors_descriptor  = flag.String("desc", "~/.config/ble2influx/sensors.json", "Sensors descriptor file")
	influx_only_connect = flag.Bool("influx-only-connect", false, "Connect InfluxDB without pushing metrics")
	period              = flag.Int("period", 60, "Duration (in sec) between two influxdB metrics updates")
	debug               = flag.Bool("debug", false, "Enable debug")

	// InfluxDB2
	client   influxdb2.Client
	writeAPI api.WriteAPI
)

// Structs and Decoder for Mijia ATC advertisements (custom firmware)
type MijiaMetrics struct {
	Mac        [6]byte
	Temp       float32
	Humi       float32
	Batt       float32
	RSSI       int
	FrameCount uint8
}

// Map for latest measurement and its Mutex
var lastMetrics = make(map[[6]byte]*MijiaMetrics)
var lockMetrics = sync.RWMutex{}

var lastUpload = time.Now()

/*
 * decodeMijia decodes BLE adv payload
 */
func decodeMijia(dat []byte) (*MijiaMetrics, error) {

	if *debug {
		hex := hex.EncodeToString(dat)
		fmt.Println("RX: %d bytes [%s]", hex)
	}

	if len(dat) != 15 {
		return nil, errors.New("Bad packet length")
	}

	ret := &MijiaMetrics{}

	for i := 0; i < 6; i++ {
		ret.Mac[i] = dat[5-i]
	}

	// See https://www.fanjoe.be/?p=3911 for negative temperatures
	ret.Temp = float32(uint32(dat[7])*0xFF+uint32(dat[6])) / 100
	ret.Humi = float32(uint32(dat[9])*0xFF+uint32(dat[8])) / 100
	ret.Batt = float32(uint32(dat[11])*0xFF+uint32(dat[10])) / 1000
	ret.FrameCount = dat[12]
	return ret, nil
}

// Helper function to drop privileges after we are bind to HCI device
func chuser(username string) (uid, gid int) {
	usr, err := user.Lookup(username)
	if err != nil {
		fmt.Printf("failed to find user %q: %s\n", username, err)
		os.Exit(3)
	}

	uid, err = strconv.Atoi(usr.Uid)

	if err != nil {
		fmt.Printf("bad user ID %q: %s\n", usr.Uid, err)
		os.Exit(3)
	}

	gid, err = strconv.Atoi(usr.Gid)

	if err != nil {
		fmt.Printf("bad group ID %q: %s", usr.Gid, err)
		os.Exit(3)
	}

	if err := syscall.Setgid(gid); err != nil {
		fmt.Printf("setgid(%d): %s", gid, err)
		os.Exit(3)
	}

	if err := syscall.Setuid(uid); err != nil {
		fmt.Printf("setuid(%d): %s", uid, err)
		os.Exit(3)
	}

	return uid, gid
}

/*
 * influxSender is a goroutine for sending Mijia metrics to influxdb
 * It assumes InfluxDB connexion is ok
 */
func influxSender(metrics map[[6]byte]*MijiaMetrics, dryRun bool) {

	cnt := int(0)
	for {
		cnt = 0
		if dryRun == true {
			fmt.Println("Sending influxdb metrics disabled (only_connect) ")
		} else {

			for mac, data := range metrics {
				hs := hex.EncodeToString(data.Mac[:])

				// Get name from json config
				dName := "unknown"
				for _, s := range mijiaConfig {
					if s.Mac == hs {
						dName = s.Name
					}
				}

				if *debug {
					fmt.Printf("TX %s: Name:%s Rssi:%d Temp:%.2f Humi:%.2f Batt:%.2f Frame:%d\n", hs, dName, data.RSSI, data.Temp, data.Humi, data.Batt, data.FrameCount)
				}
				p := influxdb2.NewPoint(*influx_measurement,
					map[string]string{"type": "mijia", "source": hs},
					map[string]interface{}{"rssi": data.RSSI, "temp": data.Temp, "humi": data.Humi, "batt": data.Batt, "name": dName}, time.Now())
				cnt++
				writeAPI.WritePoint(p)
				lockMetrics.Lock()
				delete(metrics, mac)
				lockMetrics.Unlock()
			}
			writeAPI.Flush()
			lastUpload = time.Now()
		}

		time.Sleep(time.Duration(*period) * time.Second)
		fmt.Println("influxSender: Sent ", cnt, " measurements, now sleeping ", *period, "sec.")
	}
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}

/*
 * advHandler processes BLE ads
 */
func advHandler(a ble.Advertisement) {

	if len(a.ServiceData()) > 0 {
		for _, svc := range a.ServiceData() {

			// Discard if it's not Mijia (0x181a)
			if !svc.UUID.Equal(ble.UUID16(0x181a)) {
				fmt.Println("Skipping UUID", svc.UUID)
				continue
			}

			// Try to decode payload
			mi, err := decodeMijia(svc.Data)
			if err == nil {
				hs := hex.EncodeToString(mi.Mac[:])
				if *debug {
					fmt.Printf("RX %s: Rssi:%d Temp:%.2f Humi:%.2f Batt:%.2f Frame:%d\n", hs, a.RSSI(), mi.Temp, mi.Humi, mi.Batt, mi.FrameCount)
				}
				mi.RSSI = a.RSSI()
				lockMetrics.Lock()
				lastMetrics[mi.Mac] = mi
				lockMetrics.Unlock()

			} else {
				fmt.Println("Bad Mijia payload")
			}

		}

	}
}

// MAIN
func main() {
	fmt.Println("Starting ble2influx")
	flag.Parse()

	fmt.Println("Creating BLE device")
	d, err := dev.NewDevice(*device)
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}

	//fmt.Println("Switching to ", *dropuser, " user")
	//chuser(*dropuser)

	// Try to read default
	if _, err := os.Stat(*sensors_descriptor); err == nil {
		fmt.Println("Using json sensor descriptor file ", sensors_descriptor)

		// Mijia configuration json
		err = json.Unmarshal([]byte(nodeConfigJson), &mijiaConfig)
		if err != nil {
			log.Fatalf("Can't parse json file :", err)
		}
		fmt.Println("Mijia json configuration : ", len(mijiaConfig), " entries found")
	}

	//InfluxDB connection
	fmt.Println("Connecting to influxDB server")
	fmt.Println("  server", *influx_server, " bucket:", *influx_bucket, " org:", *influx_org)
	ble.SetDefaultDevice(d)
	client = influxdb2.NewClient(*influx_server, *influx_token)
	defer client.Close()
	writeAPI = client.WriteAPI(*influx_org, *influx_bucket)

	// Run routine for sending Mijia metrics
	go influxSender(lastMetrics, *influx_only_connect)

	// Scan forever, or until interrupted by user.
	fmt.Println("Starting BLE Advertisement Listener")
	ctx := ble.WithSigHandler(context.Background(), nil)
	chkErr(ble.Scan(ctx, true, advHandler, nil))
}
