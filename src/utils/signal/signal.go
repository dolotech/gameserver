// +build !windows
// Copyright 2012-2017 Apcera Inc. All rights reserved.

package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// Signal Handling
func HandleSignals(f func(os.Signal)) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for {
			sig := <-c
			switch sig {
			case syscall.SIGTERM:
			case syscall.SIGINT:
				os.Exit(0)
			case syscall.SIGHUP:
				f(sig)
			}
		}
	}()
}
