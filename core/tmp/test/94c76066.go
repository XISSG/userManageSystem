package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	sum := 0
	for _, arg := range args {
		num, _ := strconv.Atoi(arg)
		sum += num
			}
	fmt.Println(sum)
}
