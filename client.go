package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int		//当前client 先择
}

var serverIp string
var serverPort int

func init() {
	// 使用命令行工具包 flag 配置为 ./client  -ip 127.0.0.1 -port 8888
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip地址默认(127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器ip地址默认(8888)")
}
func main() {
	flag.Parse()
	//启动Client
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>链接服务器失败")
		return
	}
	fmt.Println(">>>>>>>>链接服务器成功")
	//阻塞一下
	client.Run()
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag: 999,
	}
	// 链接Server端
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("Client net.Dial error:", err)
		return nil
	}

	//配置链接
	client.conn = conn
	return client
}

// Menu 客户端菜单功能
func (c *Client) Menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag) //阻塞监听用户输入

	if c.flag >= 0 && c.flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法数字<<<<<")
		return false
	}

}

// Run 客户端命令窗口不断提示
func (c *Client) Run() {
	for c.flag != 0 {
		for c.Menu() != true {
		}
		//根据不同模式处理不同业务
		switch c.flag {
		case 1:
			//公聊模式
			fmt.Println("进入公聊模式")
			break
		case 2:
			//私聊模式
			fmt.Println("进入私聊模式")
			break
		case 3:
			//更新用户名
			fmt.Println("使用 \"rename|姓名\" 进行更换名称")
			break
		}
	}
}

