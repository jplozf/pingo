package main

import "fmt"

// This variable will be overwritten at compile time
var Version = "0.0-dev"

func printVersion() {
	fmt.Printf("App Version: %s\n", Version)
}
