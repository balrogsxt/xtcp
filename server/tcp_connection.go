package server

import (
	"github.com/balrogsxt/xtcp/pack"
	"github.com/balrogsxt/xtcp/utils/logger"
	"io"
	"net"
	"sync"
	"time"
)

// TcpConnection tcp连接
type TcpConnection struct {
	// fd 连接标识
	fd int
	// conn tcp连接
	conn *net.TCPConn

	server *TcpServer
	//读写锁
	sync.RWMutex
	//tcp连接事件监听
	listener TcpListener
	//心跳检测定时器
	heartTimer *time.Ticker
	//上一次活动时间
	heartLastTime int64
	isClosed bool
}

func NewTcpConnection(fd int,s *TcpServer,conn *net.TCPConn) *TcpConnection {
	tcpConn := TcpConnection{
		fd: fd,
		conn: conn,
		server: s,
	}
	tcpConn.isClosed = false
	s.AddConnection(&tcpConn)
	go tcpConn.listenReader()
	go tcpConn.startHeart()
	//心跳检测机制
	return &tcpConn
}
// 监听客户端发送
func (c *TcpConnection) listenReader()  {
	if c.listener != nil {
		c.listener.OnOpen(c)
	}
	defer func(c *TcpConnection) {
		c.isClosed = true
		c.GetServer().RemoveConnection(c.fd)
		if c.listener != nil {
			c.listener.OnClose(c)
		}
	}(c)
	defer c.Close()
	//持续接收客户端发来的数据

	//心跳记录
	c.recordHeart()
	for {
		dataPack,err := pack.ParsePack(c.conn)
		if err != nil {
			if err == io.EOF {
				logger.Error("客户端连接已主动断开")
			}else{
				logger.Error("与客户端连接断开: %s",err.Error())
			}
			break
		}
		//心跳记录
		c.recordHeart()

		if c.listener != nil {
			//解包数据
			c.listener.OnMessage(c,dataPack)
		}
	}
}
func (c *TcpConnection) recordHeart(){
	//心跳记录
	c.heartLastTime = time.Now().Unix()
}
// startHeart 开启心跳检测
func (c *TcpConnection) startHeart(){
	// 检测时间、超时时间后续更改为可变配置
	c.heartTimer = time.NewTicker(time.Second * 20)
	for {
		<- c.heartTimer.C
		if time.Now().Unix() - c.heartLastTime > 60 {
			//心跳超时
			c.Close()
			logger.Error("关闭超时心跳客户端: %s",c.GetRemoteAddr().String())
		}
	}
}

//公开方法

//SetListener 设置监听事件
func (c *TcpConnection) setListener(listener TcpListener) {
	c.listener = listener
}

// GetFd 获取连接标识
func (c *TcpConnection) GetFd() int {
	return c.fd
}
// GetOriginTcpConnection 获取底层原始连接
func (c *TcpConnection) GetOriginTcpConnection() *net.TCPConn {
	return c.conn
}
// GetServer 获取TCP Server
func (c *TcpConnection) GetServer() *TcpServer {
	return c.server
}
// GetRemoteAddr 获取客户端连接地址
func (c *TcpConnection) GetRemoteAddr() net.Addr{
	return c.conn.RemoteAddr()
}

//Send 发送数据包给客户端
func (c *TcpConnection) Send(data []byte) error {

	//创建封包
	bytes,err := pack.Pack(pack.DataPack{
		Type: 1,
		Len: uint32(len(data)),
		Data: data,
	})
	if err != nil {
		return err
	}

	_,err = c.conn.Write(bytes)
	return err
}
//Close 主动断开该客户端连接
func (c *TcpConnection) Close() error {
	c.Lock()
	defer c.Unlock()
	//做一些其他操作
	//1.关闭连接
	c.conn.Close()
	//2.关闭心跳检测定时器
	if c.heartTimer != nil {
		c.heartTimer.Stop()
	}
	c.isClosed = true
	return nil
}
