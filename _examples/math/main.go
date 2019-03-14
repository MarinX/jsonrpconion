package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MarinX/jsonrpconion"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func main() {

	server := jsonrpconion.NewServer()

	server.Register(new(Arith))

	fmt.Println("Starting JSON RPC server on Tor...this will take a while")
	addr, err := server.RunNotify()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Server started. Onion address %s\n", addr)
	fmt.Println("Now you can call RPC methods Arith.Multiply and Arith.Divide")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigc

	fmt.Println("Closing server...")
	server.Stop()
	fmt.Println("Bye, bye!")

}
