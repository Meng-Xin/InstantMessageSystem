package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string //绑定主机地址
	Port int    //绑定端口号

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// NewServer 初始化server实例化对象
func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// Handler 处理当前连接的业务
func (s *Server) Handler(conn net.Conn) {
	//链接当前的业务

	//1.用户建立conn之后，把当前用户加入OnlineMap中
	user := NewUser(conn)
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//2.通过BroadCast()广播功能通知所有用户
	s.BroadCast(user, "已上线")

	//阻塞
	select {

	}
}

// ListenMessage 监听当前Message广播channel的goroutine，一旦有消息就发送给全部用户
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		//将msg发送给全部的User
		s.mapLock.RLock()
		for _, user := range s.OnlineMap {
			user.Message <- msg
		}
		s.mapLock.RUnlock()
	}
}

// BroadCast 广播方法，通知当前用户
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

// Start 启动服务器的接口
func (s *Server) Start() {
	fmt.Println("-------->开启聊天室<----------")
	//1.socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}

	//2.close listen socket
	defer listen.Close()

	//3.启动监听 BroadMessage 广播消息队列，给所有User反馈消息
	go s.ListenMessage()

	//3.不断处理socket链接传递的信息
	for {
		//4.accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.accept err:", err)
			continue
		}
		//5.do handler
		go s.Handler(conn)
		//TODO 后期加一个通知conn断开链接结束监听的Context
	}

}
