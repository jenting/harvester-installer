package console

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSSHKeysFromURL(t *testing.T) {
	keys, err := ioutil.ReadFile("testdata/keys")
	if err != nil {
		t.Fatalf("Fail to load fixture")
	}

	testCases := []struct {
		name         string
		httpResp     string
		pubKeysCount int
		expectError  string
	}{
		{
			name:         "Two public keys",
			httpResp:     string(keys),
			pubKeysCount: 2,
		},
		{
			name:        "Invalid public key",
			httpResp:    "\nooxx",
			expectError: "fail to parse on line 2: ooxx",
		},
		{
			name:        "No public key",
			httpResp:    "",
			expectError: "no key found",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, testCase.httpResp)
			}))
			defer ts.Close()

			pubKeys, err := getRemoteSSHKeys(ts.URL)
			if testCase.expectError != "" {
				assert.EqualError(t, err, testCase.expectError)
			} else {
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.pubKeysCount, len(pubKeys))
			}
		})
	}
}

func TestGetHarvesterManifestContent(t *testing.T) {
	d := map[string]string{
		"a": "b",
		"b": "\"c\"",
	}
	res := getHarvesterManifestContent(d)
	t.Log(res)
}

func TestGetHStatus(t *testing.T) {
	s := getHarvesterStatus()
	t.Log(s)
}

func TestGetFormattedServerURL(t *testing.T) {
	testCases := []struct {
		Name   string
		input  string
		output string
	}{
		{
			Name:   "ip",
			input:  "1.2.3.4",
			output: "https://1.2.3.4:6443",
		},
		{
			Name:   "domain name",
			input:  "example.org",
			output: "https://example.org:6443",
		},
		{
			Name:   "full",
			input:  "https://1.2.3.4:6443",
			output: "https://1.2.3.4:6443",
		},
	}
	for _, testCase := range testCases {
		got := getFormattedServerURL(testCase.input)
		assert.Equal(t, testCase.output, got)
	}
}

func TestF(t *testing.T) {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			if v, ok := addr.(*net.IPNet); ok && !v.IP.IsLoopback() && v.IP.To4() != nil {
				t.Log(v.IP.String())
			}
		}
	}
}

func TestGetServerURLFromEnvData(t *testing.T) {
	testCases := []struct {
		input []byte
		url   string
		err   error
	}{
		{
			input: []byte("K3S_CLUSTER_SECRET=abc\nK3S_URL=https://172.0.0.1:6443"),
			url:   "https://172.0.0.1:8443",
			err:   nil,
		},
		{
			input: []byte("K3S_CLUSTER_SECRET=abc\nK3S_URL=https://172.0.0.1:6443\nK3S_NODE_NAME=abc"),
			url:   "https://172.0.0.1:8443",
			err:   nil,
		},
	}

	for _, testCase := range testCases {
		url, err := getServerURLFromEnvData(testCase.input)
		assert.Equal(t, testCase.url, url)
		assert.Equal(t, testCase.err, err)
	}
}
