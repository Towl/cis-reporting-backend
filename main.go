package main

import (
	"fmt"
)

// Version of the service
var Version = "development"

func init() {
	fmt.Printf("Version: %s\n", Version)
}

func main() {
	StartListening()
}
