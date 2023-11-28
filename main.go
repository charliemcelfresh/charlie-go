package main

import (
	"fmt"

	"github.com/charliemcelfresh/charlie-go/cmd"
)

func init() {
	fmt.Println("Running main config")
}

func main() {
	cmd.Execute()
}
