package jsonrpconion

import (
	"encoding/json"
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewServer()
	server.Register(new(Echo))

	cli, srv := net.Pipe()
	defer cli.Close()

	go server.rpcServer.ServeCodec(jsonrpc.NewServerCodec(srv))
	dec := json.NewDecoder(cli)

	cases := []string{"a", "b", "hello"}

	for i, cas := range cases {
		reply := new(EchoResp)

		fmt.Fprintf(cli, `{"method": "Echo.Say", "id": "%d", "params": ["%s"]}`, i+1, cas)
		err := dec.Decode(reply)
		if err != nil {
			t.Errorf("Decode: %s", err)
			return
		}

		if reply.Result != cas {
			t.Errorf("Response mismatch. Expect %s got %s\n", cas, *reply)
		}
	}
}
