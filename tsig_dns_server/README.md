* TSIG DNS Server

This sample dns server responds DNS Queries with TSIG signatures.
You can add latency too

Please note it's still work in progress


* Run the server

```
$ go run tsig_dns_server.go --help
Simple TSIG DNS Sever
Usage of ./tsig_dns_server:
  -bind string
        Address to bind (default "127.0.0.1:1053")
  -latency int
        Response latency in ms 
  -tsigkey string
        TSIG Shared Key (base64)
  -tsiguser string
        TSIG User
```

```
$ go run tsig_dns_server.go  -tsiguser tsig1 -tsigkey rzerRR444343FFSfsmGhmersmcMRdfmSFDm324234G8= -latency 10
```

* Test with dig: 

```
$ dig -y hmac-sha256:tsig1:rzerRR444343FFSfsmGhmersmcMRdfmSFDm324234G8= -p 1053 A somesite1.com   @localhost 

; <<>> DiG 9.16.4 <<>> -y hmac-sha256 -p 1053 A somesite1.com @localhost
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 43102
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;somesite1.com.                 IN      A

;; ANSWER SECTION:
somesite1.com.          60      IN      A       1.2.3.4

;; TSIG PSEUDOSECTION:
tsig1.                  0       ANY     TSIG    hmac-sha256. 1597752356 300 32 5AO5nEzuxnkZjqZ0ODxl1NhLgdCCtcpyZh9VK9dNWok= 43102 NOERROR 0 

;; Query time: 120 msec
;; SERVER: 127.0.0.1#1053(127.0.0.1)
;; WHEN: Tue Aug 18 14:05:56 CEST 2020
;; MSG SIZE  rcvd: 138
```

