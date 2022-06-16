package main

import (
	"fmt"
	"proxier/pkg/logger"
)

func main() {
	f := logger.Dummy{}
	f.Printf("hello world!")
	fmt.Println("hello world")

}
