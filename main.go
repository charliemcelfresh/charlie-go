package main

import (
	"fmt"

	"github.com/charliemcelfresh/go-by-charlie/cmd"
)

func init() {
	fmt.Println("Running main config")
}

func main() {
	cmd.Execute()
}
