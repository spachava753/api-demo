package main

import (
	"api-demo/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Printf("err: %s\n", err)
		os.Exit(1)
	}
}
