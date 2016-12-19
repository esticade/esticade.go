package main

import (
	"fmt"
	"github.com/esticade/esticade.go"
)

func main() {
	fmt.Println("Connecting..")
	service, err := esticade.NewService("Multiplication Service")
	if err != nil {
		fmt.Printf("Connection failed (%s)\n", err.Error())
		return
	}

	fmt.Printf("Connected to %#v\n", service)
	defer service.Shutdown()
}
