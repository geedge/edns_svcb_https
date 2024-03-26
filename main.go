package main

import (
	"edns_svcb/method"
	"edns_svcb/utils"
	"flag"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	var input_file string
	var routine_num int
	var mode string
	flag.StringVar(&input_file, "input", "", "input file name")
	flag.IntVar(&routine_num, "routine", 500, "working goroutine number")
	flag.StringVar(&mode, "mode", "resolver", "probe mode (resolver & domain)")
	flag.Parse()
	input_pool := make(chan string, 500)
	result_pool := make(chan method.CombinedResult, 500)
	bool_pool := make(chan bool, 500)
	var test_tasks sync.WaitGroup
	var process_tasks sync.WaitGroup

	logfilename := "test/" + time.Now().Format("2006-01-02_15-04-05") + ".log"
	logfile, err := os.OpenFile(logfilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to create log file: ", err)
	}
	defer logfile.Close()
	logger := log.New(logfile, "", log.LstdFlags)

	go utils.RetrieveLines(input_pool, input_file)
	if mode == "resolver" {
		method.CombinedMeasurement(input_pool, result_pool, routine_num, &test_tasks, logger)
	} else if mode == "domain" {
		servers := utils.RetrieveLinesToSlice("data/server_list")
		method.HTTPSRecordMeasurement(servers, input_pool, bool_pool, routine_num, &test_tasks, logger)
	}
	process_tasks.Add(1)
	if mode == "resolver" {
		go method.CombinedResultProcess(result_pool, &process_tasks)
	} else if mode == "domain" {
		go method.HTTPSRecordResultProcess(bool_pool, &process_tasks)
	}
	test_tasks.Wait()
	if mode == "resolver" {
		close(result_pool)
	} else if mode == "domain" {
		close(bool_pool)
	}
	process_tasks.Wait()
	time.Sleep(time.Second)
}
