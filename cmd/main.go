package main

import (
	"fmt"

	"netcat/cmd/server"
)

func main() {
	if err := server.Start(); err != nil {
		fmt.Println(err)
		return
	}
}
