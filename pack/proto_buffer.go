package pack

import "google.golang.org/protobuf/proto"

//PackProtoBuffer 封包ProtoBuffer协议包
func PackProtoBuffer(typeId uint32,message proto.Message) ([]byte,error) {
	protoBytes,err := proto.Marshal(message)
	if err != nil {
		return nil,err
	}
	datapack,err := Pack(DataPack{
		Type: typeId,
		Len: uint32(len(protoBytes)),
		Data: protoBytes,
	})
	return datapack,err
}