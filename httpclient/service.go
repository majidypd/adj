package httpclient

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/goware/urlx"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Server struct {
	Client      HTTPClient
	WorkerCount int
	Args        []string
	ArgsLen     int
}

type Result struct {
	Url string
	Md5 string
}

// NewServer server construct
func NewServer(wc int, a []string, c HTTPClient) *Server {
	return &Server{
		Client:      c,
		WorkerCount: wc,
		Args:        a,
		ArgsLen:     len(a),
	}
}

// Run entrypoint of script by starting worker pool
func (s *Server) Run() {
	numUrls := len(flag.Args())
	urlsChan := make(chan string, numUrls)
	resultChan := make(chan Result, numUrls)

	for w := 1; w <= s.WorkerCount; w++ {
		go s.Worker(urlsChan, resultChan)
	}

	for j := 0; j < numUrls; j++ {
		urlsChan <- s.Args[j]
	}
	close(urlsChan)

	for a := 1; a <= numUrls; a++ {
		res := <-resultChan
		fmt.Printf("%s %s \n", res.Url, res.Md5)
	}
}

// Worker each worker for handle individual task
func (s *Server) Worker(urlsChan <-chan string, resultChan chan<- Result) {
	for url := range urlsChan {
		var err error

		// Normalize Url
		normalizeUrl, err := s.NormalizeUrl(url)
		if err != nil {
			log.Println(err)
			resultChan <- Result{}
			continue
		}

		// Send http request
		resp, err := s.SendHttpRequest(normalizeUrl)
		if err != nil {
			log.Println(err)
			resultChan <- Result{}
			continue
		}

		// make MD5 from response
		md5, err := s.Md5(resp.Body)
		if err != nil {
			log.Println(err)
			resultChan <- Result{}
			continue
		}

		resultChan <- Result{normalizeUrl, md5}
	}
}

func (s *Server)NormalizeUrl(url string) (string, error) {
	u, err := urlx.ParseWithDefaultScheme(url, "https")
	if err != nil {
		return "", err
	}
	normalized, err := urlx.Normalize(u)
	if err != nil {
		return "", err
	}
	return normalized, err
}

// Md5 make md5 hash
func (s *Server) Md5(resp io.Reader) (string, error) {
	buf, err := ioutil.ReadAll(resp)
	if err != nil {
		return "", err
	}
	hash := md5.Sum(buf)
	md := hex.EncodeToString(hash[:])
	return md, nil
}

// SendHttpRequest sends a Get request to the URL with the body
func (s *Server) SendHttpRequest(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return s.Client.Do(request)
}
