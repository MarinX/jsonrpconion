package jsonrpconion

import "errors"

var (
	// ErrTorNotStarted returns if Tor service is not started
	ErrTorNotStarted = errors.New("error: Tor not started")
	// ErrTorAlreadyStarted returns if we have already establish the Tor connection
	ErrTorAlreadyStarted = errors.New("error: Tor already started")
	// ErrOnionServiceAlreadyStarted returns if the onion service is already started
	ErrOnionServiceAlreadyStarted = errors.New("error: Onion service already started")
	// ErrOnionServiceNotStarted returns if we try to run server without starting onion service
	ErrOnionServiceNotStarted = errors.New("error: Onion service not started")
)
