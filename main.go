package main

import (
	"fmt"
	"os"

	"github.com/inkyblackness/hacked/ui"
	"github.com/inkyblackness/hacked/ui/native"
)

func main() {
	var app ui.Application
	deferrer := make(chan func(), 100)

	err := native.Run(app.InitializeWindow, "Test Window", 30.0, deferrer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run application: %v\n", err)
	}
}
