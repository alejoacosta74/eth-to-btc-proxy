package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/alejoacosta74/rpc-proxy/pkg/handlers"
	"github.com/alejoacosta74/rpc-proxy/pkg/rpc"

	"github.com/alejoacosta74/rpc-proxy/pkg/log"
	"github.com/alejoacosta74/rpc-proxy/pkg/qtum"

	"github.com/gorilla/mux"
)

type Server struct {
	server  *http.Server
	address string
}

func NewServer(localAddress string, backendUrl string, qcli *qtum.QtumClient, network string) (*Server, error) {
	ctx := context.Background()

	router := mux.NewRouter()

	//Create new RPC service and assign /rpc the endpoint
	rpcService, err := rpc.NewEthereumRPCService(network, qcli)
	if err != nil {
		return nil, err
	}

	router.Handle("/rpc", rpcService).Methods("POST")

	//Create new proxy handler and assign /proxy the endpoint
	proxyHandler, err := handlers.NewProxyHandler(backendUrl, ctx)
	if err != nil {
		return nil, err
	}
	router.Handle("/proxy", proxyHandler)

	server := &http.Server{
		Addr: localAddress,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	return &Server{
		server:  server,
		address: localAddress,
	}, nil
}

func (s *Server) Start() error {
	log.With("module", "server").Infof("Starting server on port: %s", s.address)
	log.With("module", "server").Infof("proxy available on: %s ", s.address+"/proxy")
	log.With("module", "server").Infof("eth jsonrpc server available on: %s ", s.address+"/rpc")
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	log.With("module", "server").Infof("Server shutdown gracefully")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	} else {
		log.With("module", "server").Infof("Server stopped")
	}
	return nil
}
