package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"io/ioutil"
	"math/rand"
	"net"
	"time"
)

// 获取节点下的所有子节点
func GetServerList(conn *zk.Conn) []string {
	paths, _, err := conn.Children("/go_zk_server")
	if err != nil {
		fmt.Println("get server children err :", err)
		return nil
	}
	return paths
}

// 连接zookeeper
func GetConnect() *zk.Conn {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*10)
	if err != nil {
		fmt.Println("conn zookeeper err : ", err)
		return nil
	}
	return conn
}

// 获取服务器地址
func getServerAddr() string {
	conn := GetConnect()
	defer conn.Close()
	addressList := GetServerList(conn)
	// 随机获取一个服务器地址
	len := len(addressList)
	if len == 0 {
		fmt.Println("server list is empty")
		return ""
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	host := addressList[r.Intn(len)]
	return host
}

//启动客户端
func StartClient()  {
	// 获取服务端的地址
	host := getServerAddr()
	fmt.Println("connected host : ", host)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		fmt.Println("resolve tcp address err:", err)
		return
	}
	//建立连接
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("dial tcp err:",err)
		return
	}
	defer conn.Close()
	//发送请求
	_, err = conn.Write([]byte("timestamp"))
	if err != nil {
		fmt.Println("client write err :", err)
		return
	}
	//获取响应
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("read response err:", err)
		//return
	}
	fmt.Println("resp:-------------->", string(resp))
}

func main() {
	for {
		StartClient()
		time.Sleep(time.Second * 2)
	}
}