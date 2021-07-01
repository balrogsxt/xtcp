package main

import (
	"fmt"
	"github.com/balrogsxt/xtcp/pack"
	"github.com/balrogsxt/xtcp/protocol"
	"github.com/balrogsxt/xtcp/server"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)
type TestListener struct {}

func (c *TestListener) OnOpen(conn *server.TcpConnection)  {
	//给所有客户端发送上线通知
	b_bytes,_ := proto.Marshal(&protocol.BroadcastMsg{
		Content: fmt.Sprintf("%d -> %s上线了",conn.GetFd(),conn.GetRemoteAddr().String()),
		SendTime: time.Now().Unix(),
	})
	conn.GetServer().Broadcast(b_bytes)
	fmt.Println(conn.GetRemoteAddr(),"客户端已连接",conn.GetFd())
}
func (c *TestListener) OnClose(conn *server.TcpConnection)  {
	b_bytes,_ := proto.Marshal(&protocol.BroadcastMsg{
		Content: fmt.Sprintf("%d -> %s离开了",conn.GetFd(),conn.GetRemoteAddr().String()),
		SendTime: time.Now().Unix(),
	})
	conn.GetServer().Broadcast(b_bytes)
	fmt.Println(conn.GetRemoteAddr(),"客户端已断开",conn.GetFd())
}
func (c *TestListener) OnMessage(conn *server.TcpConnection,data *pack.DataPack)  {

	var a protocol.BroadcastMsg
	proto.Unmarshal(data.Data,&a)
	fmt.Println(conn.GetRemoteAddr(),conn.GetFd(),"客户端发的消息:",a.Content,",时间:",time.Unix(a.SendTime,0).Format(
	"2006-01-02 15:04:05"))

	//发送pb包
	b := protocol.BroadcastMsg{
		Content: "服务端的消息",
		SendTime: time.Now().Unix(),
	}
	b_bytes,_ := proto.Marshal(&b)


	conn.Send(b_bytes)
}

//服务端部分测试入口
func TestServer(t *testing.T)  {
	s := server.NewTcpServer()
	listen := new(TestListener)
	s.SetListener(listen)
	s.Start(":6001")
}


