package method

import (
	"edns_svcb/utils"
	"fmt"
	"log"
	"sync"

	"github.com/miekg/dns"
)

type TestResult struct {
	Alive   bool
	Support bool
	Err     error
}

func EDNSSupportMeasurement(ip_pool chan string, result_pool chan TestResult, routine_num int, wg *sync.WaitGroup, logger *log.Logger) {
	for i := 0; i < routine_num; i++ {
		wg.Add(1)
		go EDNSSupportTestRoutine(ip_pool, result_pool, wg, logger)
	}
}

func EDNSSupportTestRoutine(ip_pool chan string, result_pool chan TestResult, wg *sync.WaitGroup, logger *log.Logger) {
	for {
		if ip, ok := <-ip_pool; ok {
			addr := ip + ":53"
			res := EDNSSupportTest(addr)
			result_pool <- res
			if res.Err != nil {
				logger.Printf("%v : failed ( %v )", ip, res.Err)
			} else {
				logger.Printf("%v : alive %v support %v", ip, res.Alive, res.Support)
			}
		} else {
			break
		}
	}
	wg.Done()
}

func EDNSSupportTest(addr string) TestResult {
	_, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: "www.example.net", RD: true})
	if err != nil {
		return TestResult{Err: err}
	}
	res, err := utils.DNSQuery(addr, utils.DNSOptions{Domain: "www.example.net", EDNS: true, RD: true})
	if err == nil {
		if len(res.Extra) == 1 {
			if _, ok := res.Extra[0].(*dns.OPT); ok {
				return TestResult{Alive: true, Support: true}
			}
		}
	}
	return TestResult{Alive: true}
}

func EDNSSupportResultProcess(result_pool chan TestResult, wg *sync.WaitGroup) {
	alive_cnt := 0
	valid_cnt := 0
	for {
		if res, ok := <-result_pool; ok {
			if res.Alive {
				alive_cnt++
			}
			if res.Support {
				valid_cnt++
			}
		} else {
			break
		}
	}
	fmt.Printf("EDNS support Test: alive %v support %v\n", alive_cnt, valid_cnt)
	wg.Done()
}
