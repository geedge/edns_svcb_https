package method

import (
	"edns_svcb/utils"

	"github.com/miekg/dns"
)

func HTTPSSupportTest(addr string) TestResult {
	_, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: "www.example.net", RD: true})
	if err != nil {
		return TestResult{Err: err}
	}
	res, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: "blog.cloudflare.com", Qtype: dns.TypeHTTPS, RD: true})
	if err == nil {
		if len(res.Answer) == 1 {
			if _, ok := res.Extra[0].(*dns.HTTPS); ok {
				return TestResult{Alive: true, Support: true}
			}
		}
	}
	return TestResult{Alive: true}
}

func SVCBSupportTest(addr string) TestResult {
	_, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: "www.example.net", RD: true})
	if err != nil {
		return TestResult{Err: err}
	}
	res, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: "blog.cloudflare.com", Qtype: dns.TypeSVCB, RD: true})
	if err == nil {
		if len(res.Answer) == 1 {
			if _, ok := res.Extra[0].(*dns.SVCB); ok {
				return TestResult{Alive: true, Support: true}
			}
		}
	}
	return TestResult{Alive: true}
}

func HTTPSRecordTest(addr, domain string) bool {
	res, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: domain, Qtype: dns.TypeHTTPS, RD: true})
	if err == nil {
		if len(res.Answer) > 0 {
			if _, ok := res.Answer[0].(*dns.HTTPS); ok {
				return true
			}
		}
	}
	return false
}

func SVCBRecordTest(addr, domain string) bool {
	res, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: domain, Qtype: dns.TypeSVCB, RD: true})
	if err == nil {
		if len(res.Answer) > 0 {
			if _, ok := res.Answer[0].(*dns.SVCB); ok {
				return true
			}
		}
	}
	return false
}
