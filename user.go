package main

import (
	"net"
	"strings"
)

type User struct {
	Name    string
	Addr    string
	Message chan string
	conn    net.Conn
	server *Server
}

// NewUser 初始化User实例对象
func NewUser(conn net.Conn,server *Server) *User {
	//1.获取对应链接的信息
	userAddr := conn.RemoteAddr().String()
	//2.初始化user结构体内容
	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		Message: make(chan string),
		conn:    conn,
		server: server,
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


// Online 用户上线业务
func (u *User)Online(){
	// 用户上线，把用户添加到OnlineMap 中
	u.server.mapLock.Lock()
	defer u.server.mapLock.Unlock()
	u.server.OnlineMap[u.Name] = u

	// 发送广播通知所有用户,当前用户上线
	u.server.BroadCast(u,"已上线")
}
// Offline 用户下线业务
func (u *User)Offline(){
	// 用户上线，把用户添加到OnlineMap 中
	u.server.mapLock.Lock()
	defer u.server.mapLock.Unlock()
	//将当前下线用户从map中删除
	delete(u.server.OnlineMap,u.Name)
	//关闭当前管道，释放资源
	//close(u.Message)
	// 发送广播通知所有用户,当前用户下线
	u.server.BroadCast(u,"已下线")
}
// DoMessage 处理用户信息业务
func (u *User)DoMessage(msg string){
	if msg == "who" {
		//查询当前在线用户有哪些
		u.server.mapLock.RLock()
		defer u.server.mapLock.RUnlock()
		for _,user := range u.server.OnlineMap{
			onlineMsg := "[" + user.Addr + "]" + user.Name+"在线..."
			u.SendMsg(onlineMsg)
		}
	}else if len(msg)>7 && msg[:7] == "rename|" {
		u.UpdateName(msg)
	} else if len(msg)>4 && msg[:3] == "to|"{
		u.PrivateChat(u.conn,msg)
	}else{
		u.server.BroadCast(u,msg)
	}



}

// SendMsg 发送消息给conn
func (u *User) SendMsg (msg string){
	u.conn.Write([]byte(msg))
}

// PrivateChat 用户模块实现私聊功能
func (u *User)PrivateChat (conn net.Conn,msg string){
	// 私聊消息格式： to|张三|消息内容

	//1。获取对方用户名
	remoteName := strings.Split(msg, "|")[1]
	if remoteName == "" {
		u.SendMsg("消息格式不真情古额，请使用\"to|张三|你好啊\"格式")
		return
	}

	//2. 根据用户名，得到对应用户User对象
	remoteUser,ok := u.server.OnlineMap[remoteName]
	if !ok {
		u.SendMsg("当前用户名不存在！")
		return
	}

	//3. 为对应User发送消息
	content := strings.Split(msg,"|")[2]
	if content == "" {
		u.SendMsg("没有消息内容，请重新发送")
		return
	}
	remoteUser.SendMsg(content)
}

// UpdateName 更新用户名
func (u *User) UpdateName (msg string){
	//消息格式：rename|张三
	newName := strings.Split(msg,"|")[1]
	//判断name是否存在
	_,ok := u.server.OnlineMap[newName]
	if ok {
		u.SendMsg("当前用户名已经被使用\n")
	}else{
		u.server.mapLock.Lock()
		defer u.server.mapLock.Unlock()
		u.server.OnlineMap[newName] = u

		//修改当前名字
		u.Name = newName
		u.SendMsg("您已经成功更新用户名："+u.Name+"\n")
	}
}