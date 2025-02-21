package main

import (
	"fmt"
	"os"
)

func main() {
	filePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(filePath)
}