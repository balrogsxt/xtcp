package main

import (
	"fmt"
	"github.com/balrogsxt/xtcp/pack"
	"github.com/balrogsxt/xtcp/protocol"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"testing"
	"time"
)

//客户端部分测试入口
func TestClient(t *testing.T)  {
	tcpAddr,err := net.ResolveTCPAddr("tcp","127.0.0.1:6001")
	if err != nil {
		t.Fatalf("创建tcp地址失败: %s",err.Error())
	}
	conn,err := net.DialTCP("tcp4",nil,tcpAddr)
	if err != nil {
		t.Fatalf("连接tcp服务器失败: %s",err.Error())
	}

	go func() {
		for {
			datapack,err := pack.ParsePack(conn)
			if err != nil {
				if err == io.EOF {
					t.Logf("与服务器断开连接")
					break
				}
				t.Logf("发生错误: %s",err.Error())
				continue
			}
			//测试->解析为ProtoBuffer
			var a protocol.BroadcastMsg
			proto.Unmarshal(datapack.Data,&a)
			fmt.Println(fmt.Sprintf("收到服务端发来的消息: 内容:%s ，时间: %s",a.Content,time.Unix(a.SendTime,0).Format("2006-01-02 15:04:05")))

		}
	}()
	for {
		msg := protocol.BroadcastMsg{
			Content: "这是公开消息",
			SendTime: time.Now().Unix(),
		}
		datapack,err := pack.PackProtoBuffer(0x01,&msg)
		if err != nil {
			t.Fatalf("封包失败: %s",err.Error())
		}
		//发送封包数据
		conn.Write(datapack)
		time.Sleep(time.Second * 5)
	}

}