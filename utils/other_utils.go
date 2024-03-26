package utils

import (
	"bufio"
	_ "encoding/json"
	_ "fmt"
	"io"
	_ "net"
	_ "net/http"
	"os"
	_ "strconv"
	"strings"
)

func RetrieveLines(pool chan string, filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	//fmt.Println("reading file ...")
	reader := bufio.NewReader(f)
	for {
		s, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		s = strings.Trim(s, "\n")
		pool <- s
	}
	close(pool)
}

func RetrieveLinesToSlice(filename string) []string {
	results := []string{}
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	for {
		s, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		s = strings.Trim(s, "\n")
		results = append(results, s)
	}
	return results
}
