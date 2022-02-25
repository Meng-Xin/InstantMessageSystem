package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string			//绑定主机地址
	Port int		//绑定端口号
}

// NewServer 初始化server实例化对象
func  NewServer (ip string,port int) *Server{
	return &Server{Ip: ip,Port: port}
}

// Handler 处理当前连接的业务
func (s *Server)Handler(conn net.Conn)  {
	//链接当前的业务
}
// Start 启动服务器的接口
func (s *Server) Start (){
	fmt.Println("-------->开启聊天室<----------")
	//1.socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.listen err:",err)
		return
	}
	//2.close listen socket
	defer listen.Close()
	//3.不断处理socket链接传递的信息
	for  {
		//4.accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.accept err:",err)
			continue
		}
		//5.do handler
		go s.Handler(conn)
	}

}