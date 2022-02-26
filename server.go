package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip     string    //绑定主机地址
	Port   int       //绑定端口号
	IsLive chan bool //超时强制T除用户
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
		IsLive:    make(chan bool),
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// Handler 处理当前连接的业务
func (s *Server) Handler(conn net.Conn) {

	//1.用户建立conn之后，把当前用户加入OnlineMap中
	user := NewUser(conn, s)
	user.Online()

	//3. 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("=》Server《= conn.Read() error:", err)
				return
			}

			//提取哦用户的消息去除(\n)
			msg := string(buf[:n-1])

			//将得到的消息进行广播
			user.DoMessage(msg)

			//用户的任意消息，代表当前用户是一个活跃的对象
			s.IsLive <- true
		}
	}()

	//阻塞
	for {
		select {
		case <-s.IsLive:
			// 在10s内，只要该用户操作就会一直写入数据到IsLive中，持续保活。
			// 不用做任何操作，只是为了重置定时器
		case <-time.After(time.Second * 10):
			// 强制下线通知（私聊）
			user.SendMsg("因长时间未操作，您被强制下线。。。")

			//下线广播（世界广博）
			user.Offline()

			// Server端释放当前用户的Conn链接。
			user.conn.Close()

			//退出当前的Handler
			return //也可以执行 runtime.Goexit() 释放协程资源
		}
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
