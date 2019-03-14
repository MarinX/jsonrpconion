package jsonrpconion

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// Client represents an RPC Client.
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
type Client struct {
	address   string
	rpcClient *rpc.Client
}

// NewClient returns a new Client to handle requests to the
// set of services at the other end of the Tor connection.
func NewClient(onion string) (*Client, error) {
	cli := &Client{address: onion}
	if err := cli.createClient(); err != nil {
		return nil, err
	}
	return cli, nil
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	return c.rpcClient.Call(serviceMethod, args, reply)
}

// Close rpc client and Tor connection
func (c *Client) Close() error {
	var err error
	if c.rpcClient != nil {
		err = c.rpcClient.Close()
	}
	closeTorProcess()
	return err
}

// createClient creates a dialer connection through Tor network
func (c *Client) createClient() error {
	t, err := getTorProcess()
	if err != nil {
		return err
	}

	dialer, err := t.Dialer(nil, nil)
	if err != nil {
		return err
	}
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:80", c.address))
	if err != nil {
		return err
	}
	c.rpcClient = jsonrpc.NewClient(conn)
	return err
}
