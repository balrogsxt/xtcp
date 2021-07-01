package server

type ITcpServer interface{
	// Start 服务启动
	Start(addr string) error
	// Stop 服务停止
	Stop()
	// SetListener 设置监听事件接收
	SetListener(listener TcpListener)
}

type TcpListener interface {
	//OnOpen 客户端连接事件
	OnOpen(conn *TcpConnection)
	//OnClose 客户端关闭事件
	OnClose(conn *TcpConnection)
	//OnMessage 客户端发送消息事件
	OnMessage(conn *TcpConnection,data []byte)
}

