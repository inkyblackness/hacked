package main

import (
	"fmt"
	"os"

	"github.com/inkyblackness/hacked/editor"
	"github.com/inkyblackness/hacked/ui/native"
)

func main() {
	var app editor.Application
	deferrer := make(chan func(), 100)

	err := native.Run(app.InitializeWindow, "InkyBlackness - HackEd", 30.0, deferrer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run application: %v\n", err)
	}
}
