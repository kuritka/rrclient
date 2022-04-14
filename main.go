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
	IPs    []string `json:"ip"`
}

func appendOpt(msg *dns.Msg,rr RoundRobinState) {
	json, _ := json.Marshal(rr)
	opt := new(dns.OPT)
	opt.Hdr.Name = "."
	opt.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_LOCAL)
	e.Code = dns.EDNS0LOCALSTART //.EDNS0LOCALSTART
	e.Data = append([]byte("_rr_state="),json...)
	opt.Option = append(opt.Option, e)
	msg.Extra = append(msg.Extra, opt)
}

func requestStateless(rr RoundRobinState){
	msg := new(dns.Msg)
	msg.SetQuestion("cloud.example.com.", dns.TypeA)
	appendOpt(msg, rr)
	result, err := dns.Exchange(msg, fmt.Sprintf("%s:%v", server, port))
	kingpin.FatalIfError(err,"Check if CoreDNS is running on port %s", port)
	str := fmt.Sprintf("\n-------------- https://onlinestringtools.com/convert-bytes-to-string -------\n%s\n-------------\n%s",msg.String(), result.String())
	rr = RoundRobinState{}
	for _, a := range result.Answer {
		switch a.Header().Rrtype{
		case dns.TypeA:
			rr.IPs = append(rr.IPs, a.(*dns.A).A.String())
		case dns.TypeAAAA:
			rr.IPs = append(rr.IPs, a.(*dns.AAAA).AAAA.String())
		}
	}
	fmt.Println(str)
	completed <- rr
}

func requestStateful(){
	msg := new(dns.Msg)
	msg.SetQuestion("test.example.com.", dns.TypeA)
	result, err := dns.Exchange(msg, fmt.Sprintf("%s:%v", server, port))
	kingpin.FatalIfError(err,"Check if CoreDNS is running on port %s", port)
	fmt.Println(result)
	completed <- RoundRobinState{}
}


var completed chan RoundRobinState

// Round Robin EDNS0 client
func main(){
	completed = make(chan RoundRobinState)
	//var rr = RoundRobinState{}
	for {
		go requestStateful()
		<- completed
		//rr = <-completed
		fmt.Println("Press the Enter to continue")
		_,_ = fmt.Scanln()
	}
	close(completed)
}



