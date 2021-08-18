package main

import (
	"bufio"
	"net"
	"reflect"
)

type Server struct {
	addr string
	funs map[string]reflect.Value
}

// func a() {
// 	rpc.ServeConn()
// 	client, err := rpc.Dial("tcp", "localhost:1234")
// 	if err != nil {
// 		log.Fatal("dialing:", err)
// 	}

// 	var reply string
// 	err = client.Call("HelloService.Hello", "hello", &reply)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(reply)
// }
func (s *Server) NewServer(addr string) *Server {

	server := &Server{
		addr: addr,
		funs: make(map[string]reflect.Value),
	}
	return server
}
func (s *Server) RegisterName(name string, fun interface{}) {
	//map里面有说明已经注册过了
	if _, ok := s.funs[name]; ok {
		return
	}
	//注册
	s.funs[name] = reflect.ValueOf(fun)
}
func (s *Server) ServeConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		RpcReqData, err := DeCode(reader)
		if err != nil {
			return
		}
		fc, ok := s.funs[RpcReqData.Name]
		if !ok {
			return
		}
		args := make([]reflect.Value, 0, len(RpcReqData.Args))
		for _, arg := range RpcReqData.Args {
			args = append(args, reflect.ValueOf(arg))
		}
		//args := reflect.ValueOf(RpcReqData.Args)
		callResult := fc.Call(args)
		retvalue := make([]interface{}, 0, len(callResult))
		for _, rv := range callResult {
			retvalue = append(retvalue, rv.Interface())
		}
		RpcRetData := RpcData{
			Name:  RpcReqData.Name,
			Reply: retvalue,
			Err:   nil,
		}
		RetData, err := EnCode(RpcRetData)
		if err != nil {
			return
		}
		conn.Write(RetData)
	}
}
