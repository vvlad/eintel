package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/glib"
	"log"
  "os"
)

const appID = "com.verestiuc.eintel"

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatal("Could not create application.", err)
	}

	application.Connect("activate", func() { onActivate(application) })
	os.Exit(application.Run(os.Args))
}


func onActivate(application *gtk.Application) {
	appWindow, err := gtk.ApplicationWindowNew(application)
	if err != nil {
		log.Fatal("Could not create application window.", err)
	}
	// Set ApplicationWindow Properties
	appWindow.SetTitle("Basic Application.")
	appWindow.SetDefaultSize(400, 400)
	appWindow.Show()
}
