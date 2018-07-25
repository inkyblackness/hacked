package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	_ "image/gif"
	_ "image/png"

	"github.com/inkyblackness/hacked/crash"
	"github.com/inkyblackness/hacked/editor"
	"github.com/inkyblackness/hacked/ui/native"
)

var version string

func main() {
	scale := flag.Float64("scale", 1.0, "factor for scaling the UI (0.5 .. 10.0). 1080p displays should use default. 4K most likely 2.0.")
	fontFile := flag.String("fontfile", "", "Path to font file (.TTF) to use instead of the default font. Useful for HiDPI displays.")
	fontSize := flag.Float64("fontsize", 0.0, "Size of the font to use. If not specified, a default height will be used.")
	flag.Parse()
	var app editor.Application
	app.FontFile = *fontFile
	app.FontSize = float32(*fontSize)
	app.GuiScale = float32(*scale)
	if len(version) > 0 {
		app.Version = version
	} else {
		app.Version = fmt.Sprintf("(manual build %v)", time.Now().Format("2006-01-02"))
	}
	deferrer := make(chan func(), 100)

	versionInfo := "InkyBlackness - HackEd - " + app.Version
	defer crash.Handler(versionInfo)

	err := native.Run(app.InitializeWindow, versionInfo, 30.0, deferrer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run application: %v\n", err)
	}
}
