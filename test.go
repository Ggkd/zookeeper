package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)
//创建连接
func getConn(addressList []string) *zk.Conn {
	conn, _, err := zk.Connect(addressList, time.Second)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return conn
}

// 测试连接
func test()  {
	address := []string{"127.0.0.1:2181"}
	conn := getConn(address)
	fmt.Println("conn success")
	defer conn.Close()
}

func create(conn *zk.Conn)  {
	var flag int32 = 0
	//flag有4种取值：
	//0:永久，除非手动删除
	//zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除
	//zk.FlagSequence  = 2:会自动在节点后面添加序号
	//3:Ephemeral和Sequence，即，短暂且自动添加序号

	// 创建节点
	path, err := conn.Create("/go_zk_server", []byte("test"), flag, zk.WorldACL(zk.PermAll))	// zk.WorldACL(zk.PermAll)控制访问权限模式
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("create node success, path: ", path)
}

func set(conn *zk.Conn)  {
	// 修改节点
	_, err := conn.Set("/go_zk_server", []byte("set_value"), -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("set node success")
	value, _, _ := conn.Get("/go_zk_server")
	fmt.Println("value: ", string(value))
}

func delete_node(conn *zk.Conn)  {
	// 删除节点
	_, err := conn.Set("/go_zk_server", []byte("set_value"), -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("set node success")
	value, _, _ := conn.Get("/go_zk_server")
	fmt.Println("value: ", string(value))
}

func get_all(conn *zk.Conn)  {
	// 获取所有节点
	children, stat,  err := conn.Children("/")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("children : %#v,   stat : %#v \n", children, stat)
	//e := <- ch
	//fmt.Printf("event : %#v\n ", e)
}

func get_one(conn *zk.Conn)  {
	value, _, _ := conn.Get("/go_zk_server")
	fmt.Println("value: ", string(value))
}

func main() {
	address := []string{"127.0.0.1:2181"}
	conn := getConn(address)
	defer conn.Close()
	//test()
	//create(conn)
	//set(conn)
	get_all(conn)
	//get_one(conn)
}