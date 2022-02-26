package main

import "net"

type User struct {
	Name    string
	Addr    string
	Message chan string
	conn    net.Conn
}

// NewUser 初始化User实例对象
func NewUser(conn net.Conn) *User {
	//1.获取对应链接的信息
	userAddr := conn.RemoteAddr().String()
	//2.初始化user结构体内容
	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Message: make(chan string),
		conn:    conn,
	}
	//3.启动监听当前User 消息的方法--->ListenMessage协程
	go user.ListenMessage()
	return user
}

// ListenMessage 监听 当前User Channel的方法，一旦有消息就发送给对端
func (u *User) ListenMessage() {
	for  {
		msg := <- u.Message
		u.conn.Write([]byte(msg))
	}
}
