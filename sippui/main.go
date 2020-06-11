// Copyright Raul Vera 2015-2016

// A simple UI for the SIPP library that uses QML.

package main

import (
	"fmt"
	"os"

	"gopkg.in/qml.v1"

	"github.com/Causticity/sipp/stree"
)

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

var appComponent qml.Object

var app *qml.Window

func run() error {
	engine := qml.NewEngine()

	appComponent, err := engine.LoadFile("sippui.qml")
	if err != nil {
		return err
	}

	err = stree.InitTreeComponents(engine)
	if err != nil {
		return err
	}

	app = appComponent.CreateWindow(nil)
	app.On("gotFile", stree.NewSippRootNode)

	app.Call("getFile")

	// This Show() is necessary, or the app hangs when quitting. Weird.
	app.Show()
	app.Wait()

	return nil
}
