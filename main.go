package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// visorkitExec channel to control the background process for naveo's virtualization system.
	visorkitExec := make(chan string)
	// portkitExec channel to control the background process for naveo's network port forwarding.
	portkitExec := make(chan string)

	a := app.New()
	window := a.NewWindow("Naveo Desktop")

	// fyne don't support hiding of dock icon on macOS as of yet,
	// this is a work around based on this user comment:
	// https://github.com/fyne-io/fyne/issues/3156#issuecomment-1295732800
	a.Lifecycle().SetOnStarted(func() {
		go func() {
			time.Sleep(1 * time.Second)
			setActivationPolicy()
		}()
	})

	if desk, ok := a.(desktop.App); ok {
		trayMenu := fyne.NewMenu("naveo",
			fyne.NewMenuItem("Dashboard", func() {
				window.Show()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Preference", func() {}),
			fyne.NewMenuItem("Check for Updates", func() {}),
			fyne.NewMenuItem("About", func() {}),
			fyne.NewMenuItem("Restart", func() {
				// stop all background processes before starting them again.
				visorkitExec <- "stop"
				portkitExec <- "stop"

				go naveoProcess(visorkitExec)
				visorkitExec <- "start"

				go portkitProcess(portkitExec)
			}))
		desk.SetSystemTrayMenu(trayMenu)
	}

	data := widget.NewLabel("")
	data.SetText("starting naveo")
	window.SetContent(data)

	// initial start for all background processes.
	go naveoProcess(visorkitExec)
	visorkitExec <- "start"
	go portkitProcess(portkitExec)

	// this goroutine runs keepalive on the Docker server running withing naveo.
	go func() {
		for range time.Tick(3 * time.Second) {
			updateNaveoState(data, portkitExec)
		}
	}()

	// change the behavior of the close button to hide the window instead.
	window.SetCloseIntercept(func() { window.Hide() })

	// start the naveo desktop app with main window hidden.
	a.Run()

	// stop all background processes before exiting the app.
	visorkitExec <- "quit"
	portkitExec <- "quit"
}
