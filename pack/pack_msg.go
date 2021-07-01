package pack

import (
	"bytes"
	"encoding/binary"
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