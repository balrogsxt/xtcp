package server

import (
	"github.com/balrogsxt/xtcp/utils/logger"
	"io"
	"net"
	"sync"
)

// TcpConnection tcp连接
type TcpConnection struct {
	// fd 连接标识
	fd int
	// conn tcp连接
	conn *net.TCPConn
	//读写锁
	sync.RWMutex
	//tcp连接事件监听
	listener TcpListener
}

func NewTcpConnection(fd int,conn *net.TCPConn) *TcpConnection {
	tcpConn := TcpConnection{
		fd: fd,
		conn: conn,
	}
	go tcpConn.call()
	return &tcpConn
}
func (c *TcpConnection) call()  {
	if c.listener != nil {
		c.listener.OnOpen(c)
	}
	defer func(c *TcpConnection) {
		if c.listener != nil {
			c.listener.OnClose(c)
		}
	}(c)
	defer c.Close()

	//持续接收客户端发来的数据
	for {
		buffer := make([]byte,1024)
		size,err := c.conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("读取Buffer失败: %s",err.Error())
			continue
		}
		if c.listener != nil {
			c.listener.OnMessage(c,buffer[0:size])
		}
	}
}



//公开方法

//SetListener 设置监听事件
func (c *TcpConnection) SetListener(listener TcpListener) {
	c.listener = listener
}

// GetFd 获取连接标识
func (c *TcpConnection) GetFd() int {
	return c.fd
}
// GetConnection 获取底层连接
func (c *TcpConnection) GetConnection() *net.TCPConn {
	return c.conn
}
// GetRemoteAddr 获取客户端连接地址
func (c *TcpConnection) GetRemoteAddr() net.Addr{
	return c.conn.RemoteAddr()
}

//Send 发送数据包给客户端
func (c *TcpConnection) Send(data []byte) error {
	_,err := c.conn.Write(data)
	return err
}
//Close 主动断开该客户端连接
func (c *TcpConnection) Close() error {
	return c.conn.Close()
}
