package keywords

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"uk.ac.bris.cs/gameoflife/client"
)

func KeyRun() {

	keysEvents, err := keyboard.GetKeys(10)

	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press 'ESC' or 'q' to quit,请输入命令")
	var str string
	var CommandSlice []string
	commandIndex := 0

	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		if "A" <= string(event.Rune) && string(event.Rune) <= "z" {
			str = string(event.Rune)
			fmt.Println("请输入命令:", str)
			client.SendCommandToServer(str)
			if str == "q" {
				break
			}
		} else if event.Key == keyboard.KeyEnter {
			if str != "" {
				saveCommand(str, &CommandSlice)
				commandIndex = len(CommandSlice)
			}
			str = ""
			fmt.Printf("\r\n")
			fmt.Print("请输入命令:")
		}

		if event.Key == keyboard.KeyArrowUp {
			command := getCommand(&commandIndex, CommandSlice, true)
			fmt.Print("\033[2K\r", "请输入命令:", command)
		}
		if event.Key == keyboard.KeyArrowDown {
			command := getCommand(&commandIndex, CommandSlice, false)
			fmt.Print("\033[2K\r", "请输入命令:", command)

		}

		if event.Key == keyboard.KeyEsc {
			break
		}
	}
}

// 保存输入命令
func saveCommand(str string, commandSlice *[]string) {
	*commandSlice = append(*commandSlice, str)
	if len(*commandSlice) > 10 {
		*commandSlice = (*commandSlice)[1:]
	}
}

// 提取输入命令
func getCommand(commandIndex *int, commandSlice []string, true bool) (commamd string) {
	if true {
		if *commandIndex > 0 {
			*commandIndex--
		}
	} else {
		if *commandIndex < len(commandSlice) {
			*commandIndex++
		}
	}
	if *commandIndex < len(commandSlice) {
		return commandSlice[*commandIndex]
	} else {
		return ""
	}
}
