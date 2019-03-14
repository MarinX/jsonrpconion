package main

import (
	"fmt"

	"github.com/MarinX/jsonrpconion"
)

type Service struct{}

func (Service) Echo(in *string, out *string) error {
	*out = *in
	return nil
}

func main() {

	// creates new server
	srv := jsonrpconion.NewServer()

	// register our service Service
	if err := srv.Register(&Service{}); err != nil {
		fmt.Println(err)
		return
	}

	// Run without blocking and return onion v3 address
	addr, err := srv.RunNotify()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer srv.Stop()

	fmt.Println("Onion V3 address:", addr)

	// creates new client
	cli, err := jsonrpconion.NewClient(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cli.Close()

	// send requests to our service Service.Echo
	for i := 0; i < 5; i++ {
		out := new(string)
		in := fmt.Sprintf("Echo %d times", i+1)

		// call our service
		err := cli.Call("Service.Echo", &in, out)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(*out)
	}

}
