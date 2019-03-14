package jsonrpconion

import (
	"errors"
)

type EchoResp struct {
	ID     interface{} `json:"id"`
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
}

type Echo struct{}

func (Echo) Say(in *string, out *string) error {
	*out = *in
	return nil
}
func (Echo) Err(in *string, out *string) error {
	return errors.New("error")
}
