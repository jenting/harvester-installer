package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPClient interface {
	Get(url string, timeout time.Duration) ([]byte, error)
}

type DefaultHTTPClient struct{}

func (DefaultHTTPClient) Get(url string, timeout time.Duration) ([]byte, error) {
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("got %d status code from %s, body: %s", resp.StatusCode, url, string(body))
	}

	return body, nil
}
