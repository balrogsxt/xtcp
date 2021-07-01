package pack

import (
	"bytes"
	"encoding/binary"
	"net"
)

type DataPack struct {
	Type uint32 //uint32 -> 4字节
	Len uint32 //uint32 -> 4字节
	Data []byte
}
// GetHeadLen 获取包头字节 uint32 * 2
func GetHeadLen() uint32 {
	return 8
}
//封包顺序和解包顺序必须相同

//Pack 封包
func Pack(data DataPack) ([]byte,error) {
	buffer := bytes.NewBuffer([]byte{})
	//1.类型
	if err := binary.Write(buffer,binary.LittleEndian,data.Type);err != nil {
		return nil,err
	}
	//2.数据包长度
	if err := binary.Write(buffer,binary.LittleEndian,data.Len);err != nil {
		return nil,err
	}
	//3.数据包
	if err := binary.Write(buffer,binary.LittleEndian,data.Data);err != nil {
		return nil,err
	}
	return buffer.Bytes(),nil
}
//UnPack 拆包
func UnPack(data []byte) (*DataPack,error) {
	buffer := bytes.NewBuffer(data)
	datapack := &DataPack{}
	if err := binary.Read(buffer,binary.LittleEndian,&datapack.Type);err != nil {
		return nil,err
	}
	if err := binary.Read(buffer,binary.LittleEndian,&datapack.Len);err != nil {
		return nil,err
	}
	return datapack,nil
}

//ParsePack 快速从TCP连接数据中解包
func ParsePack(conn *net.TCPConn) (*DataPack,error) {
	//1.先读取包头字节获取type、length
	buffer := make([]byte,GetHeadLen())
	_,err := conn.Read(buffer)
	if err != nil {
		return nil,err
	}
	//解析出包头后,根据包头数据结果拿到数据包
	dataPack,err := UnPack(buffer)
	if err != nil { //拆包失败
		return nil,err
	}
	//拆包完成,再次读取数据包数据
	if dataPack.Len > 0 {
		//再次读取包头
		buffer := make([]byte,dataPack.Len)
		_,err := conn.Read(buffer)
		if err != nil {
			return nil,err
		}
		dataPack.Data = buffer
		return dataPack,nil
	}
	return nil,nil //无可用的数据包
}