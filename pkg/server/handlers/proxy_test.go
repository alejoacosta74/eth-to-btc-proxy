package handlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	_ "github.com/alejoacosta74/qproxy/pkg/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestProxy(t *testing.T) {

	var (
		mu                  sync.Mutex
		forwardedHostHeader string
	)

	// create a backend server that checks the incoming headers
	// and echoes incoming requests
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		forwardedHostHeader = r.Header.Get("X-Forwarded-Host")
		w.WriteHeader(http.StatusOK)
		if r.Body != nil {
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)
			if len(body) > 0 {
				w.Write([]byte(`{"backend":` + string(body) + `}`))
			}
		}
	}))
	defer backendServer.Close()

	proxyHandler, err := NewProxyHandler(backendServer.URL, context.Background())
	if err != nil {
		t.Fatal(err)
	}
	reverseProxyServer := httptest.NewServer(proxyHandler)
	defer reverseProxyServer.Close()
	reverseProxyURL, err := url.Parse(reverseProxyServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("test GET verb", func(t *testing.T) {
		// call the reverse proxy
		_, err := http.Get(reverseProxyServer.URL)
		if err != nil {
			t.Fatal(err)
		}
		mu.Lock()
		got := forwardedHostHeader
		mu.Unlock()
		// check that the header X-Forwarded-Host has been set
		want := reverseProxyURL.Host
		if got != want {
			t.Errorf("GET %s gives header %s, got %s", reverseProxyServer.URL, want, got)
		}
	})

	t.Run("test POST verb", func(t *testing.T) {
		// call the reverse proxy
		msg := `{"msg":"hello world"}`
		resp, err := http.Post(reverseProxyServer.URL, "text/plain", ioutil.NopCloser(bytes.NewBufferString(msg)))
		if err != nil {
			t.Fatal(err)
		}
		assertResponseBody(t, resp, `{"backend":`+msg+`}`)
	})
}

func assertResponseBody(t *testing.T, resp *http.Response, want string) {
	t.Helper()
	if resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error("Error reading request body: ", err)
		}
		got := string(body)
		assert.Equal(t, want, got)
	}

}
