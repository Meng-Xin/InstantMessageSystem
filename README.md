​									**Golang 使用原生net实现聊天室**



# Server 端

## 1.server.go

### 	1.1 基础模块 

> server属性: `Ip`服务端ip,`Port`服务端端口,`Online`在线用户map,`Message`广播channel 

```go
type Server struct {
	Ip   string //绑定主机地址
	Port int    //绑定端口号

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}
```



## 2.user.go

# Client 端

