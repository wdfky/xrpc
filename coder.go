package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

type RpcData struct {
	Name  string
	Args  []interface{}
	Reply []interface{}
	Err   error
}

//编码解码时加上包长度
func EnCode(rpcdata RpcData) ([]byte, error) {
	buf := bytes.Buffer{}
	encode := gob.NewEncoder(&buf)
	err := encode.Encode(rpcdata)
	if err != nil {
		return nil, err
	}
	message := buf.Bytes()
	pkg := new(bytes.Buffer)
	length := int32(len(message))
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(pkg, binary.LittleEndian, message)
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

func DeCode(reader *bufio.Reader) (rpcdata RpcData, err error) {
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err = binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return
	}
	// Buffered返回缓冲中现有的可读取的字节数。
	if int32(reader.Buffered()) < length+4 {
		return
	}

	// 读取真正的消息数据
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(pack[4:])
	decode := gob.NewDecoder(buf)
	err = decode.Decode(&rpcdata)
	if err != nil {
		return rpcdata, err
	}
	return rpcdata, nil
}

// func main() {
// 	rpcdata := &RpcData{
// 		Name: "1234",
// 	}
// 	data, err := EnCode(*rpcdata)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(data)
// 	r, err := DeCode(data)
// 	fmt.Println(r)
// }
