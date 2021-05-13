package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testingHTTPClient(handler http.Handler, options ...bool) (*http.Client, func()) {
	s := httptest.NewServer(handler)
	fail := false
	if len(options) > 0 {
		fail = options[0]
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				if fail {
					return &net.TCPConn{}, nil
				}
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return client, s.Close
}

func TestReadData(t *testing.T) {
	assert := assert.New(t)

	testCfg := config{
		dump1090URL: &url.URL{Scheme: "http", Host: "example.com"},
	}

	t.Run("pass", func(t *testing.T) {
		h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			f, err := os.Open("./testdata/valid.json")
			if err != nil {
				t.Fatal(err)
			}
			io.Copy(rw, f)
			f.Close()
		})

		client, teardown := testingHTTPClient(h)
		defer teardown()

		_, err := readData(&testCfg, client)
		assert.Nil(err)
	})

	t.Run("http_fail", func(t *testing.T) {

		client, teardown := testingHTTPClient(nil, true)
		defer teardown()

		_, err := readData(&testCfg, client)
		if assert.NotNil(err) {
			assert.Contains(err.Error(), "retrieve")
		}
	})

	t.Run("http_status_fail", func(t *testing.T) {
		h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		})

		client, teardown := testingHTTPClient(h)
		defer teardown()

		_, err := readData(&testCfg, client)
		if assert.NotNil(err) {
			assert.Contains(err.Error(), "HTTP error code")
		}
	})

	t.Run("json_fail", func(t *testing.T) {
		h := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("not a json"))
		})

		client, teardown := testingHTTPClient(h)
		defer teardown()

		_, err := readData(&testCfg, client)
		if assert.NotNil(err) {
			assert.Contains(err.Error(), "unmarshal json")
		}
	})
}
