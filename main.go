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
	flag.StringVar(&input_file, "input", "", "input file name")
	flag.Parse()
	input_pool := make(chan string, 500)
	result_pool := make(chan method.CombinedResult)
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
	method.CombinedMeasurement(input_pool, result_pool, 500, &test_tasks, logger)
	process_tasks.Add(1)
	go method.CombinedResultProcess(result_pool, &process_tasks)
	test_tasks.Wait()
	close(result_pool)
	process_tasks.Wait()
	time.Sleep(time.Second)
}
