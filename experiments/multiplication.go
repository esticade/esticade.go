package main

import (
	"fmt"
	"github.com/esticade/esticade.go"
)

func main() {
	fmt.Println("Connecting..")
	service := esticade.NewService("Multiplication Service")
	fmt.Printf("Connected to %#v\n", service)
}
