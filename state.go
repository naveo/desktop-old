package main

import (
	"net"
	"time"

	"fyne.io/fyne/v2/widget"
)

// PortkitState represents the portkit background process state data.
type PortkitState struct {
	Pid   int
	State string
}

var Portkit PortkitState

// SetState sets the pid id for the portkit background process.
func (p *PortkitState) SetPid(pid int) {
	p.Pid = pid
}

// SetState sets the running state for the portkit background process.
func (p *PortkitState) SetState(state string) {
	p.State = state
}

// updateNaveoState tracks the state of Docker server running withing naveo.
func updateNaveoState(data *widget.Label, portkitExec chan string) {
	// initialize the address of the server.
	address := "naveo.local:2376"

	// open connection to Docker server withing naveo.
	_, err := net.DialTimeout("tcp", address, 3*time.Second)

	// check if connection was successfully established, if for somehow portkit process
	// is running before Docker stop it.
	if err != nil {
		if Portkit.State == "running" {
			portkitExec <- "stop"
		}
	} else {
		data.SetText("naveo is running")
		if Portkit.State != "running" {
			portkitExec <- "start"
		}
	}
}
