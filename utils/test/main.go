package main

import (
	"github.com/HEBNUOJ/utils"
)

func main() {
	// 输出到文件
	logger := utils.Log("test.log", 1)
	logger.Println("自定义log")
	// 输出到屏幕
	logger = utils.Log("", 3)
	logger.Println("自定义log")
}
