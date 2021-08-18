package main

import (
	"net"
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
func (c *Client) Call(name string)
