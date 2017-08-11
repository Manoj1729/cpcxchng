package main

import (
	"fmt"

	"httpstub"
)

func main() {
	fmt.Println("vim-go")
	httpstub.Start()
}
