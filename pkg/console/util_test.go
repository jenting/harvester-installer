package console

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	resp []byte
	err  error
}

func (c MockHTTPClient) Get(url string, timeout time.Duration) ([]byte, error) {
	return c.resp, c.err
}

func TestGetSSHKeysFromURL(t *testing.T) {
	validPubKeys := []string{
		`ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDbeUa9A7Kee+hcCleIXYxuaPksn2m4PZTd4T7wPcse8KbsQfttGRax6vxQXoPO6ehddqOb2nV7tkW2mEhR50OE7W7ngDHbzK2OneAyONYF44bmMsapNAGvnsBKe9rNrev1iVBwOjtmyVLhnLrJIX+2+3T3yauxdu+pmBsnD5OIKUrBrN1sdwW0rA2rHDiSnzXHNQM3m02aY6mlagdQ/Ovh96h05QFCHYxBc6oE/mIeFRaNifa4GU/oELn3a6HfbETeBQz+XOEN+IrLpnZO9riGyzsZroB/Y3Ju+cJxH06U0B7xwJCRmWZjuvfFQUP7RIJD1gRGZzmf3h8+F+oidkO2i5rbT57NaYSqkdVvR6RidVLWEzURZIGbtHjSPCi4kqD05ua8r/7CC0PvxQb1O5ILEdyJr2ZmzhF6VjjgmyrmSmt/yRq8MQtGQxyKXZhJqlPYho4d5SrHi5iGT2PvgDQaWch0I3ndEicaaPDZJHWBxVsCVAe44Wtj9g3LzXkyu3k= root@admin`,
		`ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCnFYqVCT0/QFyIHdcW9K5V2jrpxlhJ603kdY1jPvubdz1aPPln0tQdr+wQnuq3pUloe7yW97+TFylP6G6ztYeWA6ckmPZaxnEoNc9KbDcn47T3bVSGhd3u9TpfgeARKk/fpQUIxF7GhV0tM5ltdRvdNT8PGv5wP/hozykbCyss9PbP82QZQ82h+GlQGbswaRFqty+cKXSOehwVHwgB9GNrUYetanLdJVtoTn7imub/zI2aQ8YZIqxVnvjDL7ftVy/mxDWoVwzBAcLZTY0MCS4bn0164wayTtMU1WXf5WgxasPPhW2EmZe9VNOfVzCqkx/YbrE/mqW74x28qYLr3KtqniBWVyGfUB6j0n35CjJTsZ7/gpy1bOWz6/siAs69xetcONVe6IA8MnX8lSEEiewwzDT4a21P0OZ8itNvWzgDQ6fvIMeX5pWj8wIok1TeuzXRIMGlrScmRUTT6/u3BJjAB9r8zxqCZT7MssSCoU8KsAnCZ4SRnKeIpNixjuKTDrM= root@k3os-2234`,
	}

	testCases := []struct {
		httpResp     string
		pubKeysCount int
	}{
		{
			httpResp:     strings.Join(validPubKeys, "\n"),
			pubKeysCount: 2,
		},
		{
			httpResp:     "abc",
			pubKeysCount: 0,
		},
		{
			httpResp:     "",
			pubKeysCount: 0,
		},
	}

	for _, testCase := range testCases {
		mockClient := MockHTTPClient{resp: []byte(testCase.httpResp), err: nil}
		pubKeys, err := getSSHKeysFromURL(mockClient, "https://example.com/keys")
		if testCase.pubKeysCount != 0 {
			assert.Equal(t, nil, err)
			assert.Equal(t, testCase.pubKeysCount, len(pubKeys))
		} else {
			assert.EqualError(t, err, "ssh: no key found")
		}
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
