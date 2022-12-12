package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"runtime"
	"strconv"

	"uk.ac.bris.cs/gameoflife/gol"
	"uk.ac.bris.cs/gameoflife/sdl"
)

// main is the function called when starting Game of Life with 'go run .'
func MainFunc(keyword string) {
	runtime.LockOSThread()
	var params gol.Params

	flag.IntVar(
		&params.Threads,
		"t",
		8,
		"Specify the number of worker threads to use. Defaults to 8.")

	flag.IntVar(
		&params.ImageWidth,
		"w",
		512,
		"Specify the width of the image. Defaults to 512.")

	flag.IntVar(
		&params.ImageHeight,
		"h",
		512,
		"Specify the height of the image. Defaults to 512.")

	flag.IntVar(
		&params.Turns,
		"turns",
		10000000000,
		"Specify the number of turns to process. Defaults to 10000000000.")

	noVis := flag.Bool(
		"noVis",
		true,
		"Disables the SDL window, so there is no visualisation during the tests.")

	flag.Parse()

	fmt.Println("Threads:", params.Threads) // 线程数
	fmt.Println("Width:", params.ImageWidth)
	fmt.Println("Height:", params.ImageHeight)

	keyPresses := make(chan rune, 10)
	events := make(chan gol.Event, 1000)

	if len(keyword) > 0 {
		// 存入
		kw := []rune(keyword)
		keyPresses <- kw[0]
	}

	go gol.Run(params, events, keyPresses)
	if !(*noVis) {
		sdl.Run(params, events, keyPresses)
	} else {
		complete := false
		for !complete {
			event := <-events
			switch event.(type) {
			case gol.FinalTurnComplete:
				complete = true
			}
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------

//如果按s键，控制器将生成一个包含单板当前状态的PGM文件。
//如果按“q”，则关闭控制器客户端程序，不会导致GoL服务器出现错误。新的控制器应该能够接管与GoL引擎的交互。注意，您可以自由定义新控制器如何接管交互的性质。
//最有可能的状态将被重置。如果你设法继续上一个世界，这将被认为是一种扩展和容错的形式。
//如果按下k，则分布式系统的所有组件干净地关闭，系统输出最新状态的PGM图像。
//如果按下p，则暂停AWS节点上的处理，并让控制器打印正在处理的当前回合。如果再次按下p，则继续处理，并让控制器打印“持续中”。q和s在执行暂停时没有必要工作。

// ----------------------------------------------------------------------------------------------------------------------
type CSAService struct{}

func TestCSA() { //本地函数调用
	data := BytTest("1", "2", 3)
	fmt.Println("--------data:", data) //123
}

func BytTest(data1, data2 string, data3 int) string { //服务器定义函数
	return data1 + data2 + strconv.Itoa(data3)
} //调本地

func (p *CSAService) Run(request string, reply *string) error { //远程调用客户端调用服务器函数。返回对应关系
	fmt.Println("------接收到的数据：", request)
	switch request {
	case "s":
		fmt.Println("控制器将生成一个包含单板当前状态的PGM文件")
		*reply = "服务器接收到数据:" + request
		MainFunc("s")
		break
	case "q":
		fmt.Println("客户端程序退出！")
		*reply = "服务器接收到数据:" + request
		MainFunc("q")
		break
	case "k":
		fmt.Println("分布式系统的所有组件关闭，系统输出最新状态的PGM图像")
		*reply = "服务器接收到数据:" + request
		MainFunc("k")
		break
	case "p":
		fmt.Println("则暂停AWS节点上的处理，并让控制器打印正在处理的当前回合。如果再次按下p，则继续处理，并让控制器打印“持续中”")
		*reply = "服务器接收到数据:" + request
		MainFunc("p")
		break
	default:
		fmt.Println("按键请求错误！！！")
		*reply = "按键请求错误:" + request
	}
	return nil
}

// 获取细胞数据
func (p *CSAService) GetCellInfo(request string, reply *[][]byte) error { //获取细胞信息，worker传给服务器，worker在服务器里的运行
	fmt.Println("------接收到的数据：", request)
	*reply = gol.Gworld
	fmt.Println("------返回客户端的数据：", *reply)
	return nil
}

func (p *CSAService) ByteTest(request int, reply *string) error {
	*reply = "按键请求错误:" + strconv.Itoa(request)
	return nil
}

func main() {

	// 把我们的对象注册成一个rpc的 receiver
	// 其中rpc.Register函数调用会将对象类型中所有满足RPC规则的对象方法注册为RPC函数，
	// 所有注册的方法会放在“CSAService”服务空间之下
	rpc.RegisterName("CSAService", new(CSAService)) //网络连接new(CSAService）注册空间与服务

	// 然后我们建立一个唯一的TCP链接，  三次握手，持续连接
	listener, err := net.Listen("tcp", ":8666")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	// 通过rpc.ServeConn函数在该TCP链接上为对方提供RPC服务。
	// 没Accept一个请求，就创建一个goroutie进行处理
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		// 前面都是tcp的知识, 到这个RPC就接管了
		// 因此 你可以认为 rpc 帮我们封装消息到函数调用的这个逻辑,
		go rpc.ServeConn(conn)
	}
}

//----------------------------------------------------------------------------------------------------------------------
// 第一步骤：
// 1.sdl环境配置,第三发库的依赖；server --aws
// 2.本地编译好,打包命令：go build，执行文件放到server --aws
// 注：go build 前提：go env  GOOS = ? 例如GOOS=windows  打包文件包含windows执行的API
//  打包linux系统文件，先设置环境变量：go env -w GOOS=linux   go env -w GOOS=windows
// 3. AWS 连接;文件上传，文件执行：./XXXX   XXXX:需要启动的执行文件的名字，然后回车运行

// 第二步骤：本地客户端，远程服务器连接
// 1. 网络连接，RPC
// 2. 代码
