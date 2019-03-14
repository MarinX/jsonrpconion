# jsonrpconion

[![Build Status](https://travis-ci.org/MarinX/jsonrpconion.svg?branch=master)](https://travis-ci.org/MarinX/jsonrpconion)
[![Go Report Card](https://goreportcard.com/badge/github.com/MarinX/jsonrpconion)](https://goreportcard.com/report/github.com/MarinX/jsonrpconion)
[![GoDoc](https://godoc.org/github.com/MarinX/jsonrpconion?status.svg)](https://godoc.org/github.com/MarinX/jsonrpconion)
[![License MIT](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](LICENSE)

Library for building JSON RPC services on Tor network

## Install
```sh
go get github.com/MarinX/jsonrpconion
```

## Usage
Take a look at `_examples` directory which contains a simple `echo` server/client and `math` which contains JSON RPC server

### Server
Create a simple JSON RPC server on Tor.

```go
package main

import (
	"fmt"
	"time"

	"github.com/MarinX/jsonrpconion"
)

type Service struct{}

func (Service) Say(in *string, out *string) error {
	*out = fmt.Sprintf("Hello, %s", *in)
	return nil
}

func main() {

	// creates new server
	srv := jsonrpconion.NewServer()

	// register our service Service
	if err := srv.Register(new(Service)); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			addr, err := srv.Addr()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			fmt.Println("Onion V3 address:", addr)
			break
		}
	}()

	// Run with blocking
	// If you want to run a server without blocking
	// call srv.RunNotify() which will return onion address
	if err := srv.Run(); err != nil {
		fmt.Println(err)
		return
	}
	defer srv.Stop()
}
```

### Client
Create a simple client to call onion JSON RPC endpoint

```go
package main

import (
	"fmt"

	"github.com/MarinX/jsonrpconion"
)

type Service struct{}

func (Service) Say(in *string, out *string) error {
	*out = fmt.Sprintf("Hello, %s", *in)
	return nil
}

func main() {

	addr := "onion_v3_address"

	// creates new server
	cli, err := jsonrpconion.NewClient(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cli.Close()

	reply := new(string)

	// calls our simple RPC service
	err = cli.Call("Service.Say", "Marin", reply)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Should return Hello, Marin
	fmt.Println("Reply from JSON RPC endpoint:", *reply)

}

```

## Test
```sh
go test -v
```

## License
MIT