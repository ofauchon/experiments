package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/miekg/dns"
)

var domainsToAddresses map[string]string = map[string]string{
	"somesite1.com.": "1.2.3.4",
	"somesite2.com.": "5.6.7.8",
}

var latencyMs *int
var tsigUser *string
var tsigKey *string

type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)

	fmt.Println("Query :", r.Question[0].Name)

	if tsig := r.IsTsig(); tsig != nil {
		fmt.Println("TSIG signature found")
		fmt.Println("algo:", tsig.Algorithm, "MAC:", tsig.MAC)
	}

	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := domainsToAddresses[domain]
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
		}
	}

	time.Sleep(time.Duration(*latencyMs) * time.Millisecond)

	if r.IsTsig() != nil {
		if w.TsigStatus() == nil {
			msg.SetTsig(r.Extra[len(r.Extra)-1].(*dns.TSIG).Hdr.Name,
				dns.HmacSHA256, 300, time.Now().Unix())
		} else {
			fmt.Println("Status", w.TsigStatus().Error())
		}
	}

	w.WriteMsg(&msg)
}

func main() {

	fmt.Println("Simple TSIG DNS Sever")

	addr := flag.String("bind", "127.0.0.1:1053", "Address to bind")
	tsigKey = flag.String("tsigkey", "", "TSIG Shared Key (base64)")
	tsigUser = flag.String("tsiguser", "", "TSIG User")
	latencyMs = flag.Int("latency", 0, "Response latency in ms ")

	flag.Parse()

	if *tsigKey == "" || *tsigUser == "" {
		log.Fatalf("TSIG key or user not defined\n")

	}

	fmt.Println("TSIG Key: ", *tsigKey)
	fmt.Println("TSIG User: ", *tsigUser)
	fmt.Println("Response Latency (ms): ", *latencyMs)
	srv := &dns.Server{Addr: *addr, Net: "udp"}
	srv.TsigSecret = map[string]string{*tsigUser + ".": *tsigKey}

	fmt.Println("Starting Server...")
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
	fmt.Println("DNS Server started")
}
