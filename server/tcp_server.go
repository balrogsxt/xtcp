package server

import (
	"github.com/balrogsxt/xtcp/utils/logger"
	"net"
)

type TcpServer struct {
	listener TcpListener
}

func NewTcpServer() *TcpServer {
	return &TcpServer{

	}
}
func (c *TcpServer) Start(addr string) error {
	tcpAddr,err := net.ResolveTCPAddr("tcp4",addr)
	if err != nil {
		return err
	}
	listener,err := net.ListenTCP("tcp",tcpAddr)
	if err != nil {
		return err
	}
	_ = listener
	logger.Debug("Tcp客户端已启动,监听在 %s",addr)

	//监听tcp连接
	var fd = 0
	for{
		conn,err := listener.AcceptTCP()
		if err != nil {
			logger.Error("接受TCP客户端异常: %s",err.Error())
			continue
		}
		tcpConn := NewTcpConnection(fd,conn)
		tcpConn.SetListener(c.listener)
		//协程执行处理业务逻辑
		fd++
	}
	return nil
}
func (c *TcpServer) Stop() {

}

func (c *TcpServer) SetListener(listener TcpListener) {
	c.listener = listener
}
