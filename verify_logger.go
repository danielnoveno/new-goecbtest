package main

import (
	"fmt"
	"go-ecb/pkg/logging"
)

func main() {
	fmt.Println("Testing logger before Init()...")
	// This should NOT panic now
	logging.Logger().Infof("This is a test log before Init()")
	fmt.Println("Success: No panic occurred.")
}
