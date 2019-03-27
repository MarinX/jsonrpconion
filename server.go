package jsonrpconion

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"
)

// Server represents an RPC Server on Tor.
type Server struct {
	rpcServer    *rpc.Server
	onionService *tor.OnionService
}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{
		rpcServer: rpc.NewServer(),
	}
}

// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
//	- exported method of exported type
//	- two arguments, both of exported type
//	- the second argument is a pointer
//	- one return value, of type error
// It returns an error if the receiver is not an exported type or has
// no suitable methods.
// The client accesses each method using a string of the form "Type.Method",
// where Type is the receiver's concrete type.
func (s *Server) Register(rcvr interface{}) error {
	return s.rpcServer.Register(rcvr)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func (s *Server) RegisterName(name string, rcvr interface{}) error {
	return s.rpcServer.RegisterName(name, rcvr)
}

// Run runs the server by establishing tor connection
// After connection to Tor is establish, it listens for incoming connection
func (s *Server) Run() error {
	if err := s.runOnion(); err != nil {
		return err
	}
	return s.runHandler()
}

// RunNotify is the same as Run but it does not block
// Returns Onion V3 address
func (s *Server) RunNotify() (string, error) {
	if err := s.runOnion(); err != nil {
		return "", err
	}
	go s.runHandler()
	return s.Addr()
}

// Stop stops the server by closing Tor connection
func (s *Server) Stop() {
	closeTorProcess()
	s.onionService = nil
}

// Addr returns Onion V3 address
func (s *Server) Addr() (string, error) {
	if s.onionService == nil {
		return "", ErrOnionServiceNotStarted
	}
	return fmt.Sprintf("%s.onion", s.onionService.ID), nil
}

// runOnion starts the Tor and Onion V3 service
func (s *Server) runOnion() error {
	if s.onionService != nil {
		return ErrOnionServiceAlreadyStarted
	}

	t, err := getTorProcess()
	if err != nil {
		return err
	}

	cfg := &tor.ListenConf{
		Version3:    true,
		RemotePorts: []int{80},
	}

	var inPathKey bool
	if keyExist() {
		inPathKey = true
		key, err := loadKey()
		if err == nil {
			cfg.Key = key
		}
	}

	s.onionService, err = t.Listen(nil, cfg)
	if !inPathKey {
		key := s.onionService.Key.(ed25519.KeyPair)
		saveKey(key.PrivateKey())
	}

	return err
}

// runHandler handles the underlaying Tor connection and serves json-rpc methods
func (s *Server) runHandler() error {
	for {
		conn, err := s.onionService.Accept()
		if err != nil {
			return err
		}
		go s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
