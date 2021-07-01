package main

import (
	"bufio"
	"fmt"
	"github.com/balrogsxt/xtcp/pack"
	"github.com/balrogsxt/xtcp/server"
	"github.com/balrogsxt/xtcp/utils/logger"
	"io"
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
	fmt.Println(conn.GetRemoteAddr(),"客户端发的消息:",conn.GetFd(),string(data))
	conn.Send([]byte("收到数据:"+string(data)))
	//fmt.Println(fmt.Sprintf("%#v \n",data),"长度:",len(data))
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


	go func() {
		reader := bufio.NewReader(conn)
		for {

			buffer := make([]byte,pack.GetHeadLen())
			_,err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Error("读取Buffer失败: %s",err.Error())
				continue
			}
			//解析出包头后,根据包头数据结果拿到数据包
			dataPack,err := pack.UnPack(buffer)
			if err != nil {
				//拆包失败...
				logger.Error("拆包失败: %s",err.Error())
				continue
			}
			//拆包完成,再次读取数据包数据
			if dataPack.Len > 0 {
				//再次读取包头
				buffer := make([]byte,dataPack.Len)
				_,err := reader.Read(buffer)
				if err == io.EOF {
					break
				}
				if err != nil {
					logger.Error("读取Buffer失败: %s",err.Error())
					continue
				}
				fmt.Println(fmt.Sprintf("收到内容: %s",buffer))
			}
		}
	}()
	for {
		d := []byte(time.Now().Format("2006-01-02 15:04:05"))
		datapack,err := pack.Pack(pack.DataPack{
			Type: 1,
			Len: uint32(len(d)),
			Data: d,
		})
		if err != nil {
			t.Fatalf("封包失败: %s",err.Error())
		}
		//发送封包数据
		conn.Write(datapack)
		time.Sleep(time.Second)
	}

}
