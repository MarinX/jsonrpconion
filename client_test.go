package jsonrpconion

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"
)

func TestClient(t *testing.T) {
	rpc.Register(new(Echo))

	client := &Client{}

	cli, srv := net.Pipe()
	defer cli.Close()

	go rpc.ServeCodec(jsonrpc.NewServerCodec(srv))

	client.rpcClient = jsonrpc.NewClient(cli)

	cases := []string{"a", "b", "hello"}

	for _, cas := range cases {
		reply := new(string)
		err := client.Call("Echo.Say", &cas, reply)
		if err != nil {
			t.Error(err)
			continue
		}
		if *reply != cas {
			t.Errorf("Response mismatch. Expect %s got %s\n", cas, *reply)
		}
	}

	res := new(string)
	err := client.Call("Echo.Err", "Test", res)
	if err == nil {
		t.Error("Expecting Echo.Err to fail but it passed")
	}
}
