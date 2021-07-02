package server

import (
	"context"
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
	//退出标识管道
	ctx context.Context
	cancel context.CancelFunc

}

func NewTcpConnection(fd int,s *TcpServer,conn *net.TCPConn) *TcpConnection {
	tcpConn := TcpConnection{
		fd: fd,
		conn: conn,
		server: s,
	}
	ctx,cancel := context.WithCancel(context.Background())
	tcpConn.ctx = ctx
	tcpConn.cancel = cancel

	tcpConn.isClosed = false
	s.AddConnection(&tcpConn)
	go tcpConn.listenReader()
	if s.config.HeartInterval > 0 && s.config.HeartTimeout >= s.config.HeartInterval{
		go tcpConn.startHeart()
	}
	//心跳检测机制
	return &tcpConn
}
// 监听客户端发送
func (c *TcpConnection) listenReader()  {
	if c.listener != nil {
		c.listener.OnOpen(c)
	}
	defer c.Close()
	//心跳记录
	c.recordHeart()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			dataPack,err := pack.TcpConnUnPack(c.server.config.DataPack,c.conn)
			if err != nil {
				if err == io.EOF {
					logger.Error("客户端连接已主动断开")
				}else{
					logger.Error("与客户端连接断开: %s",err.Error())
				}
				c.cancel()
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
}
func (c *TcpConnection) recordHeart(){
	if c.server.config.HeartInterval > 0 && c.server.config.HeartTimeout >= c.server.config.HeartInterval{
		//心跳记录
		c.heartLastTime = time.Now().Unix()
	}
}
// startHeart 开启心跳检测
func (c *TcpConnection) startHeart(){
	// 检测时间、超时时间后续更改为可变配置
	c.heartTimer = time.NewTicker(time.Second * time.Duration(c.server.config.HeartInterval))
	for {
		<- c.heartTimer.C
		if time.Now().Unix() - c.heartLastTime > int64(c.server.config.HeartTimeout) {
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

//SendPack 发送封包数据包给客户端
func (c *TcpConnection) SendPack(typeId uint32,data []byte) error {
	//创建封包
	bytes,err := pack.Pack(pack.DataPack{
		Type: typeId,
		Len: uint32(len(data)),
		Data: data,
	})
	if err != nil {
		return err
	}
	return c.SendBytes(bytes)
}
//SendBytes 发送字节流数据给客户端
func (c *TcpConnection) SendBytes(data []byte) error {
	_,err := c.conn.Write(data)
	return err
}
//Close 主动断开该客户端连接
func (c *TcpConnection) Close() error {
	c.Lock()
	defer c.Unlock()
	//做一些其他操作
	if c.listener != nil {
		c.listener.OnClose(c)
	}
	//1.关闭连接
	c.conn.Close()
	//2.关闭心跳检测定时器
	if c.heartTimer != nil {
		c.heartTimer.Stop()
		c.heartTimer = nil
	}
	c.cancel()
	c.isClosed = true
	//删除连接管理
	c.GetServer().RemoveConnection(c.fd)

	return nil
}
