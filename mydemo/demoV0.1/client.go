package main

import (
	"fmt"
	"net"
	"time"
)

type Client struct {
	Ip string
	Port int
}

func NewClient(ip string,port int) *Client{
	return &Client{Ip: ip,Port: port}
}

func (c *Client)Start()  {
	//1.Dial
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Ip, c.Port))
	if err != nil {
		fmt.Println("net.Dial error:",err)
		return
	}
	for  {
		_, err := conn.Write([]byte("你好啊"))
		if err != nil {
			fmt.Println("conn.Write error:",err)
			return
		}
		buf := make([]byte,1024)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read buf error:",err)
			return
		}
		fmt.Printf("server call back %s,%d",buf,cnt)
		time.Sleep(time.Second * 1)
	}
}

func main (){
	client := NewClient("127.0.0.1",8888)
	client.Start()
}