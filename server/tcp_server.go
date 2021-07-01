package server

import (
	"github.com/balrogsxt/xtcp/utils/logger"
	"net"
	"sync"
)

type TcpServer struct {
	//基础事件监听
	listener TcpListener
	connections sync.Map
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
	var fd = 1
	for{
		conn,err := listener.AcceptTCP()
		if err != nil {
			logger.Error("接受TCP客户端异常: %s",err.Error())
			continue
		}
		tcpConn := NewTcpConnection(fd,c,conn)
		tcpConn.setListener(c.listener)
		//协程执行处理业务逻辑
		fd++
	}
}
func (c *TcpServer) Stop() {

}

func (c *TcpServer) SetListener(listener TcpListener) {
	c.listener = listener
}
// AddConnection 添加一个TCP连接
func (c *TcpServer) AddConnection(connection *TcpConnection){
	c.connections.Store(connection.GetFd(),connection)
}
//RemoveConnection 移除一个Tcp连接[并不是关闭一个连接]
func (c *TcpServer) RemoveConnection(fd int) {
	c.connections.Delete(fd)
}
// GetConnection 通过ID获取一个TCP连接
func (c *TcpServer) GetConnection(fd int) *TcpConnection {
	obj,has := c.connections.Load(fd)
	if !has {
		return nil
	}
	if tc,has := obj.(*TcpConnection);has{
		return tc
	}
	return nil
}
//Broadcast 发送广播消息
func (c *TcpServer) Broadcast(data []byte) {
	c.GetAllConnection(func(connection *TcpConnection) {
		connection.Send(data)
	})
}
//GetAllConnection 回调获取所有连接
func (c *TcpServer) GetAllConnection(callback func(connection *TcpConnection)) {
	c.connections.Range(func(key, value interface{}) bool {
		if v,h := value.(*TcpConnection);h {
			callback(v)
		}
		return true
	})
}
