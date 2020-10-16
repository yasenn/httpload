package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type result struct {
	err        error
	httpCodeOk bool
	duration   time.Duration
}

type HttpLoad struct {
	Req         *http.Request
	LoadCount   int
	Concurrency int
	Client      *http.Client

	start time.Time
	end   time.Time

	results chan *result
}

var (
	argHttpVerb    = flag.String("b", "GET", "http verb: {'GET','POST'}")
	argConcurrency = flag.Int("c", 10, "concurrency")
	argCount       = flag.Int("n", 20, "Count of requests to make")
)

func (httpLoad *HttpLoad) Run() {
	httpLoad.Req.Header.Add("cache-control", "no-cache")
	if httpLoad.Client == nil {
		httpLoad.Client = &http.Client{}
	}
	httpLoad.results = make(chan *result, httpLoad.LoadCount)
	httpLoad.start = time.Now()

	requestsCount := httpLoad.LoadCount
	for {
		if requestsCount < 1 {
			break
		}

		c := httpLoad.Concurrency
		if requestsCount < httpLoad.Concurrency {
			c = requestsCount
		}

		var waitGroup sync.WaitGroup
		waitGroup.Add(c)
		for i := 0; i < c; i++ {
			go func() {
				httpLoad.request()
				waitGroup.Done()
			}()
		}
		waitGroup.Wait()
		requestsCount = requestsCount - c
	}

	httpLoad.end = time.Now()
	total := httpLoad.end.Sub(httpLoad.start)
	if total.Seconds() == 0 {
		fmt.Fprintf(os.Stderr, "error: zero time of HTTP load")
		os.Exit(1)
	}
	totalSuccessful := 0

	for {
		select {
		case responce := <-httpLoad.results:
			if !responce.httpCodeOk {
				continue
			}
			totalSuccessful++
		default:
			rps := float64(httpLoad.LoadCount) / total.Seconds()
			fmt.Printf("HTTP load results:\n")
			fmt.Printf("  Sent:\t%v requests in %v seconds with concurrency %v\n", httpLoad.LoadCount, total.Seconds(), httpLoad.Concurrency)
			fmt.Printf("  successful:\t%v responces\n", totalSuccessful)
			fmt.Printf("  requests/sec:\t%v\n", rps)
			return
		}
	}
}


func (httpLoad *HttpLoad) request() {
	s := time.Now()
	httpResponce, err := httpLoad.Client.Do(httpLoad.Req)

	httpCodeOk := httpResponce != nil && httpResponce.StatusCode >= 100 && httpResponce.StatusCode < 300
	httpLoad.results <- &result{httpCodeOk: httpCodeOk, err: err, duration: time.Now().Sub(s)}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: httpload <url>\n")
		os.Exit(1)
	}

	url := flag.Args()[0]
	httpVerb := strings.ToUpper(*argHttpVerb)
	httpRequest, err := http.NewRequest(httpVerb, url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening URL: %v", err)
		os.Exit(1)
	}

	loadCount := *argCount
	concurrency := *argConcurrency
	(&HttpLoad{LoadCount: loadCount, Concurrency: concurrency, Req: httpRequest}).Run()
}
