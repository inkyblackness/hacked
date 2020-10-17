package main

import (
	"flag"
	"fmt"
	_ "image/gif"
	_ "image/png"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	_ "golang.org/x/image/bmp"

	"github.com/inkyblackness/hacked/crash"
	"github.com/inkyblackness/hacked/editor"
	"github.com/inkyblackness/hacked/ui/native"
)

var version string

func main() {
	scale := flag.Float64("scale", 1.0, "factor for scaling the UI (0.5 .. 10.0). 1080p displays should use default. 4K most likely 2.0.")
	fontFile := flag.String("fontfile", "", "Path to font file (.TTF) to use instead of the default font. Useful for HiDPI displays.")
	fontSize := flag.Float64("fontsize", 0.0, "Size of the font to use. If not specified, a default height will be used.")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
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

	configDir, err := configDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to determine config dir: %v", err)
		os.Exit(-1)
	}
	app.ConfigDir = configDir

	versionInfo := "InkyBlackness - HackEd - " + app.Version
	defer crash.Handler(versionInfo)

	profileFin, err := initProfiling(*cpuprofile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start CPU profiling: %v\n", err)
	}
	defer profileFin()

	err = native.Run(app.InitializeWindow, versionInfo, 30.0, deferrer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run application: %v\n", err)
	}
}

func initProfiling(filename string) (func(), error) {
	if filename != "" {
		f, err := os.Create(filename)
		if err != nil {
			return func() {}, err
		}
		err = pprof.StartCPUProfile(f)
		return func() { pprof.StopCPUProfile() }, err
	}
	return func() {}, nil
}

func configDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(base, "InkyBlackness", "HackEd")
	err = os.MkdirAll(fullPath, os.ModeDir|0750)
	if err != nil {
		return "", err
	}
	return fullPath, nil
}
