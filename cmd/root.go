package cmd

import (
	"adjust/httpclient"
	"flag"
	"net/http"
	"time"
)

const (
	HttpTimeout = 5
)

func Execute() {
	workerCount := flag.Int("parallel", 10, "an int")
	flag.Parse()
	server := httpclient.NewServer(*workerCount, flag.Args(), &http.Client{Timeout: HttpTimeout * time.Second})
	server.Run()
}
