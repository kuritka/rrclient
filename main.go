package main

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"gopkg.in/alecthomas/kingpin.v2"
)
const (
	port int =  5056
	server string = "localhost"
)


type RoundRobinState struct {
	Family uint16   `json:"type"`
	IPs    []string `json:"ip"`
}

func appendOpt(msg *dns.Msg) {
	var rr = RoundRobinState{
		Family: dns.TypeA,
		IPs:    []string{"10.0.0.1","10.0.0.2","10.0.0.3"},
	}
	json, _ := json.Marshal(rr)
	opt := new(dns.OPT)
	opt.Hdr.Name = "."
	opt.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_LOCAL)
	e.Code = dns.EDNS0LOCALSTART +10//.EDNS0LOCALSTART
	e.Data = append([]byte("_rr_state="),json...)
	opt.Option = append(opt.Option, e)
	msg.Extra = append(msg.Extra, opt)
}

func request(){
	msg := new(dns.Msg)
	msg.SetQuestion("cloud.example.com.", dns.TypeA)
	appendOpt(msg)
	fmt.Println("\n--------------")
	fmt.Println("online bytes converter: https://onlinestringtools.com/convert-bytes-to-string")
	fmt.Println(msg)
	fmt.Println("\n-------------")
	result, err := dns.Exchange(msg, fmt.Sprintf("%s:%v", server, port))
	kingpin.FatalIfError(err,"Check if CoreDNS is running on port %s", port)
	fmt.Println(result)
}

// Round Robin EDNS0 client
func main(){
	request()
}