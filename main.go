package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/inkyblackness/hacked/editor"
	"github.com/inkyblackness/hacked/ui/native"
)

var version string

func main() {
	scale := flag.Float64("scale", 1.0, "factor for scaling the UI (0.5 .. 10.0). 1080p displays should use default. 4K most likely 2.0.")
	flag.Parse()
	var app editor.Application
	app.GuiScale = float32(*scale)
	if len(version) > 0 {
		app.Version = version
	} else {
		app.Version = fmt.Sprintf("(manual build %v)", time.Now().Format("2006-01-02"))
	}
	deferrer := make(chan func(), 100)

	err := native.Run(app.InitializeWindow, "InkyBlackness - HackEd - "+app.Version, 30.0, deferrer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run application: %v\n", err)
	}
}
