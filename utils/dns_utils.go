// dns utils
package utils

import (
	"github.com/miekg/dns"
)

type DNSOptions struct {
	Domain string
	RD     bool
	Qclass uint16
	Qtype  uint16
	EDNS   bool
}

// build the question section of a dns packet
func questionMaker(domain string, qclass uint16, qtype uint16) *dns.Question {
	return &dns.Question{Name: dns.Fqdn(domain), Qtype: qtype, Qclass: qclass}
}

// build a specific query message
func queryMaker(domain string, rd bool, qclass uint16, qtype uint16, edns bool) *dns.Msg {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = rd
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = *questionMaker(domain, qclass, qtype)
	if edns {
		msg = msg.SetEdns0(1232, false)
	}
	return msg
}

// query and receive the response
// addr must contain dest port
// func DNSQuery(addr string, domain string, rd bool, qclass uint16, qtype uint16, edns bool) (*dns.Msg, error) {
func DNSQuery(addr string, opt DNSOptions) (*dns.Msg, error) {
	if opt.Qclass == 0 {
		opt.Qclass = 1
	}
	if opt.Qtype == 0 {
		opt.Qtype = 1
	}
	msg := queryMaker(opt.Domain, opt.RD, opt.Qclass, opt.Qtype, opt.EDNS)
	res, err := dns.Exchange(msg, addr)
	return res, err
}

func AsyncDNSQuery(addr string, domain string, rd bool, qclass uint16, qtype uint16, edns bool) error {
	msg := queryMaker(domain, rd, qclass, qtype, edns)
	client := dns.Client{Net: "udp"}
	conn, err := client.Dial(addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	err = conn.WriteMsg(msg)
	if err != nil {
		return err
	}
	return nil
}
