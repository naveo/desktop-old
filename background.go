package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

int
SetActivationPolicy(void) {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    return 0;
}
*/
import "C"
import (
	"context"
	"log"
	"os/exec"
)

func setActivationPolicy() {
	C.SetActivationPolicy()
}

func runNaveo(ctx context.Context) {
	// cmd := exec.CommandContext(ctx, "/Users/khalid/Downloads/naveoplayground/vmkit", "/Users/khalid/Downloads/naveoplayground/data.json")
	cmd := exec.CommandContext(ctx, "/Applications/naveo.app/Contents/MacOS/visorkit", "/Applications/naveo.app/Contents/MacOS/visorkit-config.json")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("naveo running")

	go func() {
		err = cmd.Wait()
		log.Printf("naveo finished with error: %v", err)
	}()
}

func naveoProcess(command <-chan string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		cmd := <-command
		log.Println(cmd)
		switch cmd {
		case "stop", "quit":
			cancel()
			return
		case "start":
			runNaveo(ctx)
		}
	}
}

func runPortkit(ctx context.Context) {
	// cmd := exec.CommandContext(ctx, "/Users/khalid/.naveo/portkit", "-config", "/Users/khalid/.naveo/config.json")
	cmd := exec.CommandContext(ctx, "/Applications/naveo.app/Contents/MacOS/portkit", "-config", "/Applications/naveo.app/Contents/MacOS/portkit-config.json")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("portkit running")

	Portkit.SetPid(cmd.Process.Pid)
	Portkit.SetState("running")

	go func() {
		err = cmd.Wait()
		log.Printf("portkit finished with error: %v", err)
	}()
}

func portkitProcess(command <-chan string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		cmd := <-command
		log.Println(cmd)
		switch cmd {
		case "stop", "quit":
			cancel()
			Portkit.SetState("")
			return
		case "start":
			runPortkit(ctx)
		}
	}
}
