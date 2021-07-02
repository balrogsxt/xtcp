package server

import (
	"errors"
	"github.com/balrogsxt/xtcp/utils/logger"
	"net"
	"sync"
)

type TcpServer struct {
	//基础事件监听
	listener TcpListener
	connections sync.Map
	config TcpServerConfig
}
type TcpServerConfig struct {
	HeartInterval int //心跳检测间隔秒数,小于等于0时不适用
	HeartTimeout int //超时时间,不可低于HeartInterval
	DataPack bool //是否进行封包,拆包处理
}

func NewTcpServer(config TcpServerConfig) *TcpServer {
	return &TcpServer{
		config: config,
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
	if c.config.HeartInterval > 0 && c.config.HeartTimeout >= c.config.HeartInterval{
		logger.Debug("[心跳检测] 已启用心跳检测机制")
		logger.Debug("[心跳检测] 检测间隔 %d 秒",c.config.HeartInterval)
		logger.Debug("[心跳检测] 超时时间 %d 秒",c.config.HeartTimeout)
	}else{
		if c.config.HeartInterval > c.config.HeartTimeout {
			return errors.New("心跳检测时间不可大于超时时间")
		}
	}

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
		connection.SendBytes(data)
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
