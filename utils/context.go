package utils

import (
	"context"
	"os"
	"os/signal"
)

// ContextWithSignal creates a context that cancels when the OS Interrupt signal is received
// Example:
// ```
// ctx := utils.ContextWithSignal(context.Background())
// go Run(ctx)
// <-ctx.Done()
// ```
func ContextWithSignal(parent context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(parent)

	go func() {
		<-c
		cancel()
	}()

	return ctx
}
