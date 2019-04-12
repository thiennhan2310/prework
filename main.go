package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

type respInfo struct {
	Status bool // e.g. "true | false"
}

type summaryInfo struct {
	succes int
	fail   int
}

func main() {
	requests := flag.Int("n", 1, "Number of requests to perform")
	concurrency := flag.Int("c", 1, "Number of multiple requests to make at a time")
	link := flag.String("url", "", "Url to test")
	flag.Parse()

	if *requests == 0 || *concurrency == 0 || *link == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	if *requests < *concurrency {
		fmt.Println("Number of request must be greater than or equal  number of request make same time")
		os.Exit(-1)
	}

	chanResp := make(chan respInfo)
	requested := 0
	responded := 0

	start := time.Now()
	for r := 0; r < *requests; r += *concurrency {
		for c := 0; requested < *requests && c < *concurrency; c++ {
			go testurl(*link, chanResp)
			requested++
		}
	}

	summary := summaryInfo{succes: 0, fail: 0}
	for resp := range chanResp {
		responded++
		if resp.Status == true {
			summary.succes++
		} else if resp.Status == false {
			summary.fail++
		}

		if responded == requested {
			fmt.Println("Success ", summary.succes)
			fmt.Println("Fail ", summary.fail)
			fmt.Println("Take ", time.Now().Sub(start))
			break
		}
	}

}

func testurl(url string, c chan respInfo) {
	_, err := http.Get(url)

	if err != nil {
		c <- respInfo{false}
		return
	}
	c <- respInfo{true}
	return
}
