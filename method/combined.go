package method

import (
	"edns_svcb/utils"
	"fmt"
	"log"
	"sync"

	"github.com/miekg/dns"
)

type CombinedResult struct {
	Alive        bool
	EDNSSupport  bool
	SVCBSupport  bool
	HTTPSSupport bool
	Err          error
}

func CombinedMeasurement(ip_pool chan string, result_pool chan CombinedResult, routine_num int, wg *sync.WaitGroup, logger *log.Logger) {
	for i := 0; i < routine_num; i++ {
		wg.Add(1)
		go CombinedTestRoutine(ip_pool, result_pool, wg, logger)
	}
}

func CombinedTestRoutine(ip_pool chan string, result_pool chan CombinedResult, wg *sync.WaitGroup, logger *log.Logger) {
	for {
		if ip, ok := <-ip_pool; ok {
			addr := ip + ":53"
			res := CombinedTest(addr)
			result_pool <- res
			if res.Err != nil {
				logger.Printf("%v : failed ( %v )", ip, res.Err)
			} else {
				logger.Printf("%v : alive %v edns %v svcb %v https %v", ip, res.Alive, res.EDNSSupport, res.SVCBSupport, res.HTTPSSupport)
			}
		} else {
			break
		}
	}
	wg.Done()
}

func CombinedTest(addr string) CombinedResult {
	result := CombinedResult{}
	_, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: "www.example.com", RD: true})
	if err != nil {
		return CombinedResult{Err: err}
	}
	result.Alive = true

	//edns support
	res, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: "www.example.com", EDNS: true, RD: true})
	if err == nil {
		if len(res.Extra) > 0 {
			if _, ok := res.Extra[0].(*dns.OPT); ok {
				result.EDNSSupport = true
			}
		}
	} else {
		result.Err = err
		return result
	}

	// svcb support
	res, err = utils.DNSQuery(addr, utils.DNSOptions{Domain: "_dns.resolver.arpa", Qtype: dns.TypeSVCB, RD: true})
	if err == nil {
		if len(res.Answer) > 0 {
			if _, ok := res.Answer[0].(*dns.SVCB); ok {
				result.SVCBSupport = true
			}
		}
	} else {
		result.Err = err
		return result
	}

	// https support
	res, err = utils.DNSQuery(addr, utils.DNSOptions{Domain: "blog.cloudflare.com", Qtype: dns.TypeHTTPS, RD: true})
	if err == nil {
		if len(res.Answer) > 0 {
			if _, ok := res.Answer[0].(*dns.HTTPS); ok {
				result.HTTPSSupport = true
			}
		}
	} else {
		result.Err = err
		return result
	}

	return result
}

func CombinedResultProcess(result_pool chan CombinedResult, wg *sync.WaitGroup) {
	alive_cnt := 0
	edns_cnt := 0
	svcb_cnt := 0
	https_cnt := 0
	for {
		if res, ok := <-result_pool; ok {
			if res.Err != nil {
				continue
			}
			if res.Alive {
				alive_cnt++
			}
			if res.EDNSSupport {
				edns_cnt++
			}
			if res.SVCBSupport {
				svcb_cnt++
			}
			if res.HTTPSSupport {
				https_cnt++
			}
		} else {
			break
		}
	}
	fmt.Printf("EDNS support Test: alive %v edns %v svcb %v https %v\n", alive_cnt, edns_cnt, svcb_cnt, https_cnt)
	wg.Done()
}
