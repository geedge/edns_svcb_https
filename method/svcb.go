package method

import (
	"edns_svcb/utils"
	"fmt"
	"log"
	"sync"

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

func HTTPSRecordMeasurement(servers []string, domain_pool chan string, result_pool chan bool, routine_num int, wg *sync.WaitGroup, logger *log.Logger) {
	for i := 0; i < routine_num; i++ {
		wg.Add(1)
		go HTTPSRecordTestRoutine(servers, domain_pool, result_pool, wg, logger)
	}
}

func HTTPSRecordTestRoutine(server_list []string, domain_pool chan string, result_pool chan bool, wg *sync.WaitGroup, logger *log.Logger) {
	server_index := 0
	for {
		if domain, ok := <-domain_pool; ok {
			server := server_list[server_index] + ":53"
			res, err := HTTPSRecordTest(server, domain)
			if err == nil {
				logger.Printf("%v : %v", domain, res)
			} else {
				logger.Printf("%v : failed (%v) server %v", domain, err, server)
			}
			result_pool <- res
			server_index = (server_index + 1) % len(server_list)
		} else {
			break
		}
	}
	wg.Done()
}

func HTTPSRecordTest(addr, domain string) (bool, error) {
	res, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: domain, Qtype: dns.TypeHTTPS, RD: true})
	if err == nil {
		if len(res.Answer) > 0 {
			if _, ok := res.Answer[0].(*dns.HTTPS); ok {
				return true, nil
			}
		}
	}
	return false, err
}

func HTTPSRecordResultProcess(result_pool chan bool, wg *sync.WaitGroup) {
	total_cnt := 0
	valid_cnt := 0
	for {
		if res, ok := <-result_pool; ok {
			total_cnt++
			if res {
				valid_cnt++
			}
		} else {
			break
		}
	}
	fmt.Printf("Domain HTTPS support Test: total %v https %v\n", total_cnt, valid_cnt)
	wg.Done()
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
