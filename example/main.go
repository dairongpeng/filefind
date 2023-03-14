package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("您好，这里我使用了中文;")

	str := "这里不, 确定能不能被, 工具发现呢。"
	if len(strings.Split(str, ",")) == 3 {
		fmt.Println("This is as expected")
	}

	fmt.Println("Hello filectl.")
}
