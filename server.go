package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"net"
	"time"
)

//创建连接
func GetConn() *zk.Conn {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*10)
	if err != nil {
		fmt.Println("conn zookeeper err:", err)
		return nil
	}
	return conn
}

//注册服务
func registerServer(conn *zk.Conn, host string)  {
	// 创建临时节点
	path, err := conn.Create("/go_zk_server/"+host, []byte("new server"), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Println("register server err :", err)
		return
	}
	fmt.Println("register new server success, path:", path)
}

//获取所有服务
func getServerList(conn *zk.Conn) []string {
	list, _, err := conn.Children("/go_zk_server")
	if err != nil {
		fmt.Println("get server list err :", err)
		return nil
	}
	return list
}

//处理服务
func handleServer(conn net.Conn)  {
	defer conn.Close()
	nowTime := time.Now().String()
	conn.Write([]byte(nowTime))
}

// 开始服务
func startServer(addr string)  {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		fmt.Println("return tcp address err :", err)
		return
	}
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("tcp listen err :", err)
		return
	}

	// 连接zookeeper
	conn := GetConn()
	defer conn.Close()

	//注册节点服务
	registerServer(conn, addr)

	//等待服务客户端
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept err : ", err)
			continue
		}
		go handleServer(conn)
	}
	fmt.Println("server over")
}

func main() {
	go startServer("127.0.0.1:8866")
	go startServer("127.0.0.1:8860")
	go startServer("127.0.0.1:8868")
	ch := make(chan struct{}, 1)
	<- ch
}