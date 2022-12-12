package client

import (
	"fmt"
	"log"
	"net/rpc"
)

var (
	Gclient  *rpc.Client //全局变量保存，客户端连服务器连一次就好，链接数据保存
	HostAdds = "local.golang.ltd"
)

func init() { //
	client, err := rpc.Dial("tcp", HostAdds+":8666") //长连接，底层
	if err != nil {
		log.Fatal("dialing:", err)
	}
	Gclient = client
} //连接服务器

// 发送命令
func SendCommandToServer(keyword string) {
	// 首先是通过rpc.Dial拨号RPC服务, 建立连接
	// 然后通过client.Call调用具体的RPC方法
	// 在调用client.Call时:
	// 第一个参数是用点号链接的RPC服务名字和方法名字，
	// 第二个参数是 请求参数
	// 第三个是请求响应, 必须是一个指针, 有底层rpc服务帮你赋值
	var reply, senddata string
	senddata = keyword
	err := Gclient.Call("CSAService.Run", senddata, &reply) //发送与接收
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)

	SendByteTest()
}

func SendByteTest() {
	var reply string
	var senddata int
	senddata = 10000
	err := Gclient.Call("CSAService.ByteTest", senddata, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
