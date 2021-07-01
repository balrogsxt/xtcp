package main

import (
	"bufio"
	"fmt"
	"github.com/balrogsxt/xtcp/server"
	"io"
	"log"
	"net"
	"testing"
	"time"
)
type TestListener struct {}

func (c *TestListener) OnOpen(conn *server.TcpConnection)  {
	fmt.Println(conn.GetRemoteAddr(),"客户端已连接",conn.GetFd())
}
func (c *TestListener) OnClose(conn *server.TcpConnection)  {
	fmt.Println(conn.GetRemoteAddr(),"客户端已断开",conn.GetFd())
}
func (c *TestListener) OnMessage(conn *server.TcpConnection,data []byte)  {
	fmt.Println(conn.GetRemoteAddr(),"客户端发送消息:",conn.GetFd(),string(data))
	conn.Send([]byte("收到数据:"+string(data)))
}

//服务端部分测试入口
func TestServer(t *testing.T)  {
	s := server.NewTcpServer()
	listen := new(TestListener)
	s.SetListener(listen)
	s.Start(":6001")
}

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
	conn.Write([]byte("字符串测试发送"))

	//go func() {
		reader := bufio.NewReader(conn)
		for {
			buffer := make([]byte,1024)
			size,err := reader.Read(buffer)
			if err == io.EOF {
				log.Printf("连接已断开")
				break
			}
			if err != nil {
				log.Printf("接收失败: %s",err.Error())
				continue
			}
			bufferBytes := buffer[0:size]
			fmt.Println("收到服务端发的数据:" + string(bufferBytes))
			conn.Write([]byte("ping -> "+time.Now().Format("2006-01-02 15:04:05")))
			time.Sleep(time.Millisecond * 500)
		}
	//}()
}
