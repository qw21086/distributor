package main

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
