package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/alejoacosta74/rpc-proxy/pkg/log"
	"github.com/pkg/errors"
	// log "github.com/alejoacosta74/gologger"
)

type ProxyHandler struct {
	backend *url.URL
	ctx     context.Context
}

func NewProxyHandler(backendUrl string, ctx context.Context) (*ProxyHandler, error) {
	backend, err := url.Parse(backendUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing backend url")
	}
	return &ProxyHandler{
		backend,
		ctx,
	}, nil
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reverseProxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Host = h.backend.Host
			r.URL.Path = "/"
			r.URL.Scheme = h.backend.Scheme
			r.Header.Set("X-Forwarded-Host", r.Host)
			r.Host = h.backend.Host
			readRequest(r)
		},
	}
	reverseProxy.ModifyResponse = readResponse()
	reverseProxy.ErrorHandler = errorHandler()
	r = r.WithContext(h.ctx)
	reverseProxy.ServeHTTP(w, r)
}

func readRequest(req *http.Request) {
	log.With("module", "proxy").Tracef("req.remoteAddr: %s, \t req.Method %s, \t req.URL %s \t req.Host: %s", req.RemoteAddr, req.Method, req.URL, req.Host)
	log.With("module", "proxy").Tracef("req.Header: %s", req.Header)

	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.With("module", "proxy").Debugf("Error reading request body: %s", err.Error())
		}
		if log.IsDebug() {
			logPretty("==> Request Body", body)
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		log.With("module", "proxy").Debugf("Error while modifying response %s ", err.Error())
	}
}

func readResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		if resp.Body != nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.With("module", "proxy").Debugf("Error reading response body: %s ", err.Error())
				return err
			}
			if log.IsDebug() {
				logPretty("<== Response body", body)
			}

			resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		return nil
	}
}

func logPretty(msg string, output []byte) {
	if len(output) > 0 {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, output, "", "    "); err != nil {
			fmt.Printf("Error decoding JSON: %v\n", err)
		} else {
			fmt.Printf("\n%s :\n%s\n", msg, prettyJSON.String())

		}
	}
}

// parseToUrl parses a "to" address to url.URL value
func parseToUrl(addr string) *url.URL {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}
	toUrl, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	return toUrl
}
