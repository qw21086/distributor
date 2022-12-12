package main

import (
	"uk.ac.bris.cs/gameoflife/keywords"
)

type CSAService struct{}

//----------------------------------------------------------------------------------------------------------------------

//如果按s键，控制器将生成一个包含单板当前状态的PGM文件。
//如果按“q”，则关闭控制器客户端程序，不会导致GoL服务器出现错误。新的控制器应该能够接管与GoL引擎的交互。注意，您可以自由定义新控制器如何接管交互的性质。
//最有可能的状态将被重置。如果你设法继续上一个世界，这将被认为是一种扩展和容错的形式。
//如果按下k，则分布式系统的所有组件干净地关闭，系统输出最新状态的PGM图像。
//如果按下p，则暂停AWS节点上的处理，并让控制器打印正在处理的当前回合。如果再次按下p，则继续处理，并让控制器打印“持续中”。q和s在执行暂停时没有必要工作。

//----------------------------------------------------------------------------------------------------------------------

func init() {
	keywords.KeyRun()
}

func main() {
}
