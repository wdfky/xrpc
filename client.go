package main

import (
	"bufio"
	"net"
	"reflect"
)

type Client struct {
	Conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}

// func (c *Client) Call(name string, args interface{}, reply interface{}) error {
// 	reflect.TypeOf(args)
// 	ReqRpcData := RpcData{
// 		Name:  name,
// 		Args:  args,
// 		Reply: reply,
// 	}
// 	data, err := EnCode(ReqRpcData)
// 	if err != nil {
// 		return err
// 	}
// 	c.Conn.Write(data)
// 	reader := bufio.NewReader(c.Conn)
// 	RetRpcData, err := DeCode(reader)
// 	if err != nil {
// 		return err
// 	}
// 	reply = RetRpcData.Reply
// }
func (c *Client) Call(name string, fun interface{}) {
	fn := reflect.ValueOf(fun).Elem()
	f := func(args []reflect.Value) (results []reflect.Value) {
		inArgs := make([]interface{}, 0, len(args))
		for _, v := range args {
			inArgs = append(inArgs, v.Interface())
		}
		rpcReqData := RpcData{
			Name: name,
			Args: inArgs,
		}
		reqdata, err := EnCode(rpcReqData)
		if err != nil {
			return
		}
		_, err = c.Conn.Write(reqdata)
		if err != nil {
			return
		}
		reader := bufio.NewReader(c.Conn)
		rpcRetData, err := DeCode(reader)
		if err != nil {
			return
		}
		outArgs := make([]reflect.Value, 0, len(rpcRetData.Args))
		for i, v := range rpcRetData.Args {
			if v == nil {
				outArgs = append(outArgs, reflect.Zero(fn.Type().Out(i)))
				continue
			}
			outArgs = append(outArgs, reflect.ValueOf(v))
		}
		return outArgs
	}
	v := reflect.MakeFunc(fn.Type(), f)
	fn.Set(v)
}
